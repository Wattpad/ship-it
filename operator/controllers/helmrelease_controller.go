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

func WithNamespace(ns string) ReconcilerOption {
	return func(c *reconcilerConfig) {
		c.Namespace = ns
	}
}

func WithGracePeriod(d time.Duration) ReconcilerOption {
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

	if helmRelease.DeletionTimestamp != nil {
		return ctrl.Result{}, r.onDelete(ctx, helmRelease)
	}

	if !contains(helmRelease.GetFinalizers(), HelmReleaseFinalizer) {
		// setting the finalizer does not change the release's
		// metadata.generation, so we have to requeue
		return ctrl.Result{Requeue: true}, r.setFinalizer(ctx, helmRelease)
	}

	return r.onUpdate(ctx, helmRelease)
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

func (r *HelmReleaseReconciler) onDelete(ctx context.Context, rls shipitv1beta1.HelmRelease) error {
	// Update HelmRelease Status to 'DELETING'
	// Delete the release with helm
	return errNotImplemented
}

func isHelmReleaseNotFound(name string, err error) bool {
	// dynamic errors can't be directly compared for equality. We use the
	// error message string, though there's no guarantee it won't change.
	return strings.Contains(err.Error(), helmerrors.ErrReleaseNotFound(name).Error())
}

func (r *HelmReleaseReconciler) onUpdate(ctx context.Context, rls shipitv1beta1.HelmRelease) (ctrl.Result, error) {
	releaseName := rls.Spec.ReleaseName

	_, err := r.helm.ReleaseStatus(releaseName)
	if err != nil {
		if isHelmReleaseNotFound(releaseName, err) {
			return r.install(ctx, rls)
		}
		return ctrl.Result{}, errors.Wrapf(err, "failed to get release status for %s", releaseName)
	}

	return ctrl.Result{}, errNotImplemented
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
	return ctrl.Result{}, errNotImplemented
}

func contains(strs []string, x string) bool {
	for _, s := range strs {
		if s == x {
			return true
		}
	}
	return false
}
