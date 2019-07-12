package service

import (
	"context"
	"os"
	"testing"
	"time"

	"ship-it/internal/api/integrations/github"
	"ship-it/internal/api/models"

	"github.com/go-kit/kit/log"
	"github.com/stretchr/testify/assert"
)

type mockK8sClient struct {
	releases []models.Release
}

func (k *mockK8sClient) ListAll(namespace string) ([]models.Release, error) {
	return k.releases, nil
}

func newMockK8sClient(name string, time time.Time) *mockK8sClient {
	r := models.Release{
		Name:    name,
		Created: time,
	}

	return &mockK8sClient{
		releases: []models.Release{r},
	}
}

type mockGitClient struct {
	Org    string
	Checks github.ChecksService
}

func newMockGitClient(org string) *mockGitClient {
	return &mockGitClient{
		Org: org,
	}
}

func (m *mockGitClient) GetTravisCIBuildURLForRef(ctx context.Context, repo string, ref string) (string, error) {
	return "www.travisci.com", nil
}

func TestListReleasesPopulatesFromCustomResource(t *testing.T) {
	name := "mock-resource"
	currentTime := time.Now()

	mockK8s := newMockK8sClient(name, currentTime)
	mockGit := newMockGitClient("wattpad")
	svc := New(mockK8s, mockGit, log.NewJSONLogger(log.NewSyncWriter(os.Stdout)))

	releases, err := svc.ListReleases(context.Background())

	assert.Nil(t, err)
	assert.Len(t, releases, 1)
	assert.Equal(t, name, releases[0].Name)
	assert.Equal(t, currentTime, releases[0].Created)
}

func TestTravisFieldIsPopulated(t *testing.T) {
	mockK8s := newMockK8sClient("word-counts", time.Now())
	mockGit := newMockGitClient("wattpad")
	svc := New(mockK8s, mockGit, log.NewJSONLogger(log.NewSyncWriter(os.Stdout)))

	ctx := context.Background()
	releases, err := svc.ListReleases(ctx)

	assert.Nil(t, err)
	assert.Len(t, releases, 1)

	assert.Equal(t, releases[0].Build.Travis, "www.travisci.com")
}
