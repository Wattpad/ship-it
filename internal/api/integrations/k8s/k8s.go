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

	"ship-it/pkg/apis/k8s.wattpad.com/v1alpha1"
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

		autoDeployStr := annotations[annotationFor("autodeploy")]
		autoDeploy, err := strconv.ParseBool(autoDeployStr)
		if err != nil {
			errors.Wrapf(err, `Unable to parse "%s" as boolean`, autoDeployStr)
			k.logger.Log("error", err)
		}

		codeURL := annotations[annotationFor("code")]
		gitRepo, err := findVCSRepo(codeURL)
		if err != nil {
			k.logger.Log("error", err)
		}

		releaseName := r.GetName()

		image, err := GetImageForRepo(releaseName, r.Spec.Values)
		if err != nil {
			k.logger.Log("error", err)
		}

		releases = append(releases, models.Release{
			Name:       releaseName,
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
				Github: *gitRepo,
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

func findVCSRepo(addr string) (*string, error) {
	address, err := url.Parse(addr)
	if err != nil {
		return nil, err
	}
	arr := strings.Split(address.Path, "/")
	if len(arr) > 2 {
		return &arr[2], nil
	}
	return nil, fmt.Errorf("")
}

func GetImageForRepo(repo string, vals map[string]interface{}) (*internal.Image, error) {
	arr := internal.GetImagePath(vals, repo)
	if len(arr) == 0 {
		return nil, fmt.Errorf("Image not found")
	}
	imgVals := internal.Table(vals, arr)
	img, err := internal.ParseImage(imgVals["repository"].(string), imgVals["tag"].(string))
	if err != nil {
		return nil, fmt.Errorf("invalid image: %v", err)
	}
	return img, nil
}
