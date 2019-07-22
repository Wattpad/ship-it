package k8s

import (
	"context"
	"time"

	"ship-it/internal/api/models"

	"k8s.io/client-go/rest"

	shipitv1beta1 "ship-it-operator/api/v1beta1"

	runtime "k8s.io/apimachinery/pkg/runtime"
	client "sigs.k8s.io/controller-runtime/pkg/client"
)

type K8sClient struct {
	client client.Client
}

func New(ctx context.Context, resync time.Duration) (*K8sClient, error) {
	scheme := runtime.NewScheme()
	shipitv1beta1.AddToScheme(scheme)

	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, err
	}

	cl, err := client.New(config, client.Options{
		Scheme: scheme,
	})

	return &K8sClient{
		client: cl,
	}, nil
}

func (k *K8sClient) ListAll(namespace string) ([]models.Release, error) {
	releaseList := &shipitv1beta1.HelmReleaseList{}

	err := k.client.List(context.Background(), releaseList, client.InNamespace(namespace))

	if err != nil {
		return nil, err
	}

	releases := make([]models.Release, 0, len(releaseList.Items))

	for _, r := range releaseList.Items {
		annotations := helmReleaseAnnotations(r.GetAnnotations())

		releases = append(releases, models.Release{
			Name:       r.ObjectMeta.GetName(),
			Created:    r.ObjectMeta.GetCreationTimestamp().Time,
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
