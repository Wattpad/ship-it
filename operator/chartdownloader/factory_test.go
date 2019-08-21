package chartdownloader

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"k8s.io/helm/pkg/proto/hapi/chart"
)

type mockS3Downloader struct {
	mock.Mock
}

func (m *mockS3Downloader) Download(ctx context.Context, chartURL string, version string) (*chart.Chart, error) {
	args := m.Called(ctx, chartURL, version)

	var ret0 *chart.Chart
	if args0 := args.Get(0); args0 != nil {
		ret0 = args0.(*chart.Chart)
	}

	return ret0, args.Error(1)
}

func TestFactoryS3Provider(t *testing.T) {
	bucket := "helm-charts"
	chartPath := fmt.Sprintf("s3://charts.wattpadhq.com/%s", bucket)
	version := "0.0.0"

	mockS3 := new(mockS3Downloader)
	mockS3.On("Download", mock.Anything, chartPath, version).Return(&chart.Chart{}, nil)

	dl := New(map[string]ChartDownloader{
		"s3": mockS3,
	})

	_, err := dl.Download(context.Background(), chartPath, version)
	require.NoError(t, err)
}

func TestFactoryUnsupportedProvider(t *testing.T) {
	chartPath := "git://github.com/Wattpad/foo"
	version := "0.0.0"
	_, err := New(nil).Download(context.Background(), chartPath, version)
	assert.Error(t, err)
}
