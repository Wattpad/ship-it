package chartdownloader

import (
	"bytes"
	"context"
	"io"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/pkg/errors"
	"k8s.io/helm/pkg/chartutil"
	"k8s.io/helm/pkg/proto/hapi/chart"
)

type downloaderAPI interface {
	DownloadWithContext(ctx aws.Context, w io.WriterAt, input *s3.GetObjectInput, options ...func(*s3manager.Downloader)) (int64, error)
}

type S3Downloader struct {
	Bucket string
	d      downloaderAPI
}

func NewS3Downloader(bucketName string, dl downloaderAPI) S3Downloader {
	return S3Downloader{
		Bucket: bucketName,
		d:      dl,
	}
}

func (dl S3Downloader) download(ctx context.Context, key string) ([]byte, error) {
	buf := aws.NewWriteAtBuffer(nil)

	_, err := dl.d.DownloadWithContext(ctx, buf, &s3.GetObjectInput{
		Bucket: aws.String(dl.Bucket),
		Key:    aws.String(key),
	})

	return buf.Bytes(), err
}

func makeS3Key(chartName string) string {
	return "/" + strings.TrimPrefix(chartName, "/")
}

func (dl S3Downloader) Download(ctx context.Context, chartName string) (*chart.Chart, error) {
	chartBytes, err := dl.download(ctx, makeS3Key(chartName))
	if err != nil {
		return nil, errors.Wrapf(err, "failed to download chart %s", chartName)
	}

	return chartutil.LoadArchive(bytes.NewBuffer(chartBytes))
}
