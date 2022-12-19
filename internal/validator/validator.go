package validator

import (
	"regexp"
	"unicode"

	"github.com/go-playground/validator/v10"
	"github.com/hatchet-dev/hatchet/internal/models/uuidutils"
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

	validate.RegisterValidation("password", func(fl validator.FieldLevel) bool {
		return passwordValidation(fl.Field().String())
	})

	validate.RegisterValidation("uuid", func(fl validator.FieldLevel) bool {
		return uuidutils.IsValidUUID(fl.Field().String())
	})

	return validate
}

func passwordValidation(pw string) bool {
	pwLen := len(pw)
	var hasNumber, hasUpper, hasLower bool

	for _, char := range pw {
		switch {
		case unicode.IsNumber(char):
			hasNumber = true
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		}
	}

	return hasNumber && hasUpper && hasLower && pwLen >= 8 && pwLen <= 32
}
