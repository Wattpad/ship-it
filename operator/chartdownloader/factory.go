package chartdownloader

import (
	"context"
	"fmt"
	"net/url"

	"github.com/pkg/errors"
	"k8s.io/helm/pkg/proto/hapi/chart"
)

type ChartDownloader interface {
	Download(context.Context, string, string) (*chart.Chart, error)
}

type factory struct {
	downloaders map[string]ChartDownloader
}

func New(downloaders map[string]ChartDownloader) ChartDownloader {
	return &factory{
		downloaders: downloaders,
	}
}

func (f *factory) Download(ctx context.Context, rawChartURL string, version string) (*chart.Chart, error) {
	repoURL, err := url.Parse(rawChartURL)
	if err != nil {
		return nil, errors.Wrapf(err, "invalid chart URL %s", rawChartURL)
	}

	if dl, ok := f.downloaders[repoURL.Scheme]; ok {
		return dl.Download(ctx, rawChartURL, version)
	}

	return nil, fmt.Errorf("unsupported chart transport protocol %s", rawChartURL)
}
