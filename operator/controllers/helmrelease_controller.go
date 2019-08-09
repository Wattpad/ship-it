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
	"errors"

	shipitv1beta1 "ship-it-operator/api/v1beta1"

	"github.com/go-logr/logr"
	apierrs "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/helm/pkg/helm"
	"k8s.io/helm/pkg/proto/hapi/chart"
	hapi "k8s.io/helm/pkg/proto/hapi/services"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// The HelmReleaseFinalizer allows the controller to clean up the associated
// release before the HelmRelease resource is deleted.
const HelmReleaseFinalizer = "HelmReleaseFinalizer"

var errNotImplemented = errors.New("not implemented")

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
	Log logr.Logger

	helm HelmClient
}

func NewHelmReleaseReconciler(l logr.Logger, client client.Client, helm HelmClient) *HelmReleaseReconciler {
	return &HelmReleaseReconciler{
		Client: client,
		Log:    l.WithName("controllers").WithName("HelmRelease"),
		helm:   helm,
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

func (r *HelmReleaseReconciler) onUpdate(ctx context.Context, rls shipitv1beta1.HelmRelease) (ctrl.Result, error) {
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
