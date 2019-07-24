package ecrconsumer

import (
	"context"
	"testing"
	"time"

	"github.com/google/go-github/v26/github"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockedObject struct {
	mock.Mock
}

func TestParseMessage(t *testing.T) {
	inputJSON := `
{
	"eventTime": "2019-07-11T14:19:59Z", 
	"repositoryName": "monolith-php", 
	"tag": "78bc9ccf64eb838c6a0e0492ded722274925e2bd",
	"registryId": "723255503624"
}
`

	deployTime, err := time.Parse(time.RFC3339, "2019-07-11T14:19:59Z")
	assert.NoError(t, err)

	expectedMessage := SQSMessage{
		EventTime:      deployTime,
		RepositoryName: "monolith-php",
		Tag:            "78bc9ccf64eb838c6a0e0492ded722274925e2bd",
		RegistryId:     "723255503624",
	}

	inputMessage, err := parseMsg(inputJSON)
	assert.NoError(t, err)
	assert.Exactly(t, expectedMessage, *inputMessage)
}

func (m *MockedObject) GetFile(branch string, path string) ([]byte, error) {
	args := m.Called(branch, path)
	return args.Get(0).([]byte), args.Error(1)
}

func (m *MockedObject) UpdateFile(msg string, branch string, path string, fileContent []byte) (*github.RepositoryContentResponse, error) {
	args := m.Called(msg, branch, path, fileContent)
	return args.Get(0).(*github.RepositoryContentResponse), args.Error(1)
}

func TestValidateTag(t *testing.T) {
	tests := []struct {
		inputString string
		expected    bool
	}{
		{"78bc9ccf64eb838c6a0e0492ded722274925e2bd", true},
		{"latest", false},
		{"78bc9ccf64eb838c6a0e0492ded722274925E2ND", false},
	}
	for _, test := range tests {
		assert.Equal(t, test.expected, validateTag(test.inputString))
	}
}

func TestProcessMessage(t *testing.T) {

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

	const inputJSON = `
{
	"eventTime": "2019-07-11T14:19:59Z", 
	"repositoryName": "bar", 
	"tag": "78bc9ccf64eb838c6a0e0492ded722274925e2bd",
	"registryId": "723255503624"
}
`

	mockGit := new(MockedObject)
	mockGit.On("GetFile", "master", "custom-resources").Return([]byte(inputYaml), error(nil))
	mockGit.On("UpdateFile", "Image Tag updated to: 78bc9ccf64eb838c6a0e0492ded722274925e2bd", "master", "custom-resources/bar.yaml", []byte(expectedYaml)).Return(&github.RepositoryContentResponse{}, error(nil))
	handler := processMessage(mockGit, "custom-resources", "master")
	err := handler(context.Background(), inputJSON)
	assert.NoError(t, err)
}
