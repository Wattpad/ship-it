package ecrconsumer

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Definition of Custom Resource YAML
type CRYaml struct {
	metav1.TypeMeta   `yaml:",inline"`
	metav1.ObjectMeta `yaml:"metadata"`
	Spec              struct {
		ReleaseName string `yaml:"releaseName"`
		Chart       struct {
			Repository string `yaml:"repository"`
			Path       string `yaml:"path"`
			Revision   string `yaml:"revision"`
		}
		Values map[string]interface{}
	}
}
