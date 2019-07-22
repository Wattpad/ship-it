package downloader

import (
	"fmt"
	"io"

	"github.com/pkg/errors"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

type downloaderAPI interface {
	Download(w io.WriterAt, input *s3.GetObjectInput, options ...func(*s3manager.Downloader)) (n int64, err error)
	DownloadWithContext(ctx aws.Context, w io.WriterAt, input *s3.GetObjectInput, options ...func(*s3manager.Downloader)) (n int64, err error)
	DownloadWithIterator(ctx aws.Context, iter s3manager.BatchDownloadIterator, opts ...func(*s3manager.Downloader)) error
}

type Downloader struct {
	Bucket string
	d downloaderAPI
}

func New(bucketName string, dl downloaderAPI) (*Downloader, error) {
	if dl == nil {
		return nil, fmt.Errorf("received nil downloader")
	}

	return &Downloader{
		Bucket: bucketName,
		d: dl,
	}, nil
}

func (dl *Downloader) Download(key string) ([]byte, error) {
	buff := aws.NewWriteAtBuffer([]byte{})

	_, err := dl.d.Download(buff, &s3.GetObjectInput{
		Bucket: aws.String(dl.Bucket),
		Key: aws.String(key),
	})

	if err != nil {
		return nil, errors.Wrap(err, "chart download failed")
	}

	return buff.Bytes(), nil
}
