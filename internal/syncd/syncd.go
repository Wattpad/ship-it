package syncd

import (
	"context"
	"sync"

	"k8s.io/helm/pkg/proto/hapi/chart"
)

// ImageReconciler is an event handler callback that reconciles a new image with
// some external state. For example, updating the values of a helm chart in a
// remote repository to use the new image's tag.
type ImageReconciler interface {
	Reconcile(context.Context, *Image) error
}

// ImageListener models a stream of `*Image` events. It calls the
// ImageReconciler to reconcile each new image with some external state.
type ImageListener interface {
	Listen(context.Context, ImageReconciler) error
}

// ChartReconciler is an event handler callback that reconciles a new chart with
// some external state. For example, deploying the chart to a kubernetes cluster.
type ChartReconciler interface {
	Reconcile(context.Context, *chart.Chart) error
}

// ChartListener models a stream of `*chart.Chart` events. It calls the
// ChartReconciler to reconcile each new chart with some external state.
type ChartListener interface {
	Listen(context.Context, ChartReconciler) error
}

// Syncd facilitates background synchronization between an image registry, a
// helm chart repository and some external state, like a kubernetes cluster.
type Syncd struct {
	chartListener   ChartListener
	chartReconciler ChartReconciler

	imageListener   ImageListener
	imageReconciler ImageReconciler
}

func New(cl ChartListener, cr ChartReconciler, il ImageListener, ir ImageReconciler) *Syncd {
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
