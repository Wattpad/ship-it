package k8s

import (
	"ship-it/internal/api/models"

	"ship-it/pkg/generated/clientset/versioned/typed/helmreleases.k8s.wattpad.com/v1alpha1"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/rest"
)

type K8sClient struct {
	helmreleases v1alpha1.HelmreleasesV1alpha1Interface
}

func New() (*K8sClient, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, err
	}

	client, err := v1alpha1.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return &K8sClient{client}, nil
}

func (k *K8sClient) ListAll(namespace string) ([]models.Release, error) {
	releaseList, err := k.helmreleases.HelmReleases(namespace).List(metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	var releases []models.Release
	for _, r := range releaseList.Items {
		releases = append(releases, models.Release{
			Name:    r.GetName(),
			Created: r.GetCreationTimestamp().Time,
		})
	}

	return releases, nil
}
