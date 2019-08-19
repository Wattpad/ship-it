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

func (m *mockS3Downloader) Download(ctx context.Context, chartURL string) (*chart.Chart, error) {
	args := m.Called(ctx, chartURL)

	var ret0 *chart.Chart
	if args0 := args.Get(0); args0 != nil {
		ret0 = args0.(*chart.Chart)
	}

	return ret0, args.Error(1)
}

func TestFactoryS3Provider(t *testing.T) {
	bucket := "helm-charts"
	repoURL := fmt.Sprintf("s3://charts.wattpadhq.com/%s", bucket)

	mockS3 := new(mockS3Downloader)
	mockS3.On("Download", mock.Anything, repoURL).Return(&chart.Chart{}, nil)

	dl := New(map[string]Interface{
		"s3": mockS3,
	})

	_, err := dl.Download(context.Background(), repoURL)
	require.NoError(t, err)
}

func TestFactoryUnsupportedProvider(t *testing.T) {
	repoURL := "git://github.com/Wattpad/foo"
	_, err := New(nil).Download(context.Background(), repoURL)
	assert.Error(t, err)
}
