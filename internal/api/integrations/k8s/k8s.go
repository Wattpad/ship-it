package k8s

import (
	"context"
	"time"

	"ship-it/internal/api/models"

	clientset "ship-it/pkg/generated/clientset/versioned"
	informers "ship-it/pkg/generated/informers/externalversions"

	listerv1alpha1 "ship-it/pkg/generated/listers/k8s.wattpad.com/v1alpha1"

	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/rest"
)

type K8sClient struct {
	lister listerv1alpha1.HelmReleaseLister
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
		annotations := helmReleaseAnnotations(r.GetAnnotations())

		releases = append(releases, models.Release{
			Name:       r.GetName(),
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
			},
			Artifacts: models.Artifacts{
				Chart: models.HelmArtifact{
					Path:       r.Spec.Chart.Path,
					Repository: r.Spec.Chart.Repository,
					Version:    r.Spec.Chart.Revision,
				},
			},
		})
	}

	return releases, nil
}
