package contracts

import (
	"github.com/gofiber/fiber/v2"
	"micro-fiber-test/pkg/commons"
	"micro-fiber-test/pkg/validation"
	"strconv"
	"strings"
)

func ConvertToInternalError(err error) commons.ApiError {
	return commons.ApiError{
		Code:    fiber.StatusInternalServerError,
		Kind:    string(commons.ErrorTypeTechnical),
		Message: err.Error(),
	}
}

func ConvertToFunctionalError(err error, targetStatus int) commons.ApiError {
	return commons.ApiError{
		Code:    targetStatus,
		Kind:    string(commons.ErrorTypeFunctional),
		Message: err.Error(),
	}
}

func ConvertValidationError(errors []validation.ErrorValidation) commons.ApiError {
	var s strings.Builder
	var details []commons.ApiErrorDetails
	for _, e := range errors {
		switch e.Error.Error() {
		case validation.ValidErrorNotBlank:
			s.WriteString("Field is null or empty")
		case validation.ValidErrorMaxLength:
			s.WriteString("Field length exceeds max value (")
			s.WriteString(strconv.Itoa(e.Size))
			s.WriteString(")")
		case validation.ValidErrorMinLength:
			s.WriteString("Field length exceeds min value (")
			s.WriteString(strconv.Itoa(e.Size))
			s.WriteString(")")
		default:
			s.WriteString("Unhandled error type")
		}
		details = append(details, commons.ApiErrorDetails{Field: e.Field, Detail: s.String()})
		s.Reset()
	}
	return commons.ApiError{
		Code:    fiber.StatusBadRequest,
		Kind:    string(commons.ErrorTypeFunctional),
		Message: s.String(),
		Details: details,
	}
}
