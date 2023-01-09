package handlerutils_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hatchet-dev/hatchet/api/serverutils/apierrors"
	"github.com/hatchet-dev/hatchet/api/serverutils/handlerutils"
	"github.com/stretchr/testify/assert"
)

const (
	requiredErrorFmt        = "validation failed on field '%s' on condition 'required'"
	simpleConditionErrorFmt = "validation failed on field '%s' on condition '%s'"
	paramErrorFmt           = "validation failed on field '%s' on condition '%s' [ %s ]: got %s"
)

type validationErrObjectTest struct {
	valErrObj *handlerutils.ValidationErrObject
	expStr    string
}

var validationErrObjectTests = []validationErrObjectTest{
	{
		valErrObj: &handlerutils.ValidationErrObject{
			Field:       "username",
			Condition:   "required",
			Param:       "",
			ActualValue: nil,
		},
		expStr: fmt.Sprintf(requiredErrorFmt, "username"),
	},
	{
		valErrObj: &handlerutils.ValidationErrObject{
			Field:       "role",
			Condition:   "oneof",
			Param:       "admin developer",
			ActualValue: "notarole",
		},
		expStr: fmt.Sprintf(paramErrorFmt, "role", "oneof", "admin developer", "'notarole'"),
	},
	{
		valErrObj: &handlerutils.ValidationErrObject{
			Field:       "role",
			Condition:   "oneof",
			Param:       "admin developer",
			ActualValue: 1,
		},
		expStr: fmt.Sprintf(paramErrorFmt, "role", "oneof", "admin developer", "1"),
	},
	{
		valErrObj: &handlerutils.ValidationErrObject{
			Field:       "role",
			Condition:   "oneof",
			Param:       "admin developer",
			ActualValue: []string{"dev1", "dev2"},
		},
		expStr: fmt.Sprintf(paramErrorFmt, "role", "oneof", "admin developer", "[ dev1 dev2 ]"),
	},
	{
		valErrObj: &handlerutils.ValidationErrObject{
			Field:       "role",
			Condition:   "oneof",
			Param:       "admin developer",
			ActualValue: []int{1, 2},
		},
		expStr: fmt.Sprintf(paramErrorFmt, "role", "oneof", "admin developer", "[ 1 2 ]"),
	},
	{
		valErrObj: &handlerutils.ValidationErrObject{
			Field:     "role",
			Condition: "oneof",
			Param:     "admin developer",
			// for nil values, we convert the actual value to null
			ActualValue: nil,
		},
		expStr: fmt.Sprintf(paramErrorFmt, "role", "oneof", "admin developer", "null"),
	},
	{
		valErrObj: &handlerutils.ValidationErrObject{
			Field:     "role",
			Condition: "oneof",
			Param:     "admin developer",
			// for unrecognized types, we don't cast to value
			ActualValue: map[string]string{
				"not": "cast",
			},
		},
		expStr: fmt.Sprintf(paramErrorFmt, "role", "oneof", "admin developer", "invalid type"),
	},
}

func TestValidationErrObject(t *testing.T) {
	assert := assert.New(t)

	for _, test := range validationErrObjectTests {
		// test that the function outputs the expected readable error message
		readableStr := test.valErrObj.SafeExternalError()
		expReadableStr := test.expStr

		assert.Equal(
			expReadableStr,
			readableStr,
			"readable string not equal: expected %s, got %s",
			expReadableStr,
			readableStr,
		)
	}
}

type validationTest struct {
	description   string
	valObj        interface{}
	expErr        bool
	expErrStrings []string
}

type validationTestObj struct {
	ID    string `form:"required"`
	Name  string `form:"required"`
	Email string `form:"email"`
	Role  string `form:"oneof=admin developer"`
}

var validationTests = []validationTest{
	{
		description: "Missing all fields",
		valObj:      &validationTestObj{},
		expErr:      true,
		expErrStrings: []string{
			fmt.Sprintf(requiredErrorFmt, "ID"),
			fmt.Sprintf(requiredErrorFmt, "Name"),
			string(handlerutils.EmailErr),
			fmt.Sprintf(paramErrorFmt, "Role", "oneof", "admin developer", "''"),
		},
	},
	{
		description: "Fails email validation",
		valObj: &validationTestObj{
			ID:    "1",
			Name:  "whatever",
			Email: "notanemail",
			Role:  "admin",
		},
		expErr: true,
		expErrStrings: []string{
			string(handlerutils.EmailErr),
		},
	},
	{
		description: "Should pass all",
		valObj: &validationTestObj{
			ID:    "1",
			Name:  "whatever",
			Email: "anemail@gmail.com",
			Role:  "admin",
		},
		expErr:        false,
		expErrStrings: []string{},
	},
}

func TestValidation(t *testing.T) {
	assert := assert.New(t)
	validator := handlerutils.NewDefaultValidator()

	for _, test := range validationTests {
		// test that the function outputs the expected readable error message
		err := validator.Validate(test.valObj)

		assert.Equal(
			err != nil,
			test.expErr,
			"[ %s ]: expected error was %t, got %t",
			test.description,
			err != nil,
			test.expErr,
		)

		if err != nil && test.expErr {
			readableStrArr := strings.Split(err.Error(), ", ")
			expReadableStrArr := test.expErrStrings

			assert.ElementsMatch(
				expReadableStrArr,
				readableStrArr,
				"[ %s ]: readable string not equal",
				test.description,
			)

			// check that external and internal errors are returned as well
			assert.Equal(
				400,
				err.GetStatusCode(),
				"[ %s ]: status code not equal",
				test.description,
			)
		}
	}
}

func TestValidationNilParam(t *testing.T) {
	validator := handlerutils.NewDefaultValidator()

	err := validator.Validate(nil)
	expErr := apierrors.NewErrInternal(fmt.Errorf("could not cast err to validator.ValidationErrors"))

	// check that error type is of type apierrors.RequestError and that
	// message is correct
	assert.EqualError(t, err, expErr.Error(), "nil param error not internal server error")

	var expErrTarget apierrors.RequestError
	assert.ErrorAs(t, err, &expErrTarget)
}

func TestErrFailedRequestValidation(t *testing.T) {
	assert := assert.New(t)

	// just check that status code is 400 and all errors are set
	expErrStr := "readable error"
	err := handlerutils.NewErrFailedRequestValidation(expErrStr)

	assert.Equal(
		expErrStr,
		err.Error(),
		"incorrect value for Error() method",
	)

	assert.Equal(
		expErrStr,
		err.APIError().Errors[0].Description,
		"incorrect value for ExternalError() method",
	)

	// check that the status code is 400
	assert.Equal(
		expErrStr,
		err.InternalError(),
		"incorrect value for InternalError() method",
	)
}
