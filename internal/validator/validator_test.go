package validator

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Validator(t *testing.T) {
	v := New()
	assert.True(t, v.Valid(), "New validator should be valid")
	v.AddError("test", "test error")
	assert.False(t, v.Valid(), "Validator with errors should be invalid")

	v = New()
	v.Check(true, "test", "test error")
	assert.True(t, v.Valid(), "Validator with no errors should be valid")
	v.Check(false, "test", "test error")
	assert.False(t, v.Valid(), "Validator with errors should be invalid")

	assert.True(t, PermittedValue("a", "a", "b", "c"), "Permitted value should return true")
	assert.False(t, PermittedValue("d", "a", "b", "c"), "Permitted value should return false")
}
