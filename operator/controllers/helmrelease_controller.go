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

	"github.com/go-logr/logr"
	apierrs "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/helm/pkg/helm"
	"k8s.io/helm/pkg/proto/hapi/chart"
	rls "k8s.io/helm/pkg/proto/hapi/services"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	shipitv1beta1 "ship-it-operator/api/v1beta1"
)

var errNotImplemented = errors.New("not implemented")

type HelmClient interface {
	DeleteRelease(rlsName string, opts ...helm.DeleteOption) (*rls.UninstallReleaseResponse, error)
	InstallReleaseFromChart(chart *chart.Chart, ns string, opts ...helm.InstallOption) (*rls.InstallReleaseResponse, error)
	RollbackRelease(rlsName string, opts ...helm.RollbackOption) (*rls.RollbackReleaseResponse, error)
	UpdateReleaseFromChart(rlsName string, chart *chart.Chart, opts ...helm.UpdateOption) (*rls.UpdateReleaseResponse, error)
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
			return ctrl.Result{}, r.onDelete(ctx, req.NamespacedName)
		}

		return ctrl.Result{}, err
	}

	return ctrl.Result{}, r.onUpdate(ctx, helmRelease)
}

func (r *HelmReleaseReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&shipitv1beta1.HelmRelease{}).
		Complete(r)
}

func (r *HelmReleaseReconciler) onDelete(ctx context.Context, name types.NamespacedName) error {
	// Update HelmRelease Status to 'PENDING_DELETE' -- but the HelmRelease doesn't exist anymore :(
	// Delete the release with helm
	return errNotImplemented
}

func (r *HelmReleaseReconciler) onUpdate(ctx context.Context, rls shipitv1beta1.HelmRelease) error {
	// Update HelmRelease Status to 'PENDING_INSTALL' or 'PENDING_UPGRADE'
	// Attempt the install/upgrade with helm
	// If it succeeded, set Status to 'DEPLOYED'
	// Else, set Status to 'PENDING_ROLLBACK'
	// Attempt the rollback with helm
	// If it succeeded, set Status to 'DEPLOYED'
	// Else, set Status to 'FAILED'
	return errNotImplemented
}
