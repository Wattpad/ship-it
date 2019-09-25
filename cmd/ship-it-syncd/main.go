package main

import (
	"context"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"ship-it/internal/syncd"
	"ship-it/internal/syncd/config"
	"ship-it/internal/syncd/integrations/ecr"
	"ship-it/internal/syncd/integrations/github"
	"ship-it/internal/syncd/integrations/k8s"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/bradleyfalzon/ghinstallation"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/metrics/dogstatsd"
	gogithub "github.com/google/go-github/v26/github"
	"k8s.io/helm/pkg/helm"
)

func main() {
	logger := log.NewJSONLogger(log.NewSyncWriter(os.Stdout))
	logger = log.With(logger, "timestamp", log.DefaultTimestampUTC, "caller", log.DefaultCaller)

	logger.Log("event", "service.start")
	defer logger.Log("event", "service.stop")

	ctx, cancel := context.WithCancel(context.Background())

	signals := make(chan os.Signal, 2)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-signals
		cancel()
	}()

	cfg, err := config.FromEnv()
	if err != nil {
		logger.Log("error", err)
		os.Exit(1)
	}

	dd := dogstatsd.New("wattpad.ship-it.", logger)
	go dd.SendLoop(ctx, time.Tick(time.Second), "udp", cfg.DataDogAddress())

	// adapted from http.DefaultTransport
	defaultTransport := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
			DualStack: true,
		}).DialContext,
		ForceAttemptHTTP2:     true,
		MaxIdleConns:          10,
		IdleConnTimeout:       time.Minute,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: time.Second,
	}

	githubTransport, err := ghinstallation.New(
		defaultTransport,
		cfg.GithubAppID,
		cfg.GithubInstallationID,
		[]byte(cfg.GithubAppSecret),
	)
	if err != nil {
		logger.Log("error", err)
		os.Exit(1)
	}

	githubClient := gogithub.NewClient(&http.Client{
		Transport: githubTransport,
	})

	registryEditor := ecr.NewChartEditor(
		githubClient.Git,
		cfg.GithubOrg,
		cfg.OperationsRepository,
		cfg.ReleaseBranch,
		cfg.RegistryChartPath,
	)

	informer, err := k8s.NewInformer(ctx)
	if err != nil {
		logger.Log("error", err)
		os.Exit(1)
	}

	imageReconciler := ecr.NewReconciler(logger, registryEditor, informer)

	chartReconciler := github.NewReconciler(
		helm.NewClient(helm.Host(cfg.TillerHost)),
		cfg.Namespace,
		cfg.ReleaseName,
		cfg.HelmTimeout(),
	)

	imageListener, chartListener, err := initListeners(logger, githubClient, dd, cfg)
	if err != nil {
		logger.Log("error", err)
		os.Exit(1)
	}

	// TODO: Allow configurable image/chart sync implementations. For now
	// we'll just use our specific ecr+sqs/github+sqs implmentations.
	syncd := syncd.New(chartListener, chartReconciler, imageListener, imageReconciler)
	if err := syncd.Run(ctx); err != nil {
		logger.Log("error", err)
		os.Exit(1)
	}
}

func initListeners(l log.Logger, githubClient *gogithub.Client, dd *dogstatsd.Dogstatsd, cfg *config.Config) (syncd.ImageListener, syncd.RegistryChartListener, error) {
	syncHist := dd.NewTiming("syncd.time", 1)

	awsSession, err := session.NewSession(cfg.AWS())
	if err != nil {
		return nil, nil, err
	}

	sqsClient := sqs.New(awsSession)

	imageListener, err := ecr.NewListener(l, syncHist, cfg.EcrQueue, sqs.New(awsSession))
	if err != nil {
		return nil, nil, err
	}

	chartListener, err := github.NewListener(l, syncHist, cfg.GithubOrg, githubClient.Repositories, cfg.GithubQueue, sqsClient)

	return imageListener, chartListener, err
}
