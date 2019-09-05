package controllers

import (
	shipitv1beta1 "ship-it-operator/api/v1beta1"

	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/tools/record"
	"k8s.io/helm/pkg/helm"
	"k8s.io/helm/pkg/proto/hapi/chart"
	"k8s.io/helm/pkg/proto/hapi/release"
)

// ReleaseManager performs release lifecycle operations using a helm client, and
// modifies HelmRelease fields like the Finalizers and Status. It broadcasts
// kubernetes events whenever the release's status condition changes.
type ReleaseManager struct {
	helm     HelmClient
	recorder record.EventRecorder
}

func (m *ReleaseManager) updateCondition(rls *shipitv1beta1.HelmRelease, cond shipitv1beta1.HelmReleaseCondition) *shipitv1beta1.HelmRelease {
	rls.Status.SetCondition(cond)
	m.recorder.Event(rls, v1.EventTypeNormal, cond.Type, cond.Message)

	return rls
}

func (m *ReleaseManager) Install(rls *shipitv1beta1.HelmRelease, chart *chart.Chart, namespace string) (*shipitv1beta1.HelmRelease, error) {
	if _, err := m.helm.InstallReleaseFromChart(
		chart,
		namespace,
		helm.InstallReuseName(true),
		helm.ReleaseName(rls.Spec.ReleaseName),
		helm.ValueOverrides(rls.Spec.Values.Raw),
	); err != nil {
		return nil, err
	}

	cond := shipitv1beta1.HelmReleaseCondition{
		Type:    release.Status_PENDING_INSTALL.String(),
		Message: "Installing release",
	}

	return m.updateCondition(rls, cond), nil
}

func (m *ReleaseManager) Delete(rls *shipitv1beta1.HelmRelease) (*shipitv1beta1.HelmRelease, error) {
	if _, err := m.helm.DeleteRelease(rls.Spec.ReleaseName); err != nil {
		return nil, err
	}

	cond := shipitv1beta1.HelmReleaseCondition{
		Type:    release.Status_DELETING.String(),
		Message: "Deleting release",
	}

	return m.updateCondition(rls, cond), nil
}

func (m *ReleaseManager) Upgrade(rls *shipitv1beta1.HelmRelease, chart *chart.Chart) (*shipitv1beta1.HelmRelease, error) {
	if _, err := m.helm.UpdateReleaseFromChart(
		rls.Spec.ReleaseName,
		chart,
		helm.UpdateValueOverrides(rls.Spec.Values.Raw),
	); err != nil {
		return nil, err
	}

	cond := shipitv1beta1.HelmReleaseCondition{
		Type:    release.Status_PENDING_UPGRADE.String(),
		Message: "Upgrading release",
	}

	return m.updateCondition(rls, cond), nil
}

func (m *ReleaseManager) Rollback(rls *shipitv1beta1.HelmRelease) (*shipitv1beta1.HelmRelease, error) {
	if _, err := m.helm.RollbackRelease(rls.Spec.ReleaseName); err != nil {
		return nil, err
	}

	cond := shipitv1beta1.HelmReleaseCondition{
		Type:    release.Status_PENDING_ROLLBACK.String(),
		Message: "Rolling back release",
	}

	return m.updateCondition(rls, cond), nil
}

func (m *ReleaseManager) Deployed(rls *shipitv1beta1.HelmRelease) *shipitv1beta1.HelmRelease {
	oldCondition := rls.Status.GetCondition()

	var reason shipitv1beta1.HelmReleaseStatusReason

	switch oldCondition.Type {
	case release.Status_PENDING_INSTALL.String():
		reason = shipitv1beta1.ReasonInstallSuccess
	case release.Status_PENDING_UPGRADE.String():
		reason = shipitv1beta1.ReasonUpdateSuccess
	case release.Status_PENDING_ROLLBACK.String():
		reason = shipitv1beta1.ReasonRollbackSuccess
	case release.Status_DEPLOYED.String():
		reason = oldCondition.Reason
	default:
		reason = shipitv1beta1.ReasonUnknown
	}

	cond := shipitv1beta1.HelmReleaseCondition{
		Type:    release.Status_DEPLOYED.String(),
		Reason:  reason,
		Message: "Release deployed",
	}

	return m.updateCondition(rls, cond)
}

func (m *ReleaseManager) Failed(rls *shipitv1beta1.HelmRelease) *shipitv1beta1.HelmRelease {
	oldCondition := rls.Status.GetCondition()

	var reason shipitv1beta1.HelmReleaseStatusReason

	switch oldCondition.Type {
	case release.Status_PENDING_INSTALL.String():
		reason = shipitv1beta1.ReasonInstallError
	case release.Status_PENDING_UPGRADE.String():
		reason = shipitv1beta1.ReasonUpdateError
	case release.Status_PENDING_ROLLBACK.String():
		reason = shipitv1beta1.ReasonRollbackError
	case release.Status_FAILED.String():
		reason = oldCondition.Reason
	default:
		reason = shipitv1beta1.ReasonUnknown
	}

	cond := shipitv1beta1.HelmReleaseCondition{
		Type:    release.Status_FAILED.String(),
		Reason:  reason,
		Message: "Release failed",
	}

	return m.updateCondition(rls, cond)
}
