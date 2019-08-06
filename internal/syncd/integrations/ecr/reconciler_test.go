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

const inputYaml = `apiVersion: helmreleases.k8s.wattpad.com/v1alpha1
kind: HelmRelease
metadata:
  creationTimestamp: null
  name: example-microservice
spec:
  chart:
    path: microservice
    repository: wattpad.s3.amazonaws.com/helm-charts
    revision: HEAD
  releaseName: example-release
  values:
    image:
      repository: 723255503624.dkr.ecr.us-east-1.amazonaws.com/bar
      tag: baz
status: {}
`

const expectedYaml = `apiVersion: helmreleases.k8s.wattpad.com/v1alpha1
kind: HelmRelease
metadata:
  creationTimestamp: null
  name: example-microservice
spec:
  chart:
    path: microservice
    repository: wattpad.s3.amazonaws.com/helm-charts
    revision: HEAD
  releaseName: example-release
  values:
    image:
      repository: 723255503624.dkr.ecr.us-east-1.amazonaws.com/bar
      tag: 78bc9ccf64eb838c6a0e0492ded722274925e2bd
status: {}
`

type MockRepoService struct {
	mock.Mock
}

func (m *MockRepoService) UpdateAndReplace(ctx context.Context, releaseName string, image *internal.Image) error {
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
	mockRepoService := new(MockRepoService)
	mockIndexService := new(MockIndexerService)
	mockReconciler := NewReconciler(mockRepoService, mockIndexService)
	inputImage := &internal.Image{
		Registry:   "723255503624.dkr.ecr.us-east-1.amazonaws.com",
		Repository: "bar",
		Tag:        "78bc9ccf64eb838c6a0e0492ded722274925e2bd",
	}

	mockIndexService.On("Lookup", "bar").Return([]types.NamespacedName{}, fmt.Errorf("some error finding release"))

	err := mockReconciler.Reconcile(context.Background(), inputImage)

	assert.Error(t, err)
	mockIndexService.AssertExpectations(t)

}

func TestReconcilerUpdateFailure(t *testing.T) {
	mockRepoService := new(MockRepoService)
	mockIndexService := new(MockIndexerService)

	mockReconciler := NewReconciler(mockRepoService, mockIndexService)

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

	mockRepoService.On("UpdateAndReplace", mock.Anything, mock.Anything, inputImage).Return(fmt.Errorf("some update image error"))

	err := mockReconciler.Reconcile(context.Background(), inputImage)

	assert.Error(t, err)
	mockIndexService.AssertExpectations(t)
	mockRepoService.AssertExpectations(t)
}

func TestReconcilerSuccess(t *testing.T) {
	mockRepoService := new(MockRepoService)
	mockIndexService := new(MockIndexerService)

	mockReconciler := NewReconciler(mockRepoService, mockIndexService)

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

	mockRepoService.On("UpdateAndReplace", mock.Anything, mock.Anything, inputImage).Return(error(nil))

	err := mockReconciler.Reconcile(context.Background(), inputImage)

	assert.NoError(t, err)
	mockIndexService.AssertExpectations(t)
	mockRepoService.AssertExpectations(t)
}
