package github

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"k8s.io/helm/pkg/helm"
	"k8s.io/helm/pkg/proto/hapi/chart"
	rls "k8s.io/helm/pkg/proto/hapi/services"
)

type helmClient interface {
	InstallReleaseFromChart(chart *chart.Chart, namespace string, opts ...helm.InstallOption) (*rls.InstallReleaseResponse, error)
	UpdateReleaseFromChart(rlsName string, chart *chart.Chart, opts ...helm.UpdateOption) (*rls.UpdateReleaseResponse, error)
}

type Reconciler struct {
	Namespace string
	Release   string
	Timeout   time.Duration

	client helmClient
}

func NewReconciler(c helmClient, namespace, release string, timeout time.Duration) *Reconciler {
	return &Reconciler{
		Namespace: namespace,
		Release:   release,
		Timeout:   timeout,
		client:    c,
	}
}

func (r *Reconciler) Reconcile(ctx context.Context, chart *chart.Chart) error {
	timeoutSeconds := int64(timeoutForDeadline(ctx, r.Timeout).Seconds())

	_, err := r.client.UpdateReleaseFromChart(r.Release, chart, helm.UpgradeTimeout(timeoutSeconds))
	err = errors.Wrap(err, "failed to update release from chart")

	if err != nil {
		_, err = r.client.InstallReleaseFromChart(chart, r.Namespace, helm.InstallTimeout(timeoutSeconds))
		err = errors.Wrap(err, "failed to install release from chart")
	}

	return err
}

func timeoutForDeadline(ctx context.Context, def time.Duration) time.Duration {
	if deadline, ok := ctx.Deadline(); ok {
		return time.Until(deadline)
	}

	return def
}
