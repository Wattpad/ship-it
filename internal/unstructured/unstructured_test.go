package unstructured

import (
	"testing"

	"github.com/stretchr/testify/mock"
)

type mockCallback struct {
	mock.Mock
}

func (m *mockCallback) VoidMethod(x interface{}) {
	m.Called(x)
}

func TestFindAll(t *testing.T) {
	obj := map[string]interface{}{
		"bar": 0,
		"foo": map[string]interface{}{
			"bar": 1,
		},
		"qux": 2,
	}

	var cb mockCallback
	cb.On("VoidMethod", interface{}(0)).Once()
	cb.On("VoidMethod", interface{}(1)).Once()

	FindAll(obj, "bar", cb.VoidMethod)

	cb.AssertExpectations(t)
}
