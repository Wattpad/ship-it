package downloader

import (
	"context"
	"fmt"
	"io"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockS3 struct {
	mock.Mock
}

func (m *mockS3) DownloadWithContext(ctx aws.Context, w io.WriterAt, input *s3.GetObjectInput, options ...func(*s3manager.Downloader)) (n int64, err error) {
	args := m.Called(ctx, w, input)
	return int64(args.Int(0)), args.Error(1)
}

func TestNewDownloader(t *testing.T) {
	t.Run("valid downloader", func(t *testing.T) {
		var m *mockS3
		dl, err := New("foo", m)
		assert.NoError(t, err)
		assert.NotNil(t, dl)
	})

	t.Run("invalid downloader", func(t *testing.T) {
		dl, err := New("foo", nil)
		assert.Error(t, err)
		assert.Nil(t, dl)
	})
}

func TestDownloadSuccess(t *testing.T) {
	mockD := &mockS3{}
	fakeCtx := context.Background()
	dl, err := New("foo", mockD)
	assert.NoError(t, err)

	mockD.On("DownloadWithContext", fakeCtx, mock.AnythingOfType("*aws.WriteAtBuffer"), &s3.GetObjectInput{
		Bucket: aws.String(dl.Bucket),
		Key:    aws.String("/some-chart"),
	}).Return(0, nil)

	outBytes, err := dl.Download("/some-chart", fakeCtx)

	assert.NoError(t, err)
	assert.NotNil(t, outBytes)
	mockD.AssertExpectations(t)
}

func TestDownloadFailure(t *testing.T) {
	mockD := &mockS3{}
	fakeCtx := context.Background()
	dl, err := New("foo", mockD)
	assert.NoError(t, err)

	mockD.On("DownloadWithContext", fakeCtx, mock.AnythingOfType("*aws.WriteAtBuffer"), &s3.GetObjectInput{
		Bucket: aws.String(dl.Bucket),
		Key:    aws.String("/some-chart"),
	}).Return(0, fmt.Errorf("some download error"))

	outBytes, err := dl.Download("/some-chart", fakeCtx)

	assert.Error(t, err)
	assert.Nil(t, outBytes)
	mockD.AssertExpectations(t)
}
