package ecr

import (
	"context"
	"fmt"
	"testing"

	"ship-it/internal"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"k8s.io/apimachinery/pkg/types"
)

type MockEditor struct {
	mock.Mock
}

func (m *MockEditor) UpdateAndReplace(ctx context.Context, releaseName string, image *internal.Image) error {
	args := m.Called(ctx, releaseName, image)
	return args.Error(0)
}

type MockIndexerService struct {
	mock.Mock
}

func (m *MockIndexerService) Lookup(repo string) ([]types.NamespacedName, error) {
	args := m.Called(repo)
	return args.Get(0).([]types.NamespacedName), args.Error(1)
}

func TestReconcileLookupFailure(t *testing.T) {
	mockEditor := new(MockEditor)
	mockIndexService := new(MockIndexerService)
	reconciler := NewReconciler(mockEditor, mockIndexService)
	inputImage := &internal.Image{
		Registry:   "723255503624.dkr.ecr.us-east-1.amazonaws.com",
		Repository: "bar",
		Tag:        "78bc9ccf64eb838c6a0e0492ded722274925e2bd",
	}

	mockIndexService.On("Lookup", "bar").Return([]types.NamespacedName{}, fmt.Errorf("some error finding release"))

	err := reconciler.Reconcile(context.Background(), inputImage)

	assert.Error(t, err)
	mockIndexService.AssertExpectations(t)

}

func TestReconcilerUpdateFailure(t *testing.T) {
	mockEditor := new(MockEditor)
	mockIndexService := new(MockIndexerService)

	reconciler := NewReconciler(mockEditor, mockIndexService)

	inputImage := &internal.Image{
		Registry:   "723255503624.dkr.ecr.us-east-1.amazonaws.com",
		Repository: "bar",
		Tag:        "78bc9ccf64eb838c6a0e0492ded722274925e2bd",
	}

	mockIndexService.On("Lookup", "bar").Return([]types.NamespacedName{
		{
			Namespace: "default",
			Name:      "bar",
		},
	}, error(nil))

	mockEditor.On("UpdateAndReplace", mock.Anything, "bar", inputImage).Return(fmt.Errorf("some update image error"))

	err := reconciler.Reconcile(context.Background(), inputImage)

	assert.Error(t, err)
	mockIndexService.AssertExpectations(t)
	mockEditor.AssertExpectations(t)
}

func TestReconcilerSuccess(t *testing.T) {
	mockEditor := new(MockEditor)
	mockIndexService := new(MockIndexerService)

	reconciler := NewReconciler(mockEditor, mockIndexService)

	inputImage := &internal.Image{
		Registry:   "723255503624.dkr.ecr.us-east-1.amazonaws.com",
		Repository: "bar",
		Tag:        "78bc9ccf64eb838c6a0e0492ded722274925e2bd",
	}

	mockIndexService.On("Lookup", "bar").Return([]types.NamespacedName{
		{
			Namespace: "default",
			Name:      "bar",
		},
	}, error(nil))

	mockEditor.On("UpdateAndReplace", mock.Anything, "bar", inputImage).Return(nil)

	err := reconciler.Reconcile(context.Background(), inputImage)

	assert.NoError(t, err)
	mockIndexService.AssertExpectations(t)
	mockEditor.AssertExpectations(t)
}
