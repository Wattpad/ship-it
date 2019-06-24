package ecrconsumer

// Definition of Custom Resource YAML
type CRYaml struct {
	ApiVersion string `yaml:"apiVersion"`
	Kind       string `yaml:"kind"`
	Metadata   struct {
		Name string `yaml:"name"`
	}
	Spec struct {
		ReleaseName string `yaml:"releaseName"`
		Chart       struct {
			Repository string `yaml:"repository"`
			Path       string `yaml:"path"`
			Revision   string `yaml:"revision"`
		}
		Values struct {
			Image struct {
				Repository string `yaml:"repository"`
				Tag        string `yaml:"tag"`
			}
			IamRoleName        string `yaml:"iamRoleName"`
			ServiceAccountName string `yaml:"serviceAccountName"`
			AutoScaler         struct {
				MinPods                     string `yaml:"minPods"`
				MaxPods                     string `yaml:"maxPods"`
				TargetCPUUtilizationPercent string `yaml:"targetCPUUtilizationPercent"`
			}
		}
	}
}