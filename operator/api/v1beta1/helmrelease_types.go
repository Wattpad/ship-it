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
	"path"
	"strconv"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

type HelmReleaseStatusReason string

const (
	ReasonUnknown         HelmReleaseStatusReason = "Unknown"
	ReasonDeleteError     HelmReleaseStatusReason = "DeleteError"
	ReasonInstallError    HelmReleaseStatusReason = "InstallError"
	ReasonInstallSuccess  HelmReleaseStatusReason = "InstallSuccess"
	ReasonRollbackError   HelmReleaseStatusReason = "RollbackError"
	ReasonRollbackSuccess HelmReleaseStatusReason = "RollbackSuccess"
	ReasonUpdateError     HelmReleaseStatusReason = "UpdateError"
	ReasonUpdateSuccess   HelmReleaseStatusReason = "UpdateSuccess"
)

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

	Conditions []HelmReleaseCondition `json:"conditions,omitempty"`
}

type HelmReleaseCondition struct {
	Type               string                  `json:"type"`
	LastTransitionTime metav1.Time             `json:"lastTransitionTime,omitempty"`
	LastUpdateTime     metav1.Time             `json:"lastUpdateTime,omitempty"`
	Message            string                  `json:"message,omitempty"`
	Reason             HelmReleaseStatusReason `json:"reason,omitempty"`
}

// ChartSpec defines the desired Helm chart
type ChartSpec struct {
	Repository string `json:"repository"`
	Path       string `json:"path"`
	Revision   string `json:"revision"`
}

func (cs *ChartSpec) URI() string {
	return path.Join(cs.Repository, cs.Path)
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

func (hr HelmRelease) HelmValues() map[string]interface{} {
	var obj map[string]interface{}
	if err := json.Unmarshal(hr.Spec.Values.Raw, &obj); err != nil {
		// this is temporary until we're using a representation of the
		// values that doesn't force us to handle this impossible error
		return make(map[string]interface{})
	}

	return obj
}

func (s *HelmReleaseStatus) SetCondition(condition HelmReleaseCondition) *HelmReleaseStatus {
	now := metav1.Now()
	condition.LastUpdateTime = now

	// if there's a matching condition, use the previous transition time
	for i, c := range s.Conditions {
		if c.Type == condition.Type && c.Reason == condition.Reason {
			condition.LastTransitionTime = c.LastTransitionTime
		} else {
			condition.LastTransitionTime = now
		}

		s.Conditions[i] = condition
		return s
	}

	// otherwise add the new condition
	condition.LastTransitionTime = now
	s.Conditions = append(s.Conditions, condition)
	return s
}

func (s *HelmReleaseStatus) GetCondition() HelmReleaseCondition {
	for _, c := range s.Conditions {
		return c
	}
	return HelmReleaseCondition{}
}

// +kubebuilder:object:root=true

// HelmReleaseList contains a list of HelmRelease
type HelmReleaseList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []HelmRelease `json:"items"`
}

type helmReleaseAnnotations map[string]string

// Annotations returns a map of annotations
func (hr HelmRelease) Annotations() helmReleaseAnnotations {
	return hr.ObjectMeta.GetAnnotations()
}

// Get returns an annotation
func (a helmReleaseAnnotations) Get(k string) string {
	return a[k]
}

// GetNamespaced returns the value of the annotation prefixed with the shipit API namespace
func (a helmReleaseAnnotations) GetNamespaced(k string) string {
	k = Resource("helmreleases").String() + "/" + k
	return a.Get(k)
}

func (a helmReleaseAnnotations) AutoDeploy() bool {
	autodeploy := a.GetNamespaced("autodeploy")

	v, err := strconv.ParseBool(autodeploy)

	if err != nil {
		return false
	}

	return v
}

func (a helmReleaseAnnotations) Code() string {
	return a.GetNamespaced("code")
}

func (a helmReleaseAnnotations) Datadog() string {
	return a.GetNamespaced("datadog")
}

func (a helmReleaseAnnotations) Squad() string {
	return a.GetNamespaced("squad")
}

func (a helmReleaseAnnotations) Slack() string {
	return a.GetNamespaced("slack")
}

func (a helmReleaseAnnotations) Sumologic() string {
	return a.GetNamespaced("sumologic")
}

func init() {
	SchemeBuilder.Register(&HelmRelease{}, &HelmReleaseList{})
}
