package ecr

import (
	"context"
	"testing"

	"ship-it/internal"

	"github.com/Wattpad/sqsconsumer"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/metrics/discard"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockReconciler struct {
	mock.Mock
}

func (m *MockReconciler) Reconcile(ctx context.Context, image *internal.Image) error {
	args := m.Called(ctx, image)
	return args.Error(0)
}

func TestMakeImageFromEvent(t *testing.T) {
	event := &pushEvent{
		RepositoryName: "ship-it",
		RegistryId:     "723255503624",
		Tag:            "shipped",
	}
	assert.Exactly(t, &internal.Image{
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
	mockSQSClient := new(MockSQS)
	fakeURL := "www.wattpad.com"
	mockSQSService := &sqsconsumer.SQSService{
		Svc:    mockSQSClient,
		URL:    &fakeURL,
		Logger: func(format string, args ...interface{}) {},
	}
	testListener := &ImageListener{
		logger:  log.NewNopLogger(),
		service: mockSQSService,
		timer:   discard.NewHistogram(),
	}

	mockReconciler := new(MockReconciler)

	err := testListener.handler(mockReconciler)(context.Background(), `some bad message`)
	assert.Error(t, err)

	mockReconciler.On("Reconcile", mock.Anything, &internal.Image{
		Registry:   "723255503624.dkr.ecr.us-east-1.amazonaws.com",
		Repository: "monolith-php",
		Tag:        "78bc9ccf64eb838c6a0e0492ded722274925e2bd",
	}).Return(error(nil))

	inputJSON := `
{
	"eventTime": "2019-07-11T14:19:59Z", 
	"repositoryName": "monolith-php", 
	"tag": "78bc9ccf64eb838c6a0e0492ded722274925e2bd",
	"registryId": "723255503624"
}
`
	err = testListener.handler(mockReconciler)(context.Background(), inputJSON)
	assert.NoError(t, err)
	mockReconciler.AssertExpectations(t)
}
