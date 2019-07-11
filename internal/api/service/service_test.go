package service

import (
	"context"
	"testing"
	"time"

	"ship-it/internal/api/integrations/github"
	"ship-it/internal/api/models"

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

func TestListReleasesPopulatesFromConfigMap(t *testing.T) {
	name := "configmap-1"
	currentTime := time.Now()

	mock := newMockK8sClient(name, currentTime)
	mockGit := github.New(context.Background(), "Wattpad", "fake")
	svc := New(mock, mockGit, "miranda", "master")

	releases, err := svc.ListReleases(context.Background())

	assert.Nil(t, err)
	assert.Len(t, releases, 1)
	assert.Equal(t, name, releases[0].Name)
	assert.Equal(t, currentTime, releases[0].Created)
}
