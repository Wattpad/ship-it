package github

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"k8s.io/helm/pkg/helm"
	"k8s.io/helm/pkg/proto/hapi/chart"
	rls "k8s.io/helm/pkg/proto/hapi/services"
	helmerrors "k8s.io/helm/pkg/storage/errors"
)

type mockHelmClient struct {
	mock.Mock
}

func (m *mockHelmClient) InstallReleaseFromChart(chart *chart.Chart, namespace string, opts ...helm.InstallOption) (*rls.InstallReleaseResponse, error) {
	args := m.Called(chart, namespace, opts)

	var ret0 *rls.InstallReleaseResponse
	if args0 := args.Get(0); args0 != nil {
		ret0 = args0.(*rls.InstallReleaseResponse)
	}

	return ret0, args.Error(1)
}

func (m *mockHelmClient) UpdateReleaseFromChart(rlsName string, chart *chart.Chart, opts ...helm.UpdateOption) (*rls.UpdateReleaseResponse, error) {
	args := m.Called(rlsName, chart, opts)

	var ret0 *rls.UpdateReleaseResponse
	if args0 := args.Get(0); args0 != nil {
		ret0 = args0.(*rls.UpdateReleaseResponse)
	}

	return ret0, args.Error(1)
}

func TestUpdateOrInstallReleaseFound(t *testing.T) {
	testChart := new(chart.Chart)
	testRelease := "release"

	var m mockHelmClient
	m.On("UpdateReleaseFromChart", testRelease, testChart, mock.Anything).Return(&rls.UpdateReleaseResponse{}, nil)

	release := newReleaseSyncer(&m, "", testRelease, 0)

	err := release.UpdateOrInstallFromChart(context.Background(), testChart)
	assert.NoError(t, err)

	m.AssertExpectations(t)
	m.AssertNotCalled(t, "InstallReleaseFromChart")
}

func TestUpdateOrInstallReleaseNotFound(t *testing.T) {
	testChart := new(chart.Chart)
	testNamespace, testRelease := "namespace", "release"

	var m mockHelmClient
	m.On("UpdateReleaseFromChart", testRelease, testChart, mock.Anything).Return(nil, helmerrors.ErrReleaseNotFound(testRelease))
	m.On("InstallReleaseFromChart", testChart, testNamespace, mock.Anything).Return(&rls.InstallReleaseResponse{}, nil)

	release := newReleaseSyncer(&m, testNamespace, testRelease, 0)

	err := release.UpdateOrInstallFromChart(context.Background(), testChart)
	assert.NoError(t, err)

	m.AssertExpectations(t)
}
