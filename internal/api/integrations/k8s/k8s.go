package k8s

import (
	"context"
	"time"

	"ship-it/internal/api/models"

	clientset "ship-it/pkg/generated/clientset/versioned"
	informers "ship-it/pkg/generated/informers/externalversions"

	"ship-it/pkg/generated/listers/helmreleases.k8s.wattpad.com/v1alpha1"

	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/rest"
)

type K8sClient struct {
	lister v1alpha1.HelmReleaseLister
}

func New(ctx context.Context) (*K8sClient, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, err
	}

	client, err := clientset.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	exit := make(chan struct{})

	go func() {
		<-ctx.Done()
		exit <- struct{}{}
	}()

	factory := informers.NewSharedInformerFactory(client, 30*time.Second)
	factory.Start(exit)

	return &K8sClient{
		lister: factory.Helmreleases().V1alpha1().HelmReleases().Lister(),
	}, nil
}

func (k *K8sClient) ListAll(namespace string) ([]models.Release, error) {
	releaseList, err := k.lister.HelmReleases(namespace).List(labels.Everything())
	if err != nil {
		return nil, err
	}

	releases := make([]models.Release, 0, len(releaseList))
	for _, r := range releaseList {
		releases = append(releases, models.Release{
			Name:    r.GetName(),
			Created: r.GetCreationTimestamp().Time,
		})
	}

	return releases, nil
}
