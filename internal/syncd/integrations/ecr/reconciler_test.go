package ecr

import (
	"context"
	"fmt"
	"path/filepath"
	"testing"

	"ship-it/internal"

	"github.com/google/go-github/v26/github"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
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

func (m *MockRepoService) GetFile(branch string, path string) (string, error) {
	args := m.Called(branch, path)
	return args.String(0), args.Error(1)
}

func (m *MockRepoService) UpdateFile(msg string, branch string, path string, fileContent []byte) (*github.RepositoryContentResponse, error) {
	args := m.Called(msg, branch, path, fileContent)
	return args.Get(0).(*github.RepositoryContentResponse), args.Error(1)
}

func TestReconcilerSuccess(t *testing.T) {
	mockRepoService := new(MockRepoService)
	fakeReconciler := &ImageReconciler{
		Org:          "Wattpad",
		ResourcePath: "foo/bar/resources",
		Branch:       "oof",
		RepoService:  mockRepoService,
	}

	inputImage := &internal.Image{
		Registry:   "723255503624.dkr.ecr.us-east-1.amazonaws.com",
		Repository: "bar",
		Tag:        "78bc9ccf64eb838c6a0e0492ded722274925e2bd",
	}
	path := filepath.Join(fakeReconciler.ResourcePath, inputImage.Repository+".yaml")
	mockRepoService.On("GetFile", fakeReconciler.Branch, path).Return(inputYaml, error(nil))
	mockRepoService.On("UpdateFile", fmt.Sprintf("Image Tag updated to: %s", inputImage.Tag), fakeReconciler.Branch, path, []byte(expectedYaml)).Return(&github.RepositoryContentResponse{}, error(nil))
	err := fakeReconciler.Reconcile(context.Background(), inputImage)
	assert.NoError(t, err)
	mockRepoService.AssertExpectations(t)

}

func TestReconcilerDownloadFailure(t *testing.T) {
	mockRepoService := new(MockRepoService)
	fakeReconciler := &ImageReconciler{
		Org:          "Wattpad",
		ResourcePath: "foo/bar/resources",
		Branch:       "oof",
		RepoService:  mockRepoService,
	}

	inputImage := &internal.Image{
		Registry:   "723255503624.dkr.ecr.us-east-1.amazonaws.com",
		Repository: "bar",
		Tag:        "78bc9ccf64eb838c6a0e0492ded722274925e2bd",
	}

	mockRepoService.On("GetFile", fakeReconciler.Branch, filepath.Join(fakeReconciler.ResourcePath, inputImage.Repository+".yaml")).Return("", fmt.Errorf("some error"))
	err := fakeReconciler.Reconcile(context.Background(), inputImage)
	assert.Error(t, err)
	mockRepoService.AssertExpectations(t)
}

func TestReconcilerUploadFailure(t *testing.T) {
	mockRepoService := new(MockRepoService)
	fakeReconciler := &ImageReconciler{
		Org:          "Wattpad",
		ResourcePath: "foo/bar/resources",
		Branch:       "oof",
		RepoService:  mockRepoService,
	}

	inputImage := &internal.Image{
		Registry:   "723255503624.dkr.ecr.us-east-1.amazonaws.com",
		Repository: "bar",
		Tag:        "78bc9ccf64eb838c6a0e0492ded722274925e2bd",
	}
	path := filepath.Join(fakeReconciler.ResourcePath, inputImage.Repository+".yaml")
	mockRepoService.On("GetFile", fakeReconciler.Branch, path).Return(inputYaml, error(nil))
	mockRepoService.On("UpdateFile", fmt.Sprintf("Image Tag updated to: %s", inputImage.Tag), fakeReconciler.Branch, path, []byte(expectedYaml)).Return((*github.RepositoryContentResponse)(nil), fmt.Errorf("some upload error"))
	err := fakeReconciler.Reconcile(context.Background(), inputImage)

	assert.Error(t, err)
	mockRepoService.AssertExpectations(t)
}

func TestReconcilerInvalidYaml(t *testing.T) {
	mockRepoService := new(MockRepoService)
	fakeReconciler := &ImageReconciler{
		Org:          "Wattpad",
		ResourcePath: "foo/bar/resources",
		Branch:       "oof",
		RepoService:  mockRepoService,
	}

	inputImage := &internal.Image{
		Registry:   "723255503624.dkr.ecr.us-east-1.amazonaws.com",
		Repository: "bar",
		Tag:        "78bc9ccf64eb838c6a0e0492ded722274925e2bd",
	}

	mockRepoService.On("GetFile", fakeReconciler.Branch, filepath.Join(fakeReconciler.ResourcePath, inputImage.Repository+".yaml")).Return("some malformed yaml", error(nil))
	err := fakeReconciler.Reconcile(context.Background(), inputImage)
	assert.Error(t, err)
	mockRepoService.AssertExpectations(t)
}

func TestNewReconciler(t *testing.T) {
	mockRepoService := new(MockRepoService)
	reconciler := NewReconciler("Wattpad", "custom-resources/path", "foo", mockRepoService)
	assert.Equal(t, &ImageReconciler{
		"Wattpad",
		"custom-resources/path",
		"foo",
		mockRepoService,
	}, reconciler)
}
