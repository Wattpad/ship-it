package v1alpha1

import (
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer"
)

// Simple Implementation of GroupVersioner
// https://godoc.org/k8s.io/apimachinery/pkg/runtime#GroupVersioner
type versioner struct {
	gvk schema.GroupVersionKind
}

func (v versioner) KindForGroupVersionKinds(kinds []schema.GroupVersionKind) (target schema.GroupVersionKind, ok bool) {
	return v.gvk, true
}

func newDecoder() runtime.Decoder {
	factory := serializer.NewCodecFactory(runtime.NewScheme())
	d := factory.UniversalDeserializer()
	v := versioner{
		gvk: schema.FromAPIVersionAndKind(SchemeGroupVersion.String(), "HelmRelease"),
	}
	return factory.DecoderToVersion(d, v)
}

var Decoder = newDecoder()
