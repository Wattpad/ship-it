package syncd

import (
	"context"
	"testing"

	"ship-it/internal/image"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"k8s.io/helm/pkg/proto/hapi/chart"
)

type blockingChartListener struct {
	ShouldFail error
}

func (l blockingChartListener) Listen(ctx context.Context, _ RegistryChartReconciler) error {
	if l.ShouldFail != nil {
		return l.ShouldFail
	}
	<-ctx.Done()
	return ctx.Err()
}

type nopChartReconciler struct{}

func (nopChartReconciler) Reconcile(context.Context, *chart.Chart) error {
	return nil
}

type blockingImageListener struct {
	ShouldFail error
}

func (l blockingImageListener) Listen(ctx context.Context, _ ImageReconciler) error {
	if l.ShouldFail != nil {
		return l.ShouldFail
	}
	<-ctx.Done()
	return ctx.Err()
}

type nopImageReconciler struct{}

func (nopImageReconciler) Reconcile(context.Context, *image.Ref) error {
	return nil
}

// asserts the listeners are cancelled together
func TestSyncdNoFailures(t *testing.T) {
	s := New(
		blockingChartListener{},
		nopChartReconciler{},
		blockingImageListener{},
		nopImageReconciler{},
	)

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	err := s.Run(ctx)

	assert.Len(t, err, 2)
	assert.Contains(t, err, context.Canceled)
}

// asserts the chart listener is cancelled when the image listener fails
func TestSyncdImageListenerFailure(t *testing.T) {
	listenerErr := errors.New("image listener internal error")

	s := New(
		blockingChartListener{},
		nopChartReconciler{},
		blockingImageListener{ShouldFail: listenerErr},
		nopImageReconciler{},
	)

	err := s.Run(context.Background())

	assert.Len(t, err, 2)
	assert.Contains(t, err, listenerErr)
	assert.Contains(t, err, context.Canceled)
}

// asserts the image listener is cancelled when the chart listener fails
func TestSyncdChartListenerFailure(t *testing.T) {
	listenerErr := errors.New("chart listener internal error")

	s := New(
		blockingChartListener{ShouldFail: listenerErr},
		nopChartReconciler{},
		blockingImageListener{},
		nopImageReconciler{},
	)

	err := s.Run(context.Background())

	assert.Len(t, err, 2)
	assert.Contains(t, err, listenerErr)
	assert.Contains(t, err, context.Canceled)
}
