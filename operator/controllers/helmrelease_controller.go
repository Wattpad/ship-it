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
)

// The HelmReleaseFinalizer allows the controller to clean up the associated
// release before the HelmRelease resource is deleted.
const HelmReleaseFinalizer = "HelmReleaseFinalizer"

var errNotImplemented = errors.New("not implemented")

type ChartDownloader interface {
	Download(ctx context.Context, chart string) (*chart.Chart, error)
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
		Complete(r)
}

func (r *HelmReleaseReconciler) setFinalizer(ctx context.Context, rls shipitv1beta1.HelmRelease) error {
	finalizers := rls.GetFinalizers()
	rls.SetFinalizers(append(finalizers, HelmReleaseFinalizer))

	return r.Update(ctx, &rls)
}

func (r *HelmReleaseReconciler) clearFinalizer(ctx context.Context, rls shipitv1beta1.HelmRelease) error {
	finalizers := []string{}

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

	if err := r.Update(ctx, &rls); err != nil {
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
	case release.Status_FAILED:
		rls.Status.SetCondition(shipitv1beta1.HelmReleaseCondition{
			Type:    releaseStatusCode.String(),
			Reason:  reasonForFailure(oldCondition),
			Message: releaseStatus.GetNotes(),
		})

		if err := r.Update(ctx, &rls); err != nil {
			r.Log.Error(err, "failed to update HelmRelease status", "release", releaseName, "status", releaseStatusCode.String())
		}

		if oldCondition.Type == release.Status_PENDING_UPGRADE.String() {
			return r.rollback(ctx, rls)
		}

		return ctrl.Result{}, nil
	}

	return ctrl.Result{}, errors.Wrap(errNotImplemented, "upgrade")
}

func (r *HelmReleaseReconciler) install(ctx context.Context, rls shipitv1beta1.HelmRelease) (ctrl.Result, error) {
	chartURI := rls.Spec.Chart.URI()
	releaseName := rls.Spec.ReleaseName

	chart, err := r.downloader.Download(ctx, chartURI)
	if err != nil {
		return ctrl.Result{}, errors.Wrapf(err, "failed to download chart %s", chartURI)
	}

	// TODO: use the returned response's `Release.Manifest` to watch and
	// receive events for the k8s resources owned by this chart
	if _, err := r.helm.InstallReleaseFromChart(chart, r.Namespace, helm.ReleaseName(releaseName)); err != nil {
		return ctrl.Result{}, errors.Wrapf(err, "failed to install release %s using chart %s", releaseName, chartURI)
	}

	rls.Status.SetCondition(shipitv1beta1.HelmReleaseCondition{
		Type:    release.Status_PENDING_INSTALL.String(),
		Message: fmt.Sprintf("installing chart %s", chartURI),
	})

	if err := r.Update(ctx, &rls); err != nil {
		r.Log.Info("failed to update HelmRelease status", "release", releaseName, "status", release.Status_PENDING_INSTALL)
	}

	return ctrl.Result{RequeueAfter: r.GracePeriod}, nil
}

func (r *HelmReleaseReconciler) rollback(ctx context.Context, rls shipitv1beta1.HelmRelease) (ctrl.Result, error) {
	_, err := r.helm.RollbackRelease(rls.Spec.ReleaseName)
	if err != nil {
		r.Log.Error(err, "unable to rollback release", "release", rls.Spec.ReleaseName)

		return ctrl.Result{
			RequeueAfter: r.GracePeriod,
		}, nil
	}

	rls.Status.SetCondition(shipitv1beta1.HelmReleaseCondition{
		Type:    release.Status_PENDING_ROLLBACK.String(),
		Message: fmt.Sprintf("rolling back %s", rls.Spec.ReleaseName),
	})

	return ctrl.Result{}, r.Update(ctx, &rls)
}

func contains(strs []string, x string) bool {
	for _, s := range strs {
		if s == x {
			return true
		}
	}
	return false
}
