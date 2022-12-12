package validator_test

import (
	"testing"

	"github.com/hatchet-dev/hatchet/internal/validator"
	"github.com/stretchr/testify/assert"
)

type nameResource struct {
	DisplayName string `form:"hatchet-name"`
}

func TestValidatorInvalidName(t *testing.T) {
	v := validator.New()

	err := v.Struct(&nameResource{
		DisplayName: "&&!!",
	})

	assert.ErrorContains(t, err, "validation for 'DisplayName' failed on the 'hatchet-name' tag", "should throw error on invalid name")
}

func TestValidatorValidName(t *testing.T) {
	v := validator.New()

	err := v.Struct(&nameResource{
		DisplayName: "test-name",
	})

	assert.NoError(t, err, "no error")
}
