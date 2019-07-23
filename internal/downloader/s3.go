package downloader

import (
	"bytes"
	"context"
	"fmt"
	"io"

	"k8s.io/helm/pkg/proto/hapi/chart"
	"k8s.io/helm/pkg/chartutil"

	"github.com/pkg/errors"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
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

func (dl Downloader) Download(ctx context.Context, key string) ([]byte, error) {
	buff := aws.NewWriteAtBuffer([]byte{})

	_, err := dl.d.DownloadWithContext(ctx, buff, &s3.GetObjectInput{
		Bucket: aws.String(dl.Bucket),
		Key:    aws.String(key),
	})

	if err != nil {
		return nil, errors.Wrap(err, "chart download failed")
	}

	return buff.Bytes(), nil
}

func makeS3Key(chartName string) string {
	if fmt.Sprintf("%c", chartName[0]) != "/" {
		return "/" + chartName
	}
	return chartName
}

func (dl Downloader) DownloadChart(ctx context.Context, chartName string) (*chart.Chart, error) {
	chartBytes, err := dl.Download(ctx, chartName)
	if err != nil {
		return nil, err
	}

	chartObj, err := chartutil.LoadArchive(bytes.NewBuffer(chartBytes))
	if err != nil {
		return nil, err
	}

	return chartObj, nil
}
