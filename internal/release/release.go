package release

import (
	"fmt"

	"gopkg.in/yaml.v2"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"ship-it/pkg/apis/helmreleases.k8s.wattpad.com/v1alpha1"
)

type versioner struct {
	gvk schema.GroupVersionKind
}

func (v versioner) KindForGroupVersionKinds(kinds []schema.GroupVersionKind) (target schema.GroupVersionKind, ok bool) {
	return v.gvk, true
}

func NewDecoder() runtime.Decoder {
	factory := serializer.NewCodecFactory(runtime.NewScheme())
	d := factory.UniversalDeserializer()
	v := versioner{
		gvk: schema.FromAPIVersionAndKind("helmreleases.k8s.wattpad.com/v1alpha1", "HelmRelease"),
	}
	return factory.DecoderToVersion(d, v)
}

func ToYaml(h v1alpha1.HelmRelease) string {
	data := make(map[string]interface{})
	// create identical yaml to original with extra code gen fields
	data["apiVersion"] = h.APIVersion
	data["kind"] = h.Kind

	metadata := make(map[string]interface{})

	// only take name as it is the only populated field
	// this function may need to filter out unpopulated fields dynamically in the future
	metadata["name"] = h.ObjectMeta.Name

	data["metadata"] = metadata
	fmt.Println(metadata)

	spec := make(map[string]interface{})
	spec["releaseName"] = h.Spec.ReleaseName
	spec["chart"] = h.Spec.Chart
	spec["values"] = h.Spec.Values.Object

	data["spec"] = spec

	out, err := yaml.Marshal(data)
	if err != nil {
		fmt.Println(err)
		return ""
	}

	return string(out)
}
