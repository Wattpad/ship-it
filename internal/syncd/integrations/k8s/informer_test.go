package k8s

import (
	"context"
	"encoding/json"
	"testing"

	shipitv1beta1 "ship-it-operator/api/v1beta1"

	"github.com/stretchr/testify/assert"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/cache/informertest"
)

func newFakeCache() *informertest.FakeInformers {
	scheme := runtime.NewScheme()
	shipitv1beta1.AddToScheme(scheme)

	return &informertest.FakeInformers{
		Scheme: scheme,
	}
}

func TestLookup(t *testing.T) {
	helmReleaseKind := "HelmRelease"
	wordCountsRelease := "word-counts-release"
	wordCountsRepository := "word-counts-repo"

	fakeCache := newFakeCache()

	fakeInformer, err := fakeCache.FakeInformerForKind(shipitv1beta1.Kind(helmReleaseKind))
	if !assert.NoError(t, err) {
		t.FailNow()
	}

	informer, err := NewInformerWithCache(context.Background(), fakeCache)
	if !assert.NoError(t, err) {
		t.FailNow()
	}

	values := map[string]interface{}{
		"image": map[string]interface{}{
			"repository": wordCountsRepository,
			"tag":        "foo",
		},
	}

	valuesRaw, err := json.Marshal(values)
	if !assert.NoError(t, err) {
		t.FailNow()
	}

	originalHR := &shipitv1beta1.HelmRelease{
		TypeMeta: metav1.TypeMeta{
			Kind: helmReleaseKind,
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      wordCountsRelease,
			Namespace: v1.NamespaceDefault,
		},
		Spec: shipitv1beta1.HelmReleaseSpec{
			Values: runtime.RawExtension{
				Raw: valuesRaw,
			},
		},
	}

	// Add the release
	fakeInformer.Add(originalHR)

	expected := types.NamespacedName{
		Name:      wordCountsRelease,
		Namespace: v1.NamespaceDefault,
	}

	names, err := informer.Lookup(wordCountsRepository)
	if assert.NoError(t, err) {
		assert.Len(t, names, 1)
		assert.Equal(t, names[0], expected)
	}

	// modify the release
	var updatedHR shipitv1beta1.HelmRelease
	originalHR.DeepCopyInto(&updatedHR)

	customNamespace := "custom-namespace"
	updatedHR.ObjectMeta.Namespace = customNamespace

	fakeInformer.Update(originalHR, &updatedHR)

	expected = types.NamespacedName{
		Name:      wordCountsRelease,
		Namespace: customNamespace,
	}

	names, err = informer.Lookup(wordCountsRepository)
	if assert.NoError(t, err) {
		assert.Len(t, names, 1)
		assert.Equal(t, names[0], expected)
	}

	// delete the HelmRelease
	fakeInformer.Delete(&updatedHR)
	names, err = informer.Lookup(wordCountsRepository)
	if assert.NoError(t, err) {
		assert.Empty(t, names)
	}
}
