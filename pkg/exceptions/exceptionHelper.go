package exceptions

import (
	"github.com/gofiber/fiber/v2"
	"micro-fiber-test/pkg/dto/commons"
	"micro-fiber-test/pkg/validation"
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
	var details []commons.ApiErrorDetails
	for _, e := range errors {
		details = append(details, commons.ApiErrorDetails{Field: e.Field, Detail: e.Error})
	}
	return commons.ApiError{
		Code:    fiber.StatusBadRequest,
		Kind:    string(commons.ErrorTypeFunctional),
		Message: validation.GlobalValidationFailed,
		Details: details,
	}
}
