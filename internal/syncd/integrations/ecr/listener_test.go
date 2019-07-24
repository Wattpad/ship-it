package ecr

import (
	"context"
	"testing"

	"ship-it/internal"

	"github.com/Wattpad/sqsconsumer"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/metrics"
	"github.com/stretchr/testify/assert"
)

func TestMakeImage(t *testing.T) {
	assert.Exactly(t, internal.Image{
		Registry:   "723255503624.dkr.ecr.us-east-1.amazonaws.com",
		Repository: "ship-it",
		Tag:        "shipped",
	}, makeImage("ship-it", "shipped", "723255503624"))
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
	var fakeHist metrics.Histogram
	mockSQSService := &sqsconsumer.SQSService{
		Svc:    mockSQSClient,
		URL:    &fakeURL,
		Logger: func(format string, args ...interface{}) {},
	}
	fakeListener := &ImageListener{
		logger:  log.NewNopLogger(),
		service: mockSQSService,
		timer:   fakeHist,
	}

	mockRepoService := new(MockRepoService)
	fakeReconciler := NewReconciler("Wattpad", "custom-resources/path", "foo", mockRepoService)

	err := fakeListener.handler(fakeReconciler)(context.Background(), `some bad message`)
	assert.Error(t, err)

	inputJSON := `
{
	"eventTime": "2019-07-11T14:19:59Z", 
	"repositoryName": "monolith-php", 
	"tag": "some invalid tag",
	"registryId": "723255503624"
}
`
	err = fakeListener.handler(fakeReconciler)(context.Background(), inputJSON)
	assert.Error(t, err)
}
