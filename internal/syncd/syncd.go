package syncd

import (
	"context"
	"sync"

	"ship-it/internal"

	"k8s.io/helm/pkg/proto/hapi/chart"
)

// ImageReconciler reconciles a new image with the state of ship-it's service
// registry chart. For example, by updating chart values in a remote
// repository to use the new image.
type ImageReconciler interface {
	Reconcile(context.Context, *internal.Image) error
}

// ImageListener models a stream of `*Image` events for service images. It calls
// the ImageReconciler to reconcile new images with ship-it's service registry chart.
type ImageListener interface {
	Listen(context.Context, ImageReconciler) error
}

// RegistryChartReconciler is reconciles the registry chart with the
// kubernetes cluster state. For example, by deploying the chart to a cluster.
type RegistryChartReconciler interface {
	Reconcile(context.Context, *chart.Chart) error
}

// RegistryChartListener models a stream of `*chart.Chart` events for ship-it's
// service registry chart. It calls the RegistryChartReconciler to reconcile
// each new chart with the kubernetes cluster state.
type RegistryChartListener interface {
	Listen(context.Context, RegistryChartReconciler) error
}

// Syncd facilitates background synchronization between a docker image registry,
// ship-it's service registry helm chart, and kubernetes cluster state.
type Syncd struct {
	chartListener   RegistryChartListener
	chartReconciler RegistryChartReconciler

	imageListener   ImageListener
	imageReconciler ImageReconciler
}

func New(cl RegistryChartListener, cr RegistryChartReconciler, il ImageListener, ir ImageReconciler) *Syncd {
	return &Syncd{
		chartListener:   cl,
		chartReconciler: cr,

		imageListener:   il,
		imageReconciler: ir,
	}
}

func (s *Syncd) Run(ctx context.Context) error {
	var wg sync.WaitGroup
	wg.Add(2)

	errs := make(chan error, 2)

	// cancel both listeners if either one exits
	ctx, cancel := context.WithCancel(ctx)

	go func() {
		errs <- s.chartListener.Listen(ctx, s.chartReconciler)
		wg.Done()
		cancel()
	}()

	go func() {
		errs <- s.imageListener.Listen(ctx, s.imageReconciler)
		wg.Done()
		cancel()
	}()

	go func() {
		wg.Wait()
		close(errs)
	}()

	var merr multiError
	for err := range errs {
		merr.Add(err)
	}

	return merr
}
