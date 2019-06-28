package helmrelease

import (
	"fmt"

	"gopkg.in/yaml.v2"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	unstructured "k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
)

// HelmReleaseSpec defines the desired state of HelmRelease
type HelmReleaseSpec struct {
	ReleaseName string                    `json:"releaseName"`
	Chart       ChartSpec                 `json:"chart"`
	Values      unstructured.Unstructured `json:"values"`
}

// HelmReleaseSpec defines the desired Helm chart
type ChartSpec struct {
	Repository string `json:"repository"`
	Path       string `json:"path"`
	Revision   string `json:"revision"`
}

// HelmReleaseStatus defines the observed state of HelmRelease
type HelmReleaseStatus struct {
	// TODO
}

// HelmRelease is the Schema for the helmreleases API
type HelmRelease struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   HelmReleaseSpec   `json:"spec,omitempty"`
	Status HelmReleaseStatus `json:"status,omitempty"`
}

// HelmReleaseList contains a list of HelmRelease
type HelmReleaseList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []HelmRelease `json:"items"`
}

func (h *HelmRelease) DeepCopyObject() runtime.Object {
	return h
}

func NewDecoder() runtime.Decoder {
	factory := serializer.NewCodecFactory(runtime.NewScheme())

	return factory.UniversalDeserializer()
}

func (h HelmRelease) Encode() []byte {
	data := make(map[string]interface{})
	// create identical json to original with extra code gen fields
	data["apiVersion"] = h.APIVersion
	data["kind"] = h.Kind
	data["metadata"] = h.ObjectMeta

	spec := make(map[string]interface{})
	spec["releaseName"] = h.Spec.ReleaseName
	spec["chart"] = h.Spec.Chart
	spec["values"] = h.Spec.Values.Object

	data["spec"] = spec

	out, err := yaml.Marshal(data)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	fmt.Println(data)
	return out
}