package helmrelease

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	unstructured "k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
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

func NewDecoder(obj runtime.Object, data []byte) (*runtime.Decoder, error) {
	factory := serializer.NewCodecFactory(runtime.NewScheme())
	// make universal deserializer
	decoder := factory.UniversalDeserializer()
	//fmt.Println(decoder)
	gvk := schema.FromAPIVersionAndKind("helmreleases.k8s.wattpad.com/v1alpha1", "HelmRelease")
	//fmt.Println()
	decoder.Decode(data, &gvk, obj)
	return nil, nil
}
