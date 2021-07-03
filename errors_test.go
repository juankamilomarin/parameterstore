package parameterstore

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewErrTagNotSetOrEmpty(t *testing.T) {
	e := newErrTagNotSetOrEmpty("testField")
	assert.Equal(t, "tag not set or empty for field testField", e.Error())
}

func TestNewErrParsingParameter(t *testing.T) {
	e := newErrParsingParameter("testField", "testTag", "testDesc")
	assert.Equal(t, "cannot parse parameter testTag to field testField of type testDesc", e.Error())
}

func TestNewErrTypeNotSupported(t *testing.T) {
	e := newErrTypeNotSupported("testField", "typeDesc")
	assert.Equal(t, "type typeDesc for field testField not supported", e.Error())
}

func TestNewErrFieldCannotBeSet(t *testing.T) {
	e := newErrFieldCannotBeSet("testField")
	assert.Equal(t, "field testField cannot be set", e.Error())
}
