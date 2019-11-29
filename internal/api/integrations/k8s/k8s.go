package k8s

import (
	"context"
	shipitv1beta1 "ship-it-operator/api/v1beta1"
	"ship-it/internal/api/models"
	"ship-it/internal/unstructured"

	runtime "k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/rest"
	client "sigs.k8s.io/controller-runtime/pkg/client"
)

type K8sClient struct {
	client client.Client
}

func New() (*K8sClient, error) {
	scheme := runtime.NewScheme()
	shipitv1beta1.AddToScheme(scheme)

	config, err := rest.InClusterConfig()

	if err != nil {
		return nil, err
	}

	cl, err := client.New(config, client.Options{
		Scheme: scheme,
	})

	if err != nil {
		return nil, err
	}

	return &K8sClient{
		client: cl,
	}, nil
}

func (k *K8sClient) Get(ctx context.Context, namespace, name string) (*models.Release, error) {
	var release shipitv1beta1.HelmRelease

	key := types.NamespacedName{
		Name:      name,
		Namespace: namespace,
	}

	err := k.client.Get(ctx, key, &release)
	if err != nil {
		return nil, err
	}

	modelRelease := transform(release)
	return &modelRelease, nil
}

func (k *K8sClient) List(ctx context.Context, namespace string) ([]models.Release, error) {
	var releaseList shipitv1beta1.HelmReleaseList

	err := k.client.List(ctx, &releaseList, client.InNamespace(namespace))
	if err != nil {
		return nil, err
	}

	releases := make([]models.Release, 0, len(releaseList.Items))

	for _, r := range releaseList.Items {
		releases = append(releases, transform(r))
	}

	return releases, nil
}

func transform(r shipitv1beta1.HelmRelease) models.Release {
	annotations := r.Annotations()

	return models.Release{
		Name:         r.ObjectMeta.GetName(),
		Created:      r.ObjectMeta.GetCreationTimestamp().Time,
		LastDeployed: r.Status.GetCondition().LastUpdateTime.Time,
		AutoDeploy:   annotations.AutoDeploy(),
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
				Path:       r.Spec.Chart.Name,
				Repository: r.Spec.Chart.Repository,
				Version:    r.Spec.Chart.Version,
			},
			Docker: dockerArtifacts(r),
		},
		Status: r.Status.GetCondition().Type,
	}
}

func dockerArtifacts(hr shipitv1beta1.HelmRelease) []models.DockerArtifact {
	var artifacts []models.DockerArtifact

	// find all "image" sections in the release's values, transforming each
	// one into a docker artifact
	unstructured.FindAll(hr.HelmValues(), "image", func(x interface{}) {
		if img, ok := x.(map[string]interface{}); ok {
			repo, ok := img["repository"].(string)
			if !ok {
				return
			}

			tag, ok := img["tag"].(string)
			if !ok {
				return
			}

			artifacts = append(artifacts, models.DockerArtifact{
				Image: repo,
				Tag:   tag,
			})
		}
	})

	return artifacts
}
