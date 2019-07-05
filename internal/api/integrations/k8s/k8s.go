package k8s

import (
	"context"
	"strconv"
	"time"

	"ship-it/internal/api/models"

	clientset "ship-it/pkg/generated/clientset/versioned"
	informers "ship-it/pkg/generated/informers/externalversions"

	"ship-it/pkg/generated/listers/k8s.wattpad.com/v1alpha1"

	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/rest"
)

type K8sClient struct {
	lister v1alpha1.HelmReleaseLister
}

func New(ctx context.Context, resync time.Duration) (*K8sClient, error) {
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
		autoDeploy, err := strconv.ParseBool(annotations["helmreleases.k8s.wattpad.com/autodeploy"])
		if err != nil {
			return nil, err
		}
		releases = append(releases, models.Release{
			Name:    r.GetName(),
			Created: r.GetCreationTimestamp().Time,
			AutoDeploy: autoDeploy,
			Owner: models.Owner{
				Squad: annotations["helmreleases.k8s.wattpad.com/squad"],
				Slack: annotations["helmreleases.k8s.wattpad.com/slack"],
			},
			Monitoring: models.Monitoring{
				Datadog: models.Datadog{
					Dashboard: annotations["helmreleases.k8s.wattpad.com/datadog"],
					Monitors: "",
				},
				Sumologic: annotations["helmreleases.k8s.wattpad.com/sumologic"],
			},
			Code: models.SourceCode{
				Github: annotations["helmreleases.k8s.wattpad.com/code"],
				Ref: "",
			},
		})
	}

	return releases, nil
}
