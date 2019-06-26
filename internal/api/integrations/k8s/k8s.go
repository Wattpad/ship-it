package k8s

import (
	"ship-it/internal/api/models"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	kv1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/rest"
)

func (k *K8sClient) ListAll(namespace string) ([]models.Release, error) {
	configMapList, err := k.core.ConfigMaps(namespace).List(metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	var releases []models.Release
	for _, configMap := range configMapList.Items {
		r := models.Release{
			Name:    configMap.GetName(),
			Created: configMap.GetCreationTimestamp().Time,
		}

		releases = append(releases, r)
	}

	return releases, nil
}

type K8sClient struct {
	core kv1.CoreV1Interface
}

func New() (*K8sClient, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, err
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return &K8sClient{clientset.CoreV1()}, nil
}
