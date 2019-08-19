package chartdownloader

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
	"github.com/stretchr/testify/require"
)

type mockS3 struct {
	mock.Mock
}

func (m *mockS3) DownloadWithContext(ctx aws.Context, w io.WriterAt, input *s3.GetObjectInput, options ...func(*s3manager.Downloader)) (n int64, err error) {
	args := m.Called(ctx, w, input)
	return int64(args.Int(0)), args.Error(1)
}

func TestChartDownloadSuccess(t *testing.T) {
	ctx := context.Background()

	testBucket := "wattpad.amazonaws.com"
	testObject := "some-chart"

	var mockD mockS3
	dl := NewS3Downloader(&mockD)

	chartBytes, err := ioutil.ReadFile("../../testdata/foo-0.1.0.tgz")
	require.NoError(t, err)

	expectedChart, err := chartutil.LoadArchive(bytes.NewBuffer(chartBytes))
	require.NoError(t, err)

	mockD.On("DownloadWithContext", ctx, mock.AnythingOfType("*aws.WriteAtBuffer"), &s3.GetObjectInput{
		Bucket: aws.String(testBucket),
		Key:    aws.String(testObject),
	}).Return(0, nil).Run(func(args mock.Arguments) {
		w := args.Get(1).(*aws.WriteAtBuffer)
		w.WriteAt(chartBytes, 0)
	})

	outChart, err := dl.Download(ctx, fmt.Sprintf("s3://%s/%s", testBucket, testObject))
	require.NoError(t, err)

	assert.Equal(t, expectedChart, outChart)
	mockD.AssertExpectations(t)
}

func TestChartDownloadFailure(t *testing.T) {
	ctx := context.Background()

	testBucket := "wattpad.amazonaws.com"
	testObject := "some-chart"

	var mockD mockS3
	dl := NewS3Downloader(&mockD)

	mockD.On("DownloadWithContext", ctx, mock.AnythingOfType("*aws.WriteAtBuffer"), &s3.GetObjectInput{
		Bucket: aws.String(testBucket),
		Key:    aws.String(testObject),
	}).Return(0, fmt.Errorf("some download error"))

	_, err := dl.Download(ctx, fmt.Sprintf("s3://%s/%s", testBucket, testObject))
	assert.Error(t, err)
	mockD.AssertExpectations(t)
}

func TestInvalidChartBytes(t *testing.T) {
	ctx := context.Background()

	testBucket := "wattpad.amazonaws.com"
	testObject := "some-chart"

	var mockD mockS3
	dl := NewS3Downloader(&mockD)

	chartBytes := []byte("some bad bytes")

	mockD.On("DownloadWithContext", ctx, mock.AnythingOfType("*aws.WriteAtBuffer"), &s3.GetObjectInput{
		Bucket: aws.String(testBucket),
		Key:    aws.String(testObject),
	}).Return(0, nil).Run(func(args mock.Arguments) {
		w := args.Get(1).(*aws.WriteAtBuffer)
		w.WriteAt(chartBytes, 0)
	})

	_, err := dl.Download(ctx, fmt.Sprintf("s3://%s/%s", testBucket, testObject))
	assert.Error(t, err)
	mockD.AssertExpectations(t)
}

func TestParseBucketObject(t *testing.T) {
	type testCase struct {
		input  string
		bucket string
		object string
	}

	tests := []testCase{
		{
			input:  "s3://charts.wattpadhq.com/microservice",
			bucket: "charts.wattpadhq.com",
			object: "microservice",
		},
		{
			input:  "s3://charts.wattpadhq.com/foo/bar",
			bucket: "charts.wattpadhq.com",
			object: "foo/bar",
		},
	}

	for _, test := range tests {
		bucket, object, err := parseBucketObject(test.input)
		assert.NoError(t, err)
		assert.Equal(t, test.bucket, bucket)
		assert.Equal(t, test.object, object)
	}
}
