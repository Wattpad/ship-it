package ecr

import (
	"context"
	"testing"

	"ship-it/internal/image"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/metrics/discard"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockReconciler struct {
	mock.Mock
}

func (m *MockReconciler) Reconcile(ctx context.Context, image *image.Ref) error {
	args := m.Called(ctx, image)
	return args.Error(0)
}

func TestMakeImageFromEvent(t *testing.T) {
	event := pushEvent{
		RepositoryName: "ship-it",
		RegistryId:     "723255503624",
		Tag:            "shipped",
	}
	assert.Exactly(t, image.Ref{
		Registry:   "723255503624.dkr.ecr.us-east-1.amazonaws.com",
		Repository: "ship-it",
		Tag:        "shipped",
	}, event.Image())
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
		assert.Equal(t, test.expected, validImageTagRegex.MatchString(test.inputString))
	}
}

func TestECRHandler(t *testing.T) {
	testListener := &ImageListener{
		logger: log.NewNopLogger(),
		timer:  discard.NewHistogram(),
	}

	mockReconciler := new(MockReconciler)

	err := testListener.handler(mockReconciler)(context.Background(), "some bad message")
	assert.Error(t, err)

	mockReconciler.On("Reconcile", mock.Anything, &image.Ref{
		Registry:   "723255503624.dkr.ecr.us-east-1.amazonaws.com",
		Repository: "monolith-php",
		Tag:        "78bc9ccf64eb838c6a0e0492ded722274925e2bd",
	}).Return(nil)

	inputJSON := `
{
	"eventTime": "2019-07-11T14:19:59Z", 
	"repositoryName": "monolith-php", 
	"tag": "78bc9ccf64eb838c6a0e0492ded722274925e2bd",
	"registryId": "723255503624"
}`
	err = testListener.handler(mockReconciler)(context.Background(), inputJSON)
	assert.NoError(t, err)
	mockReconciler.AssertExpectations(t)
}

type mockLogger struct {
	mock.Mock
}

func (m *mockLogger) Log(kvs ...interface{}) error {
	args := m.Called(kvs...)
	return args.Error(0)
}

func TestReconcilerNoRegisteredReleasesAffected(t *testing.T) {
	testErr := errors.Wrap(errNoRegisteredReleasesAffected, "wrapping with test context")

	mockLogger := new(mockLogger)
	mockLogger.On("Log", mock.AnythingOfType("string"), testErr.Error()).Return(nil)

	testListener := &ImageListener{
		logger: mockLogger,
		timer:  discard.NewHistogram(),
	}

	inputJSON := `
{
	"eventTime": "2019-07-11T14:19:59Z", 
	"repositoryName": "monolith-php", 
	"tag": "78bc9ccf64eb838c6a0e0492ded722274925e2bd",
	"registryId": "723255503624"
}`

	mockReconciler := new(MockReconciler)
	mockReconciler.On("Reconcile", mock.Anything, mock.Anything).Return(testErr)

	err := testListener.handler(mockReconciler)(context.Background(), inputJSON)
	assert.NoError(t, err)
	mockLogger.AssertExpectations(t)
}
