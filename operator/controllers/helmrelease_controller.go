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

	"ship-it-operator/notifications/slack"

	"github.com/go-logr/logr"
	"github.com/pkg/errors"
	apierrs "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/client-go/tools/record"
	"k8s.io/helm/pkg/helm"
	"k8s.io/helm/pkg/proto/hapi/chart"
	"k8s.io/helm/pkg/proto/hapi/release"
	hapi "k8s.io/helm/pkg/proto/hapi/services"
	helmerrors "k8s.io/helm/pkg/storage/errors"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
)

// HelmReleaseFinalizer allows the controller to clean up the associated release
// before the HelmRelease resource is deleted.
const HelmReleaseFinalizer = "HelmReleaseFinalizer"

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
	slack      *slack.Manager
	helm       HelmClient
	manager    ReleaseManager
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

func NewHelmReleaseReconciler(l logr.Logger, client client.Client, slackManager *slack.Manager, helm HelmClient, d ChartDownloader, rec record.EventRecorder, opts ...ReconcilerOption) *HelmReleaseReconciler {
	var cfg reconcilerConfig
	for _, opt := range opts {
		opt(&cfg)
	}

	return &HelmReleaseReconciler{
		Client: client,
		Log:    l.WithName("controllers").WithName("HelmRelease"),

		downloader:       d,
		slack:            slackManager,
		helm:             helm,
		reconcilerConfig: cfg,

		manager: ReleaseManager{
			helm:     helm,
			recorder: rec,
		},
	}
}

// +kubebuilder:rbac:groups=shipit.wattpad.com,resources=helmreleases,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=shipit.wattpad.com,resources=helmreleases/status,verbs=get;update;patch

func (r *HelmReleaseReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	ctx := context.Background()
	log := r.Log.WithValues("name", req.NamespacedName)

	helmRelease := new(shipitv1beta1.HelmRelease)

	if err := r.Get(ctx, req.NamespacedName, helmRelease); err != nil {
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

	if !hasFinalizer(helmRelease) {
		// setting the finalizer does not change the release's
		// metadata.generation, so we have to requeue
		return ctrl.Result{Requeue: true}, r.Update(ctx, setFinalizer(helmRelease))
	}

	return r.deploy(ctx, helmRelease)
}

func (r *HelmReleaseReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&shipitv1beta1.HelmRelease{}).
		WithEventFilter(predicate.GenerationChangedPredicate{}).
		Complete(r)
}

func contains(strs []string, x string) bool {
	for _, s := range strs {
		if s == x {
			return true
		}
	}
	return false
}

func hasFinalizer(rls *shipitv1beta1.HelmRelease) bool {
	return contains(rls.GetFinalizers(), HelmReleaseFinalizer)
}

func setFinalizer(rls *shipitv1beta1.HelmRelease) *shipitv1beta1.HelmRelease {
	finalizers := rls.GetFinalizers()
	rls.SetFinalizers(append(finalizers, HelmReleaseFinalizer))
	return rls
}

func clearFinalizer(rls *shipitv1beta1.HelmRelease) *shipitv1beta1.HelmRelease {
	var finalizers []string

	for _, f := range rls.GetFinalizers() {
		if f != HelmReleaseFinalizer {
			finalizers = append(finalizers, f)
		}
	}

	rls.SetFinalizers(finalizers)
	return rls
}

func (r *HelmReleaseReconciler) delete(ctx context.Context, rls *shipitv1beta1.HelmRelease) (ctrl.Result, error) {
	releaseName := rls.Spec.ReleaseName

	resp, err := r.helm.ReleaseStatus(releaseName)
	if err != nil {
		if isHelmReleaseNotFound(releaseName, err) {
			// this will only happen if a delete --purge is run
			return ctrl.Result{}, r.Update(ctx, clearFinalizer(rls))
		}

		return ctrl.Result{}, err
	}

	switch resp.GetInfo().GetStatus().GetCode() {
	case release.Status_DELETING:
		return ctrl.Result{RequeueAfter: r.GracePeriod}, nil
	case release.Status_DELETED:
		r.slack.Send(fmt.Sprintf("`%s` has been deleted.", releaseName))
		return ctrl.Result{}, r.Update(ctx, clearFinalizer(rls))
	}

	rls, err = r.manager.Delete(rls)
	if err != nil {
		return ctrl.Result{}, err
	}

	if err := r.Status().Update(ctx, rls); err != nil {
		return ctrl.Result{}, err
	}

	r.Log.Info("deleting HelmRelease", "release", releaseName)
	return ctrl.Result{RequeueAfter: r.GracePeriod}, nil
}

