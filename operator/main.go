/*

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"flag"
	"os"
	"time"

	shipitv1beta1 "ship-it-operator/api/v1beta1"

	"ship-it-operator/chartdownloader"
	"ship-it-operator/controllers"

	"ship-it-operator/notifications"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/helm/pkg/helm"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	// +kubebuilder:scaffold:imports
)

var (
	scheme   = runtime.NewScheme()
	setupLog = ctrl.Log.WithName("setup")
)

func init() {

	shipitv1beta1.AddToScheme(scheme)
	// +kubebuilder:scaffold:scheme
}

func main() {
	var (
		awsRegion            string
		gracePeriod          time.Duration
		metricsAddr          string
		targetNamespace      string
		watchNamespace       string
		tillerAddr           string
		enableLeaderElection bool
		slackChannel         string
		slackToken           string
	)

	flag.StringVar(&awsRegion, "aws-region", "us-east-1", "The AWS region where the operator's chart repository is hosted")
	flag.StringVar(&metricsAddr, "metrics-addr", ":8080", "The address the metric endpoint binds to.")
	flag.StringVar(&watchNamespace, "watch-namespace", "default", "The cluster namespace where the operator will watch HelmRelease resources")
	flag.StringVar(&targetNamespace, "target-namespace", "default", "The cluster namespace where the operator will deploy releases")
	flag.StringVar(&tillerAddr, "tiller-address", "localhost:44134", "The cluster address of the tiller service")
	flag.DurationVar(&gracePeriod, "grace-period", 10*time.Second, "The duration the operator will wait before checking a release's status after reconciling")
	flag.BoolVar(&enableLeaderElection, "enable-leader-election", false,
		"Enable leader election for controller manager. Enabling this will ensure there is only one active controller manager.")
	flag.StringVar(&slackChannel, "slack-channel", "", "The channel to send Slack notifications to")
	flag.StringVar(&slackToken, "slack-token", "", "API token for Slack")

	flag.Parse()

	ctrl.SetLogger(zap.Logger(true))

	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
		Scheme:             scheme,
		MetricsBindAddress: metricsAddr,
		Namespace:          watchNamespace,
		LeaderElection:     enableLeaderElection,
	})
	if err != nil {
		setupLog.Error(err, "unable to start manager")
		os.Exit(1)
	}

	session, err := session.NewSession(&aws.Config{
		Region: aws.String(awsRegion),
	})
	if err != nil {
		setupLog.Error(err, "unable to create AWS session")
		os.Exit(1)
	}

	// S3 is currently the only supported chart repository type
	downloaders := map[string]chartdownloader.ChartDownloader{
		"s3": chartdownloader.NewS3Downloader(s3manager.NewDownloader(session)),
	}

	reconciler := controllers.NewHelmReleaseReconciler(
		ctrl.Log,
		mgr.GetClient(),
		notifications.NewSlack(slackToken, slackChannel),
		helm.NewClient(helm.Host(tillerAddr)),
		chartdownloader.New(downloaders),
		mgr.GetEventRecorderFor("ship-it"),
		controllers.Namespace(targetNamespace),
		controllers.GracePeriod(gracePeriod),
	)

	setupLog.Info("setting up HelmRelease controller")
	if err := reconciler.SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "HelmRelease")
		os.Exit(1)
	}
	// +kubebuilder:scaffold:builder

	setupLog.Info("starting manager")
	if err := mgr.Start(ctrl.SetupSignalHandler()); err != nil {
		setupLog.Error(err, "problem running manager")
		os.Exit(1)
	}
}
