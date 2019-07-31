package downloader

import (
	"bytes"
	"context"
	"io"

	"k8s.io/helm/pkg/chartutil"
	"k8s.io/helm/pkg/proto/hapi/chart"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/pkg/errors"
)

type downloaderAPI interface {
	DownloadWithContext(ctx aws.Context, w io.WriterAt, input *s3.GetObjectInput, options ...func(*s3manager.Downloader)) (n int64, err error)
}

type Downloader struct {
	Bucket string
	d      downloaderAPI
}

func New(bucketName string, dl downloaderAPI) Downloader {
	return Downloader{
		Bucket: bucketName,
		d:      dl,
	}
}

func (dl Downloader) download(ctx context.Context, key string) ([]byte, error) {
	buff := aws.NewWriteAtBuffer([]byte{})

	_, err := dl.d.DownloadWithContext(ctx, buff, &s3.GetObjectInput{
		Bucket: aws.String(dl.Bucket),
		Key:    aws.String(key),
	})

	return buff.Bytes(), err
}

func makeS3Key(chartName string) string {
	if chartName[0] != '/' {
		return "/" + chartName
	}
	return chartName
}

func (dl Downloader) DownloadChart(ctx context.Context, chartName string) (*chart.Chart, error) {
	chartBytes, err := dl.download(ctx, makeS3Key(chartName))
	if err != nil {
		return nil, errors.Wrap(err, "chart download failed")
	}

	return chartutil.LoadArchive(bytes.NewBuffer(chartBytes))
}
