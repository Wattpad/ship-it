package ecr

import (
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/stretchr/testify/mock"
)

type MockSQS struct {
	mock.Mock
}

type SQSAPI interface {
	ChangeMessageVisibilityBatch(input *sqs.ChangeMessageVisibilityBatchInput) (*sqs.ChangeMessageVisibilityBatchOutput, error)
	CreateQueue(input *sqs.CreateQueueInput) (*sqs.CreateQueueOutput, error)
	DeleteMessageBatch(input *sqs.DeleteMessageBatchInput) (*sqs.DeleteMessageBatchOutput, error)
	GetQueueUrl(input *sqs.GetQueueUrlInput) (*sqs.GetQueueUrlOutput, error)
	ReceiveMessage(input *sqs.ReceiveMessageInput) (*sqs.ReceiveMessageOutput, error)
}

func (m *MockSQS) ChangeMessageVisibilityBatch(input *sqs.ChangeMessageVisibilityBatchInput) (*sqs.ChangeMessageVisibilityBatchOutput, error) {
	args := m.Called(input)
	return args.Get(0).(*sqs.ChangeMessageVisibilityBatchOutput), args.Error(1)
}

func (m *MockSQS) CreateQueue(input *sqs.CreateQueueInput) (*sqs.CreateQueueOutput, error) {
	args := m.Called(input)
	return args.Get(0).(*sqs.CreateQueueOutput), args.Error(1)
}

func (m *MockSQS) DeleteMessageBatch(input *sqs.DeleteMessageBatchInput) (*sqs.DeleteMessageBatchOutput, error) {
	args := m.Called(input)
	return args.Get(0).(*sqs.DeleteMessageBatchOutput), args.Error(1)
}

func (m *MockSQS) GetQueueUrl(input *sqs.GetQueueUrlInput) (*sqs.GetQueueUrlOutput, error) {
	args := m.Called(input)
	return args.Get(0).(*sqs.GetQueueUrlOutput), args.Error(1)
}

func (m *MockSQS) ReceiveMessage(input *sqs.ReceiveMessageInput) (*sqs.ReceiveMessageOutput, error) {
	args := m.Called(input)
	return args.Get(0).(*sqs.ReceiveMessageOutput), args.Error(1)
}

type MockSQSService struct {
	Svc    SQSAPI
	URL    *string
	Logger func(format string, args ...interface{})
}
