package chartdownloader

import (
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go/aws/client"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFactoryS3Provider(t *testing.T) {
	bucket := "helm-charts"
	repoURL := fmt.Sprintf("s3://wattpad.amazonaws.com/%s", bucket)

	dl, err := New(repoURL, &ProviderFuncs{
		S3Func: func() client.ConfigProvider {
			return session.Must(session.NewSession())
		},
	})

	require.NoError(t, err)

	s3, ok := dl.(*S3Downloader)
	assert.True(t, ok)
	assert.Equal(t, s3.Bucket, bucket)
}

func TestFactoryUnsupportedProvider(t *testing.T) {
	repoURL := "git://github.com/Wattpad/foo"
	_, err := New(repoURL, nil)
	assert.Error(t, err)
}
