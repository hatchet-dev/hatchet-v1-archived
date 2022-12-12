package validator

import (
	"regexp"

	"github.com/go-playground/validator/v10"
)

var NameRegex = regexp.MustCompile("^[a-zA-Z0-9\\.\\-_]+$")

// New creates a new instance of validator and sets the tag name
// to "form", instead of "validate"
func New() *validator.Validate {
	validate := validator.New()
	validate.SetTagName("form")

	validate.RegisterValidation("hatchet-name", func(fl validator.FieldLevel) bool {
		return NameRegex.MatchString(fl.Field().String())
	})

	return validate
}
