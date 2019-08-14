package chartdownloader

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	"github.com/aws/aws-sdk-go/aws/client"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"k8s.io/helm/pkg/proto/hapi/chart"
)

type ChartDownloader interface {
	Download(context.Context, string) (*chart.Chart, error)
}

// TODO: github, GCS support?
type Providers interface {
	S3() client.ConfigProvider
}

type ProviderFuncs struct {
	S3Func func() client.ConfigProvider
}

func (p ProviderFuncs) S3() client.ConfigProvider {
	if p.S3Func == nil {
		return nil
	}
	return p.S3Func()
}

func New(rawRepoURL string, p Providers) (ChartDownloader, error) {
	repoURL, err := url.Parse(rawRepoURL)
	if err != nil {
		return nil, err
	}

	switch repoURL.Scheme {
	case "s3":
		parts := strings.Split(repoURL.Path, "/")
		bucket := parts[len(parts)-1]

		s3 := p.S3()
		if s3 == nil {
			return nil, fmt.Errorf("no S3 provider for helm repository %s", rawRepoURL)
		}

		return newS3Downloader(bucket, s3manager.NewDownloader(s3)), nil
	default:
		return nil, fmt.Errorf("unsupported helm repository %s", rawRepoURL)
	}
}
