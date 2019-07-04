package v1alpha1

import (
	"encoding/json"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// HelmReleaseSpec defines the desired state of HelmRelease
type HelmReleaseSpec struct {
	ReleaseName string     `json:"releaseName"`
	Chart       ChartSpec  `json:"chart"`
	Values      HelmValues `json:"values"`
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

func (r HelmRelease) MarshalYAML() (interface{}, error) {
	rawJSON, err := json.Marshal(r)
	if err != nil {
		return nil, errors.Wrap(err, "unable to marshal HelmRelease into JSON")
	}

	var jsonObj interface{}
	if err := json.Unmarshal(rawJSON, &jsonObj); err != nil {
		return nil, errors.Wrap(err, "unable to unmarshal HelmRelease from JSON")
	}

	return jsonObj, err
}

// +k8s:deepcopy-gen:interfaces
// HelmValues allows us to implement runtime.Object for map[string]interface type.
type HelmValues map[string]interface{}

func (in *HelmValues) DeepCopyInto(out *HelmValues) {
	if in == nil {
		return
	}

	b, err := yaml.Marshal(in)
	if err != nil {
		return
	}

	var values HelmValues
	if err := yaml.Unmarshal(b, &values); err != nil {
		return
	}

	*out = values
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// HelmRelease is the Schema for the helmreleases API
type HelmRelease struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   HelmReleaseSpec   `json:"spec,omitempty"`
	Status HelmReleaseStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// HelmReleaseList contains a list of HelmRelease
type HelmReleaseList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []HelmRelease `json:"items"`
}