func isHelmReleaseNotFound(name string, err error) bool {
	// dynamic errors can't be directly compared for equality. We use the
	// error message string, though there's no guarantee it won't change.
	return strings.Contains(err.Error(), helmerrors.ErrReleaseNotFound(name).Error())
}

func (r *HelmReleaseReconciler) deploy(ctx context.Context, rls *shipitv1beta1.HelmRelease) (ctrl.Result, error) {
	releaseName := rls.Spec.ReleaseName

	resp, err := r.helm.ReleaseStatus(releaseName)
	if err != nil {
		if isHelmReleaseNotFound(releaseName, err) {
			return r.install(ctx, rls)
		}
		return ctrl.Result{}, errors.Wrapf(err, "failed to get release status for %s", releaseName)
	}

	oldCondition := rls.Status.GetCondition()

	switch statusCode := resp.GetInfo().GetStatus().GetCode(); statusCode {
	case release.Status_DELETING, release.Status_PENDING_INSTALL, release.Status_PENDING_UPGRADE, release.Status_PENDING_ROLLBACK:
		// if the release is still in transition, requeue until it settles
		return ctrl.Result{RequeueAfter: r.GracePeriod}, nil
	case release.Status_DELETED:
		return r.install(ctx, rls)
	case release.Status_DEPLOYED:
		if oldCondition.Type == release.Status_DEPLOYED.String() {
			return r.upgrade(ctx, rls)
		}

		r.slack.Send(fmt.Sprintf("`%s` is now deployed.", releaseName))
		return ctrl.Result{}, r.Status().Update(ctx, r.manager.Deployed(rls))
	case release.Status_FAILED:
		if err := r.Status().Update(ctx, r.manager.Failed(rls)); err != nil {
			return ctrl.Result{}, err
		}

		if oldCondition.Type == release.Status_PENDING_UPGRADE.String() {
			return r.rollback(ctx, rls)
		}

		return ctrl.Result{}, nil
	default: // Status_UNKNOWN
		return ctrl.Result{}, fmt.Errorf("unhandled release status code %s", statusCode)
	}
}

func (r *HelmReleaseReconciler) install(ctx context.Context, rls *shipitv1beta1.HelmRelease) (ctrl.Result, error) {
	chartSpec := rls.Spec.Chart
	releaseName := rls.Spec.ReleaseName

	chart, err := r.downloader.Download(ctx, chartSpec.URL(), chartSpec.Version)
	if err != nil {
		return ctrl.Result{}, errors.Wrapf(err, "failed to download chart %s", chartSpec.URL())
	}

	rls, err = r.manager.Install(rls, chart, r.Namespace)
	if err != nil {
		return ctrl.Result{}, errors.Wrapf(err, "failed to install release %s using chart %s", releaseName, chartSpec.URL())
	}

	if err := r.Status().Update(ctx, rls); err != nil {
		return ctrl.Result{}, err
	}

	r.Log.Info("installing HelmRelease", "release", releaseName)
	return ctrl.Result{RequeueAfter: r.GracePeriod}, nil
}

func (r *HelmReleaseReconciler) rollback(ctx context.Context, rls *shipitv1beta1.HelmRelease) (ctrl.Result, error) {
	releaseName := rls.Spec.ReleaseName

	rls, err := r.manager.Rollback(rls)
	if err != nil {
		return ctrl.Result{}, errors.Wrapf(err, "failed to roll back release %s", releaseName)
	}

	if err := r.Status().Update(ctx, rls); err != nil {
		return ctrl.Result{}, err
	}

	r.Log.Info("rolling back HelmRelease", "release", releaseName)
	r.slack.Send(fmt.Sprintf("`%s` is being rolled back.", releaseName))
	return ctrl.Result{RequeueAfter: r.GracePeriod}, nil
}

func (r *HelmReleaseReconciler) upgrade(ctx context.Context, rls *shipitv1beta1.HelmRelease) (ctrl.Result, error) {
	chartSpec := rls.Spec.Chart
	releaseName := rls.Spec.ReleaseName

	chart, err := r.downloader.Download(ctx, chartSpec.URL(), chartSpec.Version)
	if err != nil {
		return ctrl.Result{}, errors.Wrapf(err, "failed to download chart %s", chartSpec.URL())
	}

	rls, err = r.manager.Upgrade(rls, chart)
	if err != nil {
		return ctrl.Result{}, errors.Wrapf(err, "failed to upgrade release %s using chart %s", releaseName, chartSpec.URL())
	}

	if err := r.Status().Update(ctx, rls); err != nil {
		return ctrl.Result{}, err
	}

	r.Log.Info("upgrading HelmRelease", "release", releaseName)
	return ctrl.Result{RequeueAfter: r.GracePeriod}, nil
}
