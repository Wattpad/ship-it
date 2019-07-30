package service

import (
	"context"
	"fmt"
	"testing"
	"time"

	"ship-it/internal/api/integrations/github"
	"ship-it/internal/api/models"

	"github.com/go-kit/kit/log"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockK8sClient struct {
	releases []models.Release
}

func (k *mockK8sClient) ListAll(ctx context.Context, namespace string) ([]models.Release, error) {
	return k.releases, nil
}

func newMockK8sClient(name string, time time.Time, repo string) *mockK8sClient {
	r := models.Release{
		Name:    name,
		Created: time,
		Code: models.SourceCode{
			Github: repo,
			Ref:    "master",
		},
	}

	return &mockK8sClient{
		releases: []models.Release{r},
	}
}

type mockGitClient struct {
	mock.Mock
	Checks github.ChecksService
}

func (m *mockGitClient) GetTravisCIBuildURLForRef(ctx context.Context, repo string, ref string) (string, error) {
	args := m.Called(ctx, repo, ref)
	return args.String(0), args.Error(1)
}

func TestListReleasesPopulatesFromCustomResource(t *testing.T) {
	name := "mock-resource"
	currentTime := time.Now()

	mockK8s := newMockK8sClient(name, currentTime, "highlander")
	mockGit := new(mockGitClient)

	svc := New(mockK8s, mockGit, log.NewNopLogger())

	mockGit.On("GetTravisCIBuildURLForRef", context.Background(), "highlander", "master").Return("", nil)
	releases, err := svc.ListReleases(context.Background())

	assert.Nil(t, err)
	assert.Len(t, releases, 1)
	assert.Equal(t, name, releases[0].Name)
	assert.Equal(t, currentTime, releases[0].Created)
}

func TestTravisFieldIsPopulated(t *testing.T) {
	t.Run("valid input", func(t *testing.T) {
		mockK8s := newMockK8sClient("word-counts", time.Now(), "highlander")
		mockGit := new(mockGitClient)
		svc := New(mockK8s, mockGit, log.NewNopLogger())

		mockGit.On("GetTravisCIBuildURLForRef", context.Background(), "highlander", "master").Return("travisci.com/build/highlander/master", nil)
		releases, err := svc.ListReleases(context.Background())

		assert.Nil(t, err)
		assert.Len(t, releases, 1)

		assert.Equal(t, releases[0].Build.Travis, "travisci.com/build/highlander/master")
	})

	t.Run("invalid input", func(t *testing.T) {
		mockK8s := newMockK8sClient("word-counts", time.Now(), "miranda")
		mockGit := new(mockGitClient)
		svc := New(mockK8s, mockGit, log.NewNopLogger())

		mockGit.On("GetTravisCIBuildURLForRef", context.Background(), "miranda", "master").Return("", fmt.Errorf("unable to fetch travis url"))
		releases, _ := svc.ListReleases(context.Background())

		assert.Equal(t, releases[0].Build.Travis, "")
	})
}
