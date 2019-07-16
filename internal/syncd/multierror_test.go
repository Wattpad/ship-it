package syncd

import (
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestMultiErrorEmpty(t *testing.T) {
	var me multiError
	assert.EqualError(t, me, "")
}

func TestMultiErrorSingle(t *testing.T) {
	errFoo := errors.New("foo")

	var me multiError
	me.Add(errFoo)

	assert.EqualError(t, me, errFoo.Error())
}

func TestMultiErrorMultiple(t *testing.T) {
	errFoo := errors.New("foo")
	errBar := errors.New("bar")
	errBaz := errors.New(`failed to baz:
quux`)

	var me multiError
	me.Add(errFoo)
	me.Add(errBar)
	me.Add(errBaz)

	expected := `multiple errors (3):
	1: foo
	2: bar
	3: failed to baz:
quux`

	assert.EqualError(t, me, expected)
}
