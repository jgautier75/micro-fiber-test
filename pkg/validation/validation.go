package validation

import (
	"fmt"
	"github.com/go-playground/validator"
)

type ErrorType string

const (
	FieldRequired          = "Field %s is required"
	FieldMinLength         = "Field %s min length [%s] validation failed"
	FieldMaxLength         = "Field %s max length [%s] validation failed"
	GlobalValidationFailed = "validation_failed"
)

type ErrorValidation struct {
	Field string
	Error string
}

type ErrorResponse struct {
	FailedField string
	Tag         string
	Value       string
}

func ConvertValidationErrors(verr interface{}) []ErrorValidation {
	var errorsList []ErrorValidation
	for _, err := range verr.(validator.ValidationErrors) {
		var element ErrorValidation
		element.Field = err.StructNamespace()

		switch err.Tag() {
		case "required":
			element.Error = fmt.Sprintf(FieldRequired, err.Field())
			break
		case "min":
			element.Error = fmt.Sprintf(FieldMinLength, err.Field(), err.Param())
			break
		case "max":
			element.Error = fmt.Sprintf(FieldMaxLength, err.Field(), err.Param())
			break
		default:
			break
		}
		errorsList = append(errorsList, element)
	}
	return errorsList
}
