package downloader

import (
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

func (m *mockS3) Download(w io.WriterAt, input *s3.GetObjectInput, options ...func(*s3manager.Downloader)) (n int64, err error) {
	args := m.Called(w, input)
	w.WriteAt([]byte("this is some bytes"), 0)
	return int64(args.Int(0)), args.Error(1)
}

func (m *mockS3) DownloadWithContext(ctx aws.Context, w io.WriterAt, input *s3.GetObjectInput, options ...func(*s3manager.Downloader)) (n int64, err error) {
	args := m.Called(ctx, w, input)
	return int64(args.Int(0)), args.Error(1)
}

func (m *mockS3) DownloadWithIterator(ctx aws.Context, iter s3manager.BatchDownloadIterator, opts ...func(*s3manager.Downloader)) error {
	args := m.Called(ctx, iter)
	return args.Error(0)
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
	dl, err := New("foo", mockD)
	assert.NoError(t, err)

	mockD.On("Download", mock.AnythingOfType("*aws.WriteAtBuffer"), &s3.GetObjectInput{
		Bucket: aws.String(dl.Bucket),
		Key:    aws.String("/some-chart"),
	}).Return(0, nil)

	outBytes, err := dl.Download("/some-chart")

	assert.NoError(t, err)
	assert.NotNil(t, outBytes)
	mockD.AssertExpectations(t)
}

func TestDownloadFailure(t *testing.T) {
	mockD := &mockS3{}
	dl, err := New("foo", mockD)
	assert.NoError(t, err)

	mockD.On("Download", mock.AnythingOfType("*aws.WriteAtBuffer"), &s3.GetObjectInput{
		Bucket: aws.String(dl.Bucket),
		Key:    aws.String("/some-chart"),
	}).Return(0, fmt.Errorf("some download error"))

	outBytes, err := dl.Download("/some-chart")

	assert.Error(t, err)
	assert.Nil(t, outBytes)
	mockD.AssertExpectations(t)
}
