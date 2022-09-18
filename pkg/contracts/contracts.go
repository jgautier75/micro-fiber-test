package contracts

import (
	"github.com/gofiber/fiber/v2"
	"micro-fiber-test/pkg/commons"
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

type IdResponse struct {
	ID int64 `json:"id"`
}

type OrganizationListResponse struct {
	Organizations []OrganizationResponse `json:"organizations"`
}

type OrganizationResponse struct {
	Code   string `json:"code"`
	Label  string `json:"label"`
	Kind   string `json:"type"`
	Status int    `json:"status"`
}
