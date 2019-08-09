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

package v1beta1

import (
	"encoding/json"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"golang.org/x/net/context"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
)

// These tests are written in BDD-style using Ginkgo framework. Refer to
// http://onsi.github.io/ginkgo to learn more.

var _ = Describe("HelmRelease", func() {
	var (
		key              types.NamespacedName
		created, fetched *HelmRelease
	)

	BeforeEach(func() {
		// Add any setup steps that needs to be executed before each test
	})

	AfterEach(func() {
		// Add any teardown steps that needs to be executed after each test
	})

	// Add Tests for OpenAPI validation (or additonal CRD features) specified in
	// your API definition.
	// Avoid adding tests for vanilla CRUD operations because they would
	// test Kubernetes API server, which isn't the goal here.
	Context("Create API", func() {

		It("should create an object successfully", func() {

			key = types.NamespacedName{
				Name:      "foo",
				Namespace: "default",
			}
			created = &HelmRelease{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "foo",
					Namespace: "default",
				},
				Spec: HelmReleaseSpec{
					ReleaseName: "test",
					Chart: ChartSpec{
						Repository: "blah",
						Path:       "blah",
						Revision:   "blah",
					},
					Values: runtime.RawExtension{
						Raw: []byte(`{"test":1}`),
					},
				},
			}

			By("creating an API obj")
			Expect(k8sClient.Create(context.TODO(), created)).To(Succeed())

			fetched = &HelmRelease{}
			Expect(k8sClient.Get(context.TODO(), key, fetched)).To(Succeed())
			Expect(fetched).To(Equal(created))

			By("deleting the created object")
			Expect(k8sClient.Delete(context.TODO(), created)).To(Succeed())
			Expect(k8sClient.Get(context.TODO(), key, created)).ToNot(Succeed())
		})

	})

	Context("Serializing Helm values", func() {
		It("should unmarshal values as a JSON object", func() {
			expectedVals := map[string]interface{}{"hello": "world"}

			expectedRaw, _ := json.Marshal(expectedVals)

			hr := HelmRelease{
				Spec: HelmReleaseSpec{
					Values: runtime.RawExtension{
						Raw: expectedRaw,
					},
				},
			}

			By("calling HelmValues method")
			Expect(hr.HelmValues()).To(Equal(expectedVals))
		})
	})

	Context("Getting annotations", func() {
		It("should return the annotation", func() {
			hr := &HelmRelease{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "foo",
					Namespace: "default",
					Annotations: map[string]string{
						"test": "annotation",
						"helmreleases.shipit.wattpad.com/autodeploy": "true",
						"helmreleases.shipit.wattpad.com/code":       "code",
					},
				},
				Spec: HelmReleaseSpec{},
			}

			annotations := hr.GetAnnotations()

			By("calling AutoDeploy")
			Expect(annotations.AutoDeploy()).To(BeTrue())

			By("calling Get")
			Expect(annotations.Get("test")).To(Equal("annotation"))

			By("calling GetNamespaced")
			Expect(annotations.GetNamespaced("code")).To(Equal("code"))
		})
	})
})
