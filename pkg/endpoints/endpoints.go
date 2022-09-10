package endpoints

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"micro-fiber-test/pkg/commons"
	"micro-fiber-test/pkg/contracts"
	"micro-fiber-test/pkg/model"
	"micro-fiber-test/pkg/service/api"
)

func MakeOrgCreateEndpoint(rdbmsUrl string, defaultTenantId int64, orgSvc api.OrganizationServiceInterface) func(ctx *fiber.Ctx) error {
	return func(ctx *fiber.Ctx) error {
		payload := struct {
			Code   string `json:"code"`
			Label  string `json:"label"`
			Kind   string `json:"type"`
			Status int    `json:"status"`
		}{}
		if err := ctx.BodyParser(&payload); err != nil {
			fmt.Println("error = ", err)
			return ctx.SendStatus(200)
		}
		org := model.Organization{}
		org.SetTenantId(defaultTenantId)
		org.SetCode(payload.Code)
		org.SetLabel(payload.Label)
		org.SetType(model.OrganizationType(payload.Kind))
		org.SetStatus(model.OrganizationStatus(payload.Status))
		id, err := orgSvc.Create(rdbmsUrl, &org)
		if err != nil {
			if err.Error() == commons.OrgAlreadyExistsByCode {
				ctx.SendStatus(fiber.StatusConflict)
				apiErr := commons.ApiError{
					Code:    fiber.StatusConflict,
					Kind:    string(commons.ErrorTypeFunctional),
					Message: err.Error(),
				}
				return ctx.JSON(apiErr)
			} else {
				ctx.SendStatus(fiber.StatusInternalServerError)
				apiErr := commons.ApiError{
					Code:    fiber.StatusInternalServerError,
					Kind:    string(commons.ErrorTypeTechnical),
					Message: err.Error(),
				}
				return ctx.JSON(apiErr)
			}

		} else {
			ctx.SendStatus(fiber.StatusCreated)
			idResponse := contracts.IdResponse{ID: id}
			return ctx.JSON(idResponse)
		}

	}
}
