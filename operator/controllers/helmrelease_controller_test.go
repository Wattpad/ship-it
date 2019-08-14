package controllers

import (
	"context"
	"fmt"
	shipitv1beta1 "ship-it-operator/api/v1beta1"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/mock"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/helm/pkg/helm"
	"k8s.io/helm/pkg/proto/hapi/chart"
	hapi "k8s.io/helm/pkg/proto/hapi/release"
	ctrl "sigs.k8s.io/controller-runtime"
)

type mockDownloader struct {
	mock.Mock
}

func (m *mockDownloader) Download(ctx context.Context, chartName string) (*chart.Chart, error) {
	args := m.Called(ctx, chartName)

	var ret0 *chart.Chart
	if args0 := args.Get(0); args0 != nil {
		ret0 = args0.(*chart.Chart)
	}

	return ret0, args.Error(1)
}

var _ = Describe("HelmReleaseReconciler", func() {
	ctx := context.Background()
	log := ctrl.Log.WithName("helmrelease_controller_test")

	releaseName := "test-release"
	releaseNamespace := "test-namespace"

	chart := &chart.Chart{
		Metadata: &chart.Metadata{
			Name: releaseName,
		},
	}

	releaseKey := types.NamespacedName{
		Name:      releaseName,
		Namespace: releaseNamespace,
	}

	request := ctrl.Request{NamespacedName: releaseKey}

	var (
		helmClient  *helm.FakeClient
		downloader  *mockDownloader
		reconciler  *HelmReleaseReconciler
		testRelease *shipitv1beta1.HelmRelease
	)

	BeforeEach(func() {
		downloader = new(mockDownloader)
		helmClient = new(helm.FakeClient)
		reconciler = NewHelmReleaseReconciler(log, k8sClient, helmClient, downloader)

		testRelease = &shipitv1beta1.HelmRelease{
			TypeMeta: metav1.TypeMeta{
				Kind: "HelmRelease",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:              releaseName,
				Namespace:         releaseNamespace,
				DeletionTimestamp: nil,
				Annotations: map[string]string{
					"helmreleases.shipit.wattpad.com/autodeploy": "true",
				},
			},
			Spec: shipitv1beta1.HelmReleaseSpec{
				ReleaseName: releaseName,
				Chart: shipitv1beta1.ChartSpec{
					Repository: "github.com/example/foo",
					Path:       "bar/baz",
				},
				Values: runtime.RawExtension{Raw: []byte("{}")},
			},
			Status: shipitv1beta1.HelmReleaseStatus{
				Conditions: []shipitv1beta1.HelmReleaseCondition{},
			},
		}
	})

	JustBeforeEach(func() {
		k8sClient.Delete(ctx, testRelease)
	})

	When("the HelmRelease isn't found", func() {
		It("should do nothing", func() {
			res, err := reconciler.Reconcile(request)
			Expect(err).To(BeNil())
			Expect(res).To(BeZero())
		})
	})

	When("the HelmRelease has autodeploy disabled", func() {
		It("should ignore the release", func() {
			ann := fmt.Sprintf("%s/autodeploy", shipitv1beta1.Resource("helmreleases"))
			testRelease.ObjectMeta.Annotations[ann] = "false"
			k8sClient.Create(ctx, testRelease)

			res, err := reconciler.Reconcile(request)
			Expect(err).To(BeNil())
			Expect(res).To(BeZero())
		})
	})

	When("the HelmRelease has autodeploy enabled", func() {
		It("should manage the HelmReleaseFinalizer", func() {
			err := k8sClient.Create(ctx, testRelease)

			By("reconciling autodeployed release")
			res, err := reconciler.Reconcile(request)
			Expect(err).To(BeNil())
			Expect(res.Requeue).To(BeTrue())

			var got shipitv1beta1.HelmRelease
			err = k8sClient.Get(ctx, releaseKey, &got)
			Expect(err).To(BeNil())
			Expect(got.GetFinalizers()).To(ContainElement(HelmReleaseFinalizer))
		})

		It("should install the helm release if it doesn't already exist", func() {
			downloader.On("Download", ctx, testRelease.Spec.Chart.URI()).Return(chart, nil)

			res, err := reconciler.Reconcile(request)
			Expect(err).To(BeNil())
			Expect(res.RequeueAfter).To(Equal(reconciler.GracePeriod))

			resp, err := helmClient.ReleaseStatus(releaseName)
			Expect(err).To(BeNil())
			Expect(resp.GetInfo().GetStatus().GetCode()).To(Equal(hapi.Status_DEPLOYED))
		})

		It("should update the helm release if it already exists", func() {
			Skip("TODO")
		})

		It("should rollback the helm release if it fails to update", func() {
			Skip("TODO")
		})

		It("should delete the helm release when the deletion timestamp is set", func() {
			Skip("TODO")
		})
	})
})
