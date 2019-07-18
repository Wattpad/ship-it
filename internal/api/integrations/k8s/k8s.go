package k8s

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"

	"ship-it/internal"
	"ship-it/internal/api/models"

	clientset "ship-it/pkg/generated/clientset/versioned"
	informers "ship-it/pkg/generated/informers/externalversions"

	listerv1alpha1 "ship-it/pkg/generated/listers/k8s.wattpad.com/v1alpha1"

	"github.com/go-kit/kit/log"

	"github.com/pkg/errors"

	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/rest"
)

type K8sClient struct {
	lister listerv1alpha1.HelmReleaseLister
	logger log.Logger
}

func New(ctx context.Context, resync time.Duration, logger log.Logger) (*K8sClient, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, err
	}

	client, err := clientset.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	factory := informers.NewSharedInformerFactory(client, resync)

	helmReleaseLister := factory.Helmreleases().V1alpha1().HelmReleases().Lister()

	// factory must be started after all informers/listers have been created
	factory.Start(ctx.Done())

	return &K8sClient{
		lister: helmReleaseLister,
		logger: logger,
	}, nil
}

func (k *K8sClient) ListAll(namespace string) ([]models.Release, error) {
	releaseList, err := k.lister.HelmReleases(namespace).List(labels.Everything())
	if err != nil {
		return nil, err
	}

	releases := make([]models.Release, 0, len(releaseList))
	for _, r := range releaseList {
		annotations := r.GetAnnotations()

		gitRepo, err := findVCSRepo(annotations.Code())
		if err != nil {
			k.logger.Log("error", err)
		}

		releaseName := r.GetName()

		image, err := GetImageForRepo(releaseName, r.Spec.Values)
		if err != nil {
			k.logger.Log("error", err)
		}

		repository := ""
		tag := ""
		if image != nil {
			repository = image.Repository
			tag = image.Tag
		}

		releases = append(releases, models.Release{
			Name:       releaseName,
			Created:    r.GetCreationTimestamp().Time,
			AutoDeploy: annotations.AutoDeploy(),
			Owner: models.Owner{
				Squad: annotations.Squad(),
				Slack: annotations.Slack(),
			},
			Monitoring: models.Monitoring{
				Datadog: models.Datadog{
					Dashboard: annotations.Datadog(),
				},
				Sumologic: annotations.Sumologic(),
			},
			Code: models.SourceCode{
				Github: annotations.Code(),
				Ref:    tag,
			},
			Artifacts: models.Artifacts{
				Chart: models.HelmArtifact{
					Path:       r.Spec.Chart.Path,
					Repository: r.Spec.Chart.Repository,
					Version:    r.Spec.Chart.Revision,
				},
				Docker: models.DockerArtifact{
					Image: repository,
					Tag:   tag,
				},
			},
		})
	}
	return releases, nil
}

func findVCSRepo(addr string) (string, error) {
	address, err := url.Parse(addr)
	if err != nil {
		return "", errors.Wrap(err, "url parsing failure")
	}
	arr := strings.Split(address.Path, "/")
	if len(arr) > 2 {
		return arr[2], nil
	}
	return "", fmt.Errorf("no repository found")
}

func GetImageForRepo(repo string, vals map[string]interface{}) (*internal.Image, error) {
	arr := internal.GetImagePath(vals, repo)
	if len(arr) == 0 {
		return nil, fmt.Errorf("image not found")
	}
	imgVals := internal.Table(vals, arr)
	img, err := internal.ParseImage(imgVals["repository"].(string), imgVals["tag"].(string))
	if err != nil {
		return nil, errors.Wrap(err, "invalid image")
	}
	return img, nil
}
