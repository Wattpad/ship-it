package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"ship-it/internal/api/models"

	"github.com/stretchr/testify/assert"
)

type mockHelmClient struct {
	resources map[string]string
}

func (m *mockHelmClient) Get(name string) (string, error) {
	if res, ok := m.resources[name]; ok {
		return res, nil
	}
	return "", errors.New("release not found")
}

type mockK8sClient struct {
	releases []models.Release
}

func (k *mockK8sClient) List(ctx context.Context, namespace string) ([]models.Release, error) {
	return k.releases, nil
}

func (k *mockK8sClient) Get(ctx context.Context, namespace, name string) (*models.Release, error) {
	for _, r := range k.releases {
		if r.Name == name {
			return &r, nil
		}
	}
	return nil, errors.New("release not found")
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

func TestGetAndListReleases(t *testing.T) {
	name := "releaseName"
	currentTime := time.Now()

	mock := newMockK8sClient(name, currentTime)
	svc := New(mock, nil)

	releases, err := svc.ListReleases(context.Background())
	if !assert.NoError(t, err) {
		assert.Len(t, releases, 1)
		assert.Equal(t, name, releases[0].Name)
		assert.Equal(t, currentTime, releases[0].Created)
	}

	release, err := svc.GetRelease(context.Background(), name)
	if !assert.NoError(t, err) {
		assert.Equal(t, name, release.Name)
		assert.Equal(t, currentTime, release.Created)
	}
}

func TestGetReleaseResources(t *testing.T) {
	name := "releaseName"
	resources := "foobarbaz"

	mock := mockHelmClient{
		resources: map[string]string{
			name: resources,
		},
	}

	svc := New(nil, &mock)

	res, err := svc.GetReleaseResources(context.Background(), name)
	if assert.NoError(t, err) {
		assert.Equal(t, res.Name, name)
		assert.Equal(t, res.Resources, resources)
	}
}
