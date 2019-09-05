package controllers

import (
	"ship-it-operator/api/v1beta1"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"k8s.io/client-go/tools/record"
	"k8s.io/helm/pkg/helm"
	"k8s.io/helm/pkg/proto/hapi/chart"
	hapi "k8s.io/helm/pkg/proto/hapi/release"
)

var _ = Describe("ReleaseManager", func() {
	releaseName := "test-release"

	var (
		fakeHelm     *helm.FakeClient
		fakeRecorder *record.FakeRecorder
		manager      *ReleaseManager
		release      *v1beta1.HelmRelease
	)

	_ = BeforeEach(func() {
		fakeHelm = new(helm.FakeClient)
		fakeRecorder = record.NewFakeRecorder(1)

		release = &v1beta1.HelmRelease{
			Spec: v1beta1.HelmReleaseSpec{
				ReleaseName: releaseName,
			},
		}

		manager = &ReleaseManager{
			helm:     fakeHelm,
			recorder: fakeRecorder,
		}

	})

	It("should manage the release's lifecycle", func() {
		By("installing a new release")
		got, err := manager.Install(release, &chart.Chart{}, releaseName)
		Expect(err).To(BeNil())
		Expect(got.Status.GetCondition().Type).To(Equal(hapi.Status_PENDING_INSTALL.String()))
		Expect(<-fakeRecorder.Events).To(ContainSubstring(hapi.Status_PENDING_INSTALL.String()))

		resp, err := fakeHelm.ReleaseStatus(releaseName)
		Expect(err).To(BeNil())
		Expect(resp.GetInfo().GetStatus().GetCode()).To(Equal(hapi.Status_DEPLOYED))

		got = manager.Deployed(release)
		Expect(got.Status.GetCondition().Type).To(Equal(hapi.Status_DEPLOYED.String()))
		Expect(<-fakeRecorder.Events).To(ContainSubstring(hapi.Status_DEPLOYED.String()))

		By("upgrading an installed release")
		got, err = manager.Upgrade(release, &chart.Chart{})
		Expect(err).To(BeNil())
		Expect(got.Status.GetCondition().Type).To(Equal(hapi.Status_PENDING_UPGRADE.String()))
		Expect(<-fakeRecorder.Events).To(ContainSubstring(hapi.Status_PENDING_UPGRADE.String()))

		resp, err = fakeHelm.ReleaseStatus(releaseName)
		Expect(err).To(BeNil())
		Expect(resp.GetInfo().GetStatus().GetCode()).To(Equal(hapi.Status_DEPLOYED))

		got = manager.Deployed(release)
		Expect(got.Status.GetCondition().Type).To(Equal(hapi.Status_DEPLOYED.String()))
		Expect(<-fakeRecorder.Events).To(ContainSubstring(hapi.Status_DEPLOYED.String()))

		By("rolling back a failed upgraded release")
		got = manager.Failed(release)
		Expect(got.Status.GetCondition().Type).To(Equal(hapi.Status_FAILED.String()))
		Expect(<-fakeRecorder.Events).To(ContainSubstring(hapi.Status_FAILED.String()))

		got, err = manager.Rollback(release)
		Expect(err).To(BeNil())
		Expect(got.Status.GetCondition().Type).To(Equal(hapi.Status_PENDING_ROLLBACK.String()))
		Expect(<-fakeRecorder.Events).To(ContainSubstring(hapi.Status_PENDING_ROLLBACK.String()))

		resp, err = fakeHelm.ReleaseStatus(releaseName)
		Expect(err).To(BeNil())
		Expect(resp.GetInfo().GetStatus().GetCode()).To(Equal(hapi.Status_DEPLOYED))

		got = manager.Deployed(release)
		Expect(got.Status.GetCondition().Type).To(Equal(hapi.Status_DEPLOYED.String()))
		Expect(<-fakeRecorder.Events).To(ContainSubstring(hapi.Status_DEPLOYED.String()))

		By("deleting a release")
		got, err = manager.Delete(release)
		Expect(err).To(BeNil())
		Expect(got.Status.GetCondition().Type).To(Equal(hapi.Status_DELETING.String()))
		Expect(<-fakeRecorder.Events).To(ContainSubstring(hapi.Status_DELETING.String()))

		// deleted release not found
		_, err = fakeHelm.ReleaseStatus(releaseName)
		Expect(err).To(Not(BeNil()))
	})
})
