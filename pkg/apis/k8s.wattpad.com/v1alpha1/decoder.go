package v1alpha1

import (
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
)

var Decoder = newDecoder()

func newDecoder() runtime.Decoder {
	factory := serializer.NewCodecFactory(runtime.NewScheme())
	d := factory.UniversalDeserializer()

	return factory.DecoderToVersion(d, SchemeGroupVersion)
}
