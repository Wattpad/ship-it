package chartdownloader

import (
	"bytes"
	"context"
	"io"
	"net/url"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/pkg/errors"
	"k8s.io/helm/pkg/chartutil"
	"k8s.io/helm/pkg/proto/hapi/chart"
)

type S3DownloadManager interface {
	DownloadWithContext(ctx aws.Context, w io.WriterAt, input *s3.GetObjectInput, options ...func(*s3manager.Downloader)) (int64, error)
}

type S3Downloader struct {
	manager S3DownloadManager
}

func NewS3Downloader(manager S3DownloadManager) *S3Downloader {
	return &S3Downloader{
		manager: manager,
	}
}

func (dl S3Downloader) download(ctx context.Context, bucket, key string) ([]byte, error) {
	buf := aws.NewWriteAtBuffer(nil)

	_, err := dl.manager.DownloadWithContext(ctx, buf, &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})

	return buf.Bytes(), err
}

func (dl S3Downloader) Download(ctx context.Context, chartURL string) (*chart.Chart, error) {
	bucket, object, err := parseBucketObject(chartURL)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to parse S3 bucket and object from URL %s", chartURL)
	}

	chartBytes, err := dl.download(ctx, bucket, object)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to download chart %s", chartURL)
	}

	return chartutil.LoadArchive(bytes.NewBuffer(chartBytes))
}

func parseBucketObject(rawChartURL string) (bucket string, object string, err error) {
	chartURL, err := url.Parse(rawChartURL)
	if err != nil {
		return "", "", err
	}

	return chartURL.Host, strings.TrimPrefix(chartURL.Path, "/"), nil
}
