package ecr

import (
	"context"
	"fmt"
	"testing"

	"ship-it/internal/image"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"k8s.io/apimachinery/pkg/types"
)

type MockReleaseEditor struct {
	mock.Mock
}

func (m *MockReleaseEditor) Edit(ctx context.Context, releases []types.NamespacedName, image *image.Ref) error {
	args := m.Called(ctx, releases, image)
	return args.Error(0)
}

type MockReleaseIndexer struct {
	mock.Mock
}

func (m *MockReleaseIndexer) Lookup(image *image.Ref) ([]types.NamespacedName, error) {
	args := m.Called(image)
	return args.Get(0).([]types.NamespacedName), args.Error(1)
}

func TestReconcileLookupFailure(t *testing.T) {
	mockEditor := new(MockReleaseEditor)
	mockIndexer := new(MockReleaseIndexer)
	reconciler := NewReconciler(mockEditor, mockIndexer)
	inputImage := &image.Ref{
		Registry:   "723255503624.dkr.ecr.us-east-1.amazonaws.com",
		Repository: "bar",
		Tag:        "78bc9ccf64eb838c6a0e0492ded722274925e2bd",
	}

	mockIndexer.On("Lookup", inputImage).Return([]types.NamespacedName{}, fmt.Errorf("some error finding release"))

	err := reconciler.Reconcile(context.Background(), inputImage)

	assert.Error(t, err)
	mockIndexer.AssertExpectations(t)
}

func TestReconcilerUpdateFailure(t *testing.T) {
	mockEditor := new(MockReleaseEditor)
	mockIndexer := new(MockReleaseIndexer)

	reconciler := NewReconciler(mockEditor, mockIndexer)

	inputImage := &image.Ref{
		Registry:   "723255503624.dkr.ecr.us-east-1.amazonaws.com",
		Repository: "bar",
		Tag:        "78bc9ccf64eb838c6a0e0492ded722274925e2bd",
	}

	releaseNames := []types.NamespacedName{
		{
			Namespace: "default",
			Name:      "bar",
		},
	}

	mockIndexer.On("Lookup", inputImage).Return(releaseNames, error(nil))

	mockEditor.On("Edit", mock.Anything, releaseNames, inputImage).Return(fmt.Errorf("some update image error"))

	err := reconciler.Reconcile(context.Background(), inputImage)

	assert.Error(t, err)
	mockIndexer.AssertExpectations(t)
	mockEditor.AssertExpectations(t)
}

func TestReconcilerSuccess(t *testing.T) {
	mockEditor := new(MockReleaseEditor)
	mockIndexer := new(MockReleaseIndexer)

	reconciler := NewReconciler(mockEditor, mockIndexer)

	inputImage := &image.Ref{
		Registry:   "723255503624.dkr.ecr.us-east-1.amazonaws.com",
		Repository: "bar",
		Tag:        "78bc9ccf64eb838c6a0e0492ded722274925e2bd",
	}

	releaseNames := []types.NamespacedName{
		{
			Namespace: "default",
			Name:      "bar",
		},
	}

	mockIndexer.On("Lookup", inputImage).Return(releaseNames, error(nil))

	mockEditor.On("Edit", mock.Anything, releaseNames, inputImage).Return(nil)

	err := reconciler.Reconcile(context.Background(), inputImage)

	assert.NoError(t, err)
	mockIndexer.AssertExpectations(t)
	mockEditor.AssertExpectations(t)
}
