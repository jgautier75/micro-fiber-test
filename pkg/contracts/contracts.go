package contracts

import (
	"github.com/gofiber/fiber"
	"micro-fiber-test/pkg/commons"
)

func ConvertToInternalError(err error) commons.ApiError {
	return commons.ApiError{
		Code:    fiber.StatusInternalServerError,
		Kind:    string(commons.ErrorTypeTechnical),
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
