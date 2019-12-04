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
	"k8s.io/client-go/tools/record"
	"k8s.io/helm/pkg/helm"
	"k8s.io/helm/pkg/proto/hapi/chart"
	"k8s.io/helm/pkg/proto/hapi/release"
	hapi "k8s.io/helm/pkg/proto/hapi/release"
	ctrl "sigs.k8s.io/controller-runtime"
)

var sentNotifications []string

type fakeNotifier struct{}

func (m *fakeNotifier) Send(msg string) error {
	sentNotifications = append(sentNotifications, msg)
	return nil
}

type mockDownloader struct {
	mock.Mock
}

func (m *mockDownloader) Download(ctx context.Context, chartName string, version string) (*chart.Chart, error) {
	args := m.Called(ctx, chartName, version)

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

	testChart := &chart.Chart{
		Metadata: &chart.Metadata{
			Name: releaseName,
		},
	}

	releaseKey := types.NamespacedName{
		Name:      releaseName,
		Namespace: releaseNamespace,
	}

	request := ctrl.Request{NamespacedName: releaseKey}

	var recorder record.FakeRecorder

	go func() {
		for range recorder.Events {
			// discard events
		}
	}()

	var (
		downloader  *mockDownloader
		helmClient  *helm.FakeClient
		reconciler  *HelmReleaseReconciler
		testRelease *shipitv1beta1.HelmRelease
	)

	BeforeEach(func() {
		downloader = new(mockDownloader)
		helmClient = new(helm.FakeClient)
		reconciler = NewHelmReleaseReconciler(log, k8sClient, &fakeNotifier{}, helmClient, downloader, &recorder, GracePeriod(42), Namespace("test"))

		testRelease = &shipitv1beta1.HelmRelease{
			TypeMeta: metav1.TypeMeta{
				APIVersion: "shipit.wattpad.com/v1beta1",
				Kind:       "HelmRelease",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      releaseName,
				Namespace: releaseNamespace,
				Annotations: map[string]string{
					"helmreleases.shipit.wattpad.com/autodeploy": "true",
				},
			},
			Spec: shipitv1beta1.HelmReleaseSpec{
				ReleaseName: releaseName,
				Chart: shipitv1beta1.ChartSpec{
					Repository: "github.com/example/foo",
					Name:       "bar",
					Version:    "0.0.0",
				},
				Values: runtime.RawExtension{Raw: []byte("{}")},
			},
		}
	})

	AfterEach(func() {
		k8sClient.Delete(ctx, testRelease)
	})

	Context("HelmReleaseFinalizer", func() {
		It("should set and clear the finalizer", func() {
			Expect(hasFinalizer(testRelease)).To(BeFalse())

			By("setting the finalizer")
			got := setFinalizer(testRelease)
			Expect(hasFinalizer(testRelease)).To(BeTrue())
			Expect(got.GetFinalizers()).To(ContainElement(HelmReleaseFinalizer))

			By("clearing the finalizer")
			got = clearFinalizer(testRelease)
			Expect(hasFinalizer(testRelease)).To(BeFalse())
			Expect(got.GetFinalizers()).To(Not(ContainElement(HelmReleaseFinalizer)))
		})
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
			Expect(k8sClient.Create(ctx, testRelease)).To(Succeed())

			res, err := reconciler.Reconcile(request)
			Expect(err).To(BeNil())
			Expect(res).To(BeZero())
		})
	})

	When("the HelmRelease has autodeploy enabled", func() {
		It("should reconcile the release", func() {
			Expect(k8sClient.Create(ctx, testRelease)).To(Succeed())

			By("reconciling a new release without the HelmReleaseFinalizer")

			res, err := reconciler.Reconcile(request)
			Expect(err).To(BeNil())
			Expect(res.Requeue).To(BeTrue())

			var got shipitv1beta1.HelmRelease
			Expect(k8sClient.Get(ctx, releaseKey, &got)).To(Succeed())
			Expect(got.GetFinalizers()).To(ContainElement(HelmReleaseFinalizer))

			By("reconciling a new release")
			downloader.On("Download", ctx, testRelease.Spec.Chart.URL(), testRelease.Spec.Chart.Version).Return(testChart, nil)

			_, err = helmClient.ReleaseStatus(releaseName)
			Expect(isHelmReleaseNotFound(releaseName, err)).To(BeTrue())

			// first reconcile updates the HelmRelease to be
			// PENDING_INSTALL while it asks tiller to install the release
			res, err = reconciler.Reconcile(request)
			Expect(err).To(BeNil())
			Expect(res.RequeueAfter).To(Equal(reconciler.GracePeriod))

			Expect(k8sClient.Get(ctx, releaseKey, &got)).To(Succeed())
			Expect(got.Status.GetCondition().Type).To(Equal(hapi.Status_PENDING_INSTALL.String()))

			resp, err := helmClient.ReleaseStatus(releaseName)
			Expect(err).To(BeNil())
			Expect(resp.GetInfo().GetStatus().GetCode()).To(Equal(hapi.Status_DEPLOYED))

			// second reconcile observes the completed release
			// install and updates the HelmRelease to be DEPLOYED
			res, err = reconciler.Reconcile(request)
			Expect(err).To(BeNil())
			Expect(res).To(BeZero())

			Expect(k8sClient.Get(ctx, releaseKey, &got)).To(Succeed())
			Expect(got.Status.GetCondition().Type).To(Equal(hapi.Status_DEPLOYED.String()))

			By("reconciling an installed release")

			resp, err = helmClient.ReleaseStatus(releaseName)
			Expect(err).To(BeNil())
			Expect(resp.GetInfo().GetStatus().GetCode()).To(Equal(hapi.Status_DEPLOYED))

			downloader.On("Download", ctx, testRelease.Spec.Chart.URL(), testRelease.Spec.Chart.Version).Return(testChart, nil)

			// first reconcile updates the HelmRelease to be
			// PENDING_UPGRADE while it asks tiller to upgrade the release
			res, err = reconciler.Reconcile(request)
			Expect(err).To(BeNil())
			Expect(res.RequeueAfter).To(Equal(reconciler.GracePeriod))

			Expect(k8sClient.Get(ctx, releaseKey, &got)).To(Succeed())
			Expect(got.Status.GetCondition().Type).To(Equal(hapi.Status_PENDING_UPGRADE.String()))

			resp, err = helmClient.ReleaseStatus(releaseName)
			Expect(err).To(BeNil())
			Expect(resp.GetInfo().GetStatus().GetCode()).To(Equal(hapi.Status_DEPLOYED))

			// second reconcile observes the completed release
			// upgrade and updates the HelmRelease to be DEPLOYED
			res, err = reconciler.Reconcile(request)
			Expect(err).To(BeNil())
			Expect(res).To(BeZero())

			Expect(k8sClient.Get(ctx, releaseKey, &got)).To(Succeed())
			Expect(got.Status.GetCondition().Type).To(Equal(hapi.Status_DEPLOYED.String()))

			By("reconciling a failed updated release")

			// fake a failed release upgrade
			got.Status.SetCondition(shipitv1beta1.HelmReleaseCondition{
				Type: hapi.Status_PENDING_UPGRADE.String(),
			})

			Expect(k8sClient.Status().Update(ctx, &got)).To(Succeed())
			failedUpgrade(helmClient, releaseName)

			// first reconcile updates the HelmRelease to be
			// PENDING_ROLLBACK while it asks tiller to rollback the release
			res, err = reconciler.Reconcile(request)
			Expect(err).To(BeNil())
			Expect(res.RequeueAfter).To(Equal(reconciler.GracePeriod))

			// the fake helm client doesn't actually implement a
			// rollback, so we do it live
			fakedRollback(helmClient, releaseName)

			Expect(k8sClient.Get(ctx, releaseKey, &got)).To(Succeed())
			Expect(got.Status.GetCondition().Type).To(Equal(hapi.Status_PENDING_ROLLBACK.String()))

			// second reconcile observes the completed rollback and
			// updates the HelmRelease to be DEPLOYED
			res, err = reconciler.Reconcile(request)
			Expect(err).To(BeNil())
			Expect(res).To(BeZero())

			Expect(k8sClient.Get(ctx, releaseKey, &got)).To(Succeed())
			Expect(got.Status.GetCondition().Type).To(Equal(hapi.Status_DEPLOYED.String()))

			By("reconciling a release that has been deleted")
			Expect(k8sClient.Delete(ctx, testRelease)).To(Succeed())

			// first reconcile updates the HelmRelease to be
			// DELETING while it asks tiller to delete the release
			_, err = reconciler.Reconcile(request)
			Expect(err).To(BeNil())

			Expect(k8sClient.Get(ctx, releaseKey, &got)).To(Succeed())
			Expect(got.Status.GetCondition().Type).To(Equal(hapi.Status_DELETING.String()))

			// second reconcile observes the completed release
			// deletion and clears the HelmRelease's finalizer so it
			// can be deleted by kubernetes
			_, err = reconciler.Reconcile(request)
			Expect(err).To(BeNil())

			Expect(k8sClient.Get(ctx, releaseKey, &got)).To(Not(Succeed()))

			resp, err = helmClient.ReleaseStatus(releaseName)
			Expect(isHelmReleaseNotFound(releaseName, err)).To(BeTrue())
		})
	})
})

// reaching into the fake client internals to fake a failed release
func failedUpgrade(client *helm.FakeClient, name string) {
	for _, rls := range client.Rels {
		if rls.Name == name {
			rls.Info.Status.Code = release.Status_FAILED
			return
		}
	}
}

func fakedRollback(client *helm.FakeClient, name string) {
	for _, rls := range client.Rels {
		if rls.Name == name {
			rls.Info.Status.Code = release.Status_DEPLOYED
			return
		}
	}
}
