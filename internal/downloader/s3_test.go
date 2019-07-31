package downloader

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"testing"

	"k8s.io/helm/pkg/chartutil"

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
	var m *mockS3
	dl := New("foo", m)
	assert.NotNil(t, dl)
	assert.Equal(t, "foo", dl.Bucket)
	assert.Equal(t, m, dl.d)

}

func TestDownloadSuccess(t *testing.T) {
	mockD := &mockS3{}
	fakeCtx := context.Background()
	dl := New("foo", mockD)

	chartBytes, err := ioutil.ReadFile("../../testdata/foo-0.1.0.tgz")
	assert.NoError(t, err)

	mockD.On("DownloadWithContext", fakeCtx, mock.AnythingOfType("*aws.WriteAtBuffer"), &s3.GetObjectInput{
		Bucket: aws.String(dl.Bucket),
		Key:    aws.String("/some-chart"),
	}).Return(0, nil).Run(func(args mock.Arguments) {
		w := args.Get(1).(*aws.WriteAtBuffer)
		w.WriteAt(chartBytes, 0)
	})

	outBytes, err := dl.download(fakeCtx, "/some-chart")

	assert.NoError(t, err)
	assert.NotNil(t, outBytes)
	assert.Equal(t, chartBytes, outBytes)
	mockD.AssertExpectations(t)
}

func TestDownloadFailure(t *testing.T) {
	mockD := &mockS3{}
	fakeCtx := context.Background()
	dl := New("foo", mockD)

	mockD.On("DownloadWithContext", fakeCtx, mock.AnythingOfType("*aws.WriteAtBuffer"), &s3.GetObjectInput{
		Bucket: aws.String(dl.Bucket),
		Key:    aws.String("/some-chart"),
	}).Return(0, fmt.Errorf("some download error"))

	_, err := dl.download(fakeCtx, "/some-chart")

	assert.Error(t, err)
	mockD.AssertExpectations(t)
}

func TestChartDownloadSuccess(t *testing.T) {
	mockD := &mockS3{}
	fakeCtx := context.Background()
	dl := New("foo", mockD)

	chartBytes, err := ioutil.ReadFile("../../testdata/foo-0.1.0.tgz")
	assert.NoError(t, err)

	expectedChart, err := chartutil.LoadArchive(bytes.NewBuffer(chartBytes))
	assert.NoError(t, err)

	mockD.On("DownloadWithContext", fakeCtx, mock.AnythingOfType("*aws.WriteAtBuffer"), &s3.GetObjectInput{
		Bucket: aws.String(dl.Bucket),
		Key:    aws.String("/some-chart"),
	}).Return(0, nil).Run(func(args mock.Arguments) {
		w := args.Get(1).(*aws.WriteAtBuffer)
		w.WriteAt(chartBytes, 0)
	})

	outChart, err := dl.DownloadChart(fakeCtx, "/some-chart")
	assert.NoError(t, err)

	assert.Equal(t, expectedChart, outChart)
	mockD.AssertExpectations(t)
}

func TestChartDownloadFailure(t *testing.T) {
	mockD := &mockS3{}
	fakeCtx := context.Background()
	dl := New("foo", mockD)

	mockD.On("DownloadWithContext", fakeCtx, mock.AnythingOfType("*aws.WriteAtBuffer"), &s3.GetObjectInput{
		Bucket: aws.String(dl.Bucket),
		Key:    aws.String("/some-chart"),
	}).Return(0, fmt.Errorf("some download error"))

	outChart, err := dl.DownloadChart(fakeCtx, "/some-chart")
	assert.Error(t, err)

	assert.Nil(t, outChart)
	mockD.AssertExpectations(t)

}

func TestInvalidChartBytes(t *testing.T) {
	mockD := &mockS3{}
	fakeCtx := context.Background()
	dl := New("foo", mockD)

	chartBytes := []byte("some bad bytes")

	mockD.On("DownloadWithContext", fakeCtx, mock.AnythingOfType("*aws.WriteAtBuffer"), &s3.GetObjectInput{
		Bucket: aws.String(dl.Bucket),
		Key:    aws.String("/some-chart"),
	}).Return(0, nil).Run(func(args mock.Arguments) {
		w := args.Get(1).(*aws.WriteAtBuffer)
		w.WriteAt(chartBytes, 0)
	})

	outChart, err := dl.DownloadChart(fakeCtx, "/some-chart")
	assert.Error(t, err)
	assert.Nil(t, outChart)
	mockD.AssertExpectations(t)
}

func TestMakeS3Key(t *testing.T) {
	type testCase struct {
		input    string
		expected string
	}

	tests := []testCase{
		{"some-chart", "/some-chart"},
		{"/some-chart", "/some-chart"},
	}

	for _, test := range tests {
		assert.Equal(t, test.expected, makeS3Key(test.input))
	}
}
