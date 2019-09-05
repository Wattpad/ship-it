/*

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controllers

import (
	"context"
	"fmt"
	"strings"
	"time"

	appsv1 "k8s.io/api/apps/v1"

	shipitv1beta1 "ship-it-operator/api/v1beta1"

	"github.com/go-logr/logr"
	"github.com/pkg/errors"
	apierrs "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/helm/pkg/helm"
	"k8s.io/helm/pkg/proto/hapi/chart"
	"k8s.io/helm/pkg/proto/hapi/release"
	hapi "k8s.io/helm/pkg/proto/hapi/services"
	helmerrors "k8s.io/helm/pkg/storage/errors"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
)

// The HelmReleaseFinalizer allows the controller to clean up the associated
// release before the HelmRelease resource is deleted.
const HelmReleaseFinalizer = "HelmReleaseFinalizer"

var errNotImplemented = errors.New("not implemented")

type ChartDownloader interface {
	Download(ctx context.Context, chart string, version string) (*chart.Chart, error)
}

type HelmClient interface {
	DeleteRelease(rlsName string, opts ...helm.DeleteOption) (*hapi.UninstallReleaseResponse, error)
	InstallReleaseFromChart(chart *chart.Chart, ns string, opts ...helm.InstallOption) (*hapi.InstallReleaseResponse, error)
	ReleaseStatus(rlsName string, opts ...helm.StatusOption) (*hapi.GetReleaseStatusResponse, error)
	RollbackRelease(rlsName string, opts ...helm.RollbackOption) (*hapi.RollbackReleaseResponse, error)
	UpdateReleaseFromChart(rlsName string, chart *chart.Chart, opts ...helm.UpdateOption) (*hapi.UpdateReleaseResponse, error)
}

// HelmReleaseReconciler reconciles a HelmRelease object
type HelmReleaseReconciler struct {
	client.Client
	reconcilerConfig

	Log logr.Logger

	downloader ChartDownloader
	helm       HelmClient
}

type ReconcilerOption func(*reconcilerConfig)

type reconcilerConfig struct {
	GracePeriod time.Duration
	Namespace   string
}

func Namespace(ns string) ReconcilerOption {
	return func(c *reconcilerConfig) {
		c.Namespace = ns
	}
}

func GracePeriod(d time.Duration) ReconcilerOption {
	return func(c *reconcilerConfig) {
		c.GracePeriod = d
	}
}

func NewHelmReleaseReconciler(l logr.Logger, client client.Client, helm HelmClient, d ChartDownloader, opts ...ReconcilerOption) *HelmReleaseReconciler {
	var cfg reconcilerConfig
	for _, opt := range opts {
		opt(&cfg)
	}

	return &HelmReleaseReconciler{
		Client: client,
		Log:    l.WithName("controllers").WithName("HelmRelease"),

		downloader:       d,
		helm:             helm,
		reconcilerConfig: cfg,
	}
}

// +kubebuilder:rbac:groups=shipit.wattpad.com,resources=helmreleases,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=shipit.wattpad.com,resources=helmreleases/status,verbs=get;update;patch

func (r *HelmReleaseReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	ctx := context.Background()
	log := r.Log.WithValues("name", req.NamespacedName)

	var helmRelease shipitv1beta1.HelmRelease

	if err := r.Get(ctx, req.NamespacedName, &helmRelease); err != nil {
		if apierrs.IsNotFound(err) {
			log.Info("HelmRelease doesn't exist")
			return ctrl.Result{}, nil
		}

		return ctrl.Result{}, err
	}

	if !helmRelease.Annotations().AutoDeploy() {
		log.Info("AutoDeploy is disabled for this HelmRelease")
		return ctrl.Result{}, nil
	}

	if !helmRelease.DeletionTimestamp.IsZero() {
		return r.delete(ctx, helmRelease)
	}

	if !contains(helmRelease.GetFinalizers(), HelmReleaseFinalizer) {
		// setting the finalizer does not change the release's
		// metadata.generation, so we have to requeue
		return ctrl.Result{Requeue: true}, r.setFinalizer(ctx, helmRelease)
	}

	return r.update(ctx, helmRelease)
}

func (r *HelmReleaseReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&shipitv1beta1.HelmRelease{}).
		WithEventFilter(predicate.GenerationChangedPredicate{}).
		Complete(r)
}

func (r *HelmReleaseReconciler) setFinalizer(ctx context.Context, rls shipitv1beta1.HelmRelease) error {
	finalizers := rls.GetFinalizers()
	rls.SetFinalizers(append(finalizers, HelmReleaseFinalizer))

	return r.Update(ctx, &rls)
}

func (r *HelmReleaseReconciler) clearFinalizer(ctx context.Context, rls shipitv1beta1.HelmRelease) error {
	var finalizers []string

	for _, f := range rls.GetFinalizers() {
		if f != HelmReleaseFinalizer {
			finalizers = append(finalizers, f)
		}
	}

	rls.SetFinalizers(finalizers)
	return r.Update(ctx, &rls)
}

func (r *HelmReleaseReconciler) delete(ctx context.Context, rls shipitv1beta1.HelmRelease) (ctrl.Result, error) {
	releaseName := rls.Spec.ReleaseName
	rlsStatus, err := r.helm.ReleaseStatus(releaseName)

	if err != nil {
		if isHelmReleaseNotFound(releaseName, err) {
			// this will only happen if a delete --purge is run
			return ctrl.Result{}, r.clearFinalizer(ctx, rls)
		}

		return ctrl.Result{}, err
	}

	info := rlsStatus.GetInfo()

	if info.Status.Code == release.Status_DELETING {
		return ctrl.Result{RequeueAfter: r.GracePeriod}, nil
	}

	if info.Status.Code == release.Status_DELETED {
		return ctrl.Result{}, r.clearFinalizer(ctx, rls)
	}

	_, err = r.helm.DeleteRelease(releaseName)

	if err != nil {
		return ctrl.Result{}, err
	}

	rls.Status.SetCondition(shipitv1beta1.HelmReleaseCondition{
		Type:    release.Status_DELETING.String(),
		Message: "Release is being deleted",
	})

	if err := r.Status().Update(ctx, &rls); err != nil {
		return ctrl.Result{}, err
	}

	return ctrl.Result{RequeueAfter: r.GracePeriod}, nil
}

func isHelmReleaseNotFound(name string, err error) bool {
	// dynamic errors can't be directly compared for equality. We use the
	// error message string, though there's no guarantee it won't change.
	return strings.Contains(err.Error(), helmerrors.ErrReleaseNotFound(name).Error())
}

// reasonForFailure determines the reason for a release's state transition to
// FAILED, given the previous release state. If the old state was FAILED, the
// original reason is re-used.
func reasonForFailure(oldCondition shipitv1beta1.HelmReleaseCondition) shipitv1beta1.HelmReleaseStatusReason {
	switch oldCondition.Type {
	case release.Status_DELETING.String():
		return shipitv1beta1.ReasonDeleteError
	case release.Status_PENDING_INSTALL.String():
		return shipitv1beta1.ReasonInstallError
	case release.Status_PENDING_UPGRADE.String():
		return shipitv1beta1.ReasonUpdateError
	case release.Status_PENDING_ROLLBACK.String():
		return shipitv1beta1.ReasonRollbackError
	case release.Status_FAILED.String():
		return oldCondition.Reason
	default:
		return shipitv1beta1.ReasonUnknown
	}
}

func (r *HelmReleaseReconciler) update(ctx context.Context, rls shipitv1beta1.HelmRelease) (ctrl.Result, error) {
	releaseName := rls.Spec.ReleaseName
	oldCondition := rls.Status.GetCondition()

	resp, err := r.helm.ReleaseStatus(releaseName)
	if err != nil {
		if isHelmReleaseNotFound(releaseName, err) {
			return r.install(ctx, rls)
		}
		return ctrl.Result{}, errors.Wrapf(err, "failed to get release status for %s", releaseName)
	}

	releaseStatus := resp.GetInfo().GetStatus()
	releaseStatusCode := releaseStatus.GetCode()

	switch releaseStatusCode {
	case release.Status_DELETING, release.Status_PENDING_INSTALL, release.Status_PENDING_ROLLBACK, release.Status_PENDING_UPGRADE:
		// if the release is still in transition, requeue until it settles
		return ctrl.Result{RequeueAfter: r.GracePeriod}, nil
	case release.Status_DEPLOYED:
		if oldCondition.Type == release.Status_DEPLOYED.String() {
			return r.upgrade(ctx, rls)
		}

		var reason shipitv1beta1.HelmReleaseStatusReason

		switch oldCondition.Type {
		case release.Status_PENDING_INSTALL.String():
			reason = shipitv1beta1.ReasonInstallSuccess
		case release.Status_PENDING_UPGRADE.String():
			reason = shipitv1beta1.ReasonUpdateSuccess
		case release.Status_PENDING_ROLLBACK.String():
			reason = shipitv1beta1.ReasonRollbackSuccess
		}

		r.Log.Info("Release is deployed, but we need to check if its fully rolled out")

		// Check all Deployments have rolled out
		// get Deployments by label instance=releaseName
		// check status.conditions for type: Available, status: False
		// if all pods aren't rolled out Requeue until X minutes has passed

		var deploymentList appsv1.DeploymentList

		err = r.List(
			ctx,
			&deploymentList,
			client.InNamespace(rls.ObjectMeta.Namespace),
			client.MatchingLabels(map[string]string{
				"instance": releaseName, // this makes an assumption
			}),
		)

		if err != nil {
			return ctrl.Result{}, err
		}

		// TODO this logic can live in its own func
		deployments := deploymentList.Items

		for _, deployment := range deployments {
			conditions := deployment.Status.Conditions

			var available bool
			var progressing bool

			for _, condition := range conditions {
				if condition.Type == "Available" {
					available = condition.Status == "True"
				} else if condition.Type == "Progressing" {
					progressing = condition.Status == "True"
				}
			}

			if !available && !progressing {
				if oldCondition.Type == release.Status_PENDING_INSTALL.String() {
					rls.Status.SetCondition(shipitv1beta1.HelmReleaseCondition{
						Type:    release.Status_FAILED.String(),
						Reason:  oldCondition.Reason,
						Message: "Deployments in the release did not roll out successfully",
					})
					r.Log.Info("Chart install was successful but Deployments did not complete.")
					return ctrl.Result{}, r.Status().Update(ctx, &rls)
				}

				r.Log.Info("Chart upgrade was successful but Deployments did not complete. Rolling back.")
				return r.rollback(ctx, rls)
			}

			if !available && progressing {
				rls.Status.SetCondition(shipitv1beta1.HelmReleaseCondition{
					Type:    oldCondition.Type,
					Reason:  oldCondition.Reason,
					Message: "Deployments in the release are still rolling out",
				})
				return ctrl.Result{RequeueAfter: r.GracePeriod}, r.Status().Update(ctx, &rls)
			}
		}

		rls.Status.SetCondition(shipitv1beta1.HelmReleaseCondition{
			Type:    releaseStatusCode.String(),
			Reason:  reason,
			Message: releaseStatus.GetNotes(),
		})

		return ctrl.Result{}, r.Status().Update(ctx, &rls)
	case release.Status_FAILED:
		rls.Status.SetCondition(shipitv1beta1.HelmReleaseCondition{
			Type:    releaseStatusCode.String(),
			Reason:  reasonForFailure(oldCondition),
			Message: releaseStatus.GetNotes(),
		})

		if err := r.Status().Update(ctx, &rls); err != nil {
			return ctrl.Result{}, err
		}

		if oldCondition.Type == release.Status_PENDING_UPGRADE.String() {
			return r.rollback(ctx, rls)
		}

		return ctrl.Result{}, nil
	default: // Status_UNKNOWN, Status_SUPERSEDED, Status_DELETED
		return ctrl.Result{}, nil
	}
}

func (r *HelmReleaseReconciler) install(ctx context.Context, rls shipitv1beta1.HelmRelease) (ctrl.Result, error) {
	releaseChart := rls.Spec.Chart
	releaseName := rls.Spec.ReleaseName

	chartVersion := fmt.Sprintf("%s@%s", releaseChart.URL(), releaseChart.Version)

	chart, err := r.downloader.Download(ctx, releaseChart.URL(), releaseChart.Version)
	if err != nil {
		return ctrl.Result{}, errors.Wrapf(err, "failed to download chart %s", chartVersion)
	}

	// TODO: use the returned response's `Release.Manifest` to watch and
	// receive events for the k8s resources owned by this chart
	if _, err := r.helm.InstallReleaseFromChart(chart, r.Namespace, helm.ReleaseName(releaseName), helm.ValueOverrides(rls.Spec.Values.Raw)); err != nil {
		return ctrl.Result{}, errors.Wrapf(err, "failed to install release %s using chart %s", releaseName, chartVersion)
	}

	rls.Status.SetCondition(shipitv1beta1.HelmReleaseCondition{
		Type:    release.Status_PENDING_INSTALL.String(),
		Message: fmt.Sprintf("installing chart %s", chartVersion),
	})

	if err := r.Status().Update(ctx, &rls); err != nil {
		return ctrl.Result{}, err
	}

	r.Log.Info("HelmRelease installed", "release", releaseName)
	return ctrl.Result{RequeueAfter: r.GracePeriod}, nil
}

func (r *HelmReleaseReconciler) rollback(ctx context.Context, rls shipitv1beta1.HelmRelease) (ctrl.Result, error) {
	releaseName := rls.Spec.ReleaseName

	if _, err := r.helm.RollbackRelease(releaseName); err != nil {
		return ctrl.Result{}, errors.Wrapf(err, "failed to rollback release %s", releaseName)
	}

	rls.Status.SetCondition(shipitv1beta1.HelmReleaseCondition{
		Type:    release.Status_PENDING_ROLLBACK.String(),
		Message: fmt.Sprintf("rolling back %s", releaseName),
	})

	r.Log.Info("HelmRelease rolled back", "release", releaseName)
	return ctrl.Result{}, r.Status().Update(ctx, &rls)
}

func (r *HelmReleaseReconciler) upgrade(ctx context.Context, rls shipitv1beta1.HelmRelease) (ctrl.Result, error) {
	releaseChart := rls.Spec.Chart
	releaseName := rls.Spec.ReleaseName

	chartVersion := fmt.Sprintf("%s@%s", releaseChart.URL(), releaseChart.Version)

	chart, err := r.downloader.Download(ctx, releaseChart.URL(), releaseChart.Version)
	if err != nil {
		return ctrl.Result{}, errors.Wrapf(err, "failed to download chart %s", chartVersion)
	}

	// TODO: use the returned response's `Release.Manifest` to watch and
	// receive events for the k8s resources owned by this chart
	if _, err := r.helm.UpdateReleaseFromChart(releaseName, chart, helm.UpdateValueOverrides(rls.Spec.Values.Raw)); err != nil {
		return ctrl.Result{}, errors.Wrapf(err, "failed to update release %s using chart %s", releaseName, chartVersion)
	}

	rls.Status.SetCondition(shipitv1beta1.HelmReleaseCondition{
		Type:    release.Status_PENDING_UPGRADE.String(),
		Message: fmt.Sprintf("upgrading chart %s", chartVersion),
	})

	if err := r.Status().Update(ctx, &rls); err != nil {
		return ctrl.Result{}, err
	}

	r.Log.Info("HelmRelease upgraded", "release", releaseName)
	return ctrl.Result{RequeueAfter: r.GracePeriod}, nil
}

func contains(strs []string, x string) bool {
	for _, s := range strs {
		if s == x {
			return true
		}
	}
	return false
}
