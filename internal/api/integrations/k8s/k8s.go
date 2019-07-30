package k8s

import (
	"context"
	"fmt"
	shipitv1beta1 "ship-it-operator/api/v1beta1"
	"ship-it/internal"
	"ship-it/internal/api/models"
	"ship-it/internal/unstructured"

	"github.com/pkg/errors"
	runtime "k8s.io/apimachinery/pkg/runtime"
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

func (k *K8sClient) ListAll(ctx context.Context, namespace string) ([]models.Release, error) {
	var releaseList shipitv1beta1.HelmReleaseList

	err := k.client.List(ctx, &releaseList, client.InNamespace(namespace))

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
				Ref:    "",
			},
			Artifacts: models.Artifacts{
				Chart: models.HelmArtifact{
					Path:       r.Spec.Chart.Path,
					Repository: r.Spec.Chart.Repository,
					Version:    r.Spec.Chart.Revision,
				},
				Docker: dockerArtifacts(r),
			},
		})
	}
	return releases, nil
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
