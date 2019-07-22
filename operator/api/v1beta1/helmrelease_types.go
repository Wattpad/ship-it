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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// HelmReleaseSpec defines the desired state of HelmRelease
type HelmReleaseSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	ReleaseName string               `json:"releaseName"`
	Chart       ChartSpec            `json:"chart"`
	Values      runtime.RawExtension `json:"values"`
}

// HelmReleaseStatus defines the observed state of HelmRelease
type HelmReleaseStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

// ChartSpec defines the desired Helm chart
type ChartSpec struct {
	Repository string `json:"repository"`
	Path       string `json:"path"`
	Revision   string `json:"revision"`
}

// +kubebuilder:object:root=true
// +kubebuilder:resource:path=helmreleases,shortName=rls

// HelmRelease is the Schema for the helmreleases API
type HelmRelease struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   HelmReleaseSpec   `json:"spec,omitempty"`
	Status HelmReleaseStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// HelmReleaseList contains a list of HelmRelease
type HelmReleaseList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []HelmRelease `json:"items"`
}

func init() {
	SchemeBuilder.Register(&HelmRelease{}, &HelmReleaseList{})
}
