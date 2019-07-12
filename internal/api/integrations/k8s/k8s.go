package k8s

import (
	"context"
	"strconv"
	"strings"
	"time"

	"ship-it/internal"
	"ship-it/internal/api/models"

	clientset "ship-it/pkg/generated/clientset/versioned"
	informers "ship-it/pkg/generated/informers/externalversions"

	"ship-it/pkg/apis/k8s.wattpad.com/v1alpha1"
	listerv1alpha1 "ship-it/pkg/generated/listers/k8s.wattpad.com/v1alpha1"

	"github.com/go-kit/kit/log"

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

func annotationFor(k string) string {
	return v1alpha1.Resource("helmreleases").String() + "/" + k
}

func (k *K8sClient) ListAll(namespace string) ([]models.Release, error) {
	releaseList, err := k.lister.HelmReleases(namespace).List(labels.Everything())
	if err != nil {
		return nil, err
	}

	releases := make([]models.Release, 0, len(releaseList))
	for _, r := range releaseList {
		annotations := r.GetAnnotations()
		autoDeploy, err := strconv.ParseBool(annotations[annotationFor("autodeploy")])

		k.logger.Log(err)

		codeURL := annotations[annotationFor("code")]
		serviceName := r.GetName()
		image := GetImageForRepo(serviceName, r.Spec.Values, k.logger)
		releases = append(releases, models.Release{
			Name:       serviceName,
			Created:    r.GetCreationTimestamp().Time,
			AutoDeploy: autoDeploy,
			Owner: models.Owner{
				Squad: annotations[annotationFor("squad")],
				Slack: annotations[annotationFor("slack")],
			},
			Monitoring: models.Monitoring{
				Datadog: models.Datadog{
					Dashboard: annotations[annotationFor("datadog")],
				},
				Sumologic: annotations[annotationFor("sumologic")],
			},
			Code: models.SourceCode{
				Github: findVCSRepo(codeURL),
				Ref:    image.Tag,
			},
			Artifacts: models.Artifacts{
				Chart: models.HelmArtifact{
					Path:    r.Spec.Chart.Path,
					Version: r.Spec.Chart.Revision,
				},
				Docker: models.DockerArtifact{
					Image: image.Repository,
					Tag:   image.Tag,
				},
			},
		})
	}

	return releases, nil
}

func findVCSRepo(url string) string {
	arr := strings.Split(url, "/")
	if len(arr) > 4 {
		return arr[4]
	}
	return ""
}

func GetImageForRepo(repo string, vals map[string]interface{}, logger log.Logger) internal.Image {
	arr := internal.GetImagePath(vals, repo)
	if len(arr) == 0 {
		logger.Log("No Image Found for repository")
		return internal.Image{}
	}
	imgVals := internal.Table(vals, arr)
	img, err := internal.ParseImage(imgVals["repository"].(string), imgVals["tag"].(string))
	if err != nil {
		logger.Log("Unable to parse valid image")
		return internal.Image{}
	}
	return *img
}
