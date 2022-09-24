package endpoints

import (
	"github.com/gofiber/fiber/v2"
	jsoniter "github.com/json-iterator/go"
	"micro-fiber-test/pkg/commons"
	"micro-fiber-test/pkg/contracts"
	"micro-fiber-test/pkg/model"
	"micro-fiber-test/pkg/service/api"
	"micro-fiber-test/pkg/validation"
)

func MakeOrgCreateEndpoint(rdbmsUrl string, defaultTenantId int64, orgSvc api.OrganizationServiceInterface) func(ctx *fiber.Ctx) error {
	return func(ctx *fiber.Ctx) error {
		payload := struct {
			Code   *string `json:"code"`
			Label  *string `json:"label"`
			Kind   *string `json:"type"`
			Status int     `json:"status"`
		}{}
		var json = jsoniter.ConfigCompatibleWithStandardLibrary
		if err := json.Unmarshal(ctx.Body(), &payload); err != nil {
			ctx.SendStatus(fiber.StatusInternalServerError)
			apiErr := contracts.ConvertToInternalError(err)
			return ctx.JSON(apiErr)
		}

		validErr := validation.Validate(payload)
		if validErr != nil && len(validErr) > 0 {
			apiError := contracts.ConvertValidationError(validErr)
			ctx.SendStatus(fiber.StatusBadRequest)
			return ctx.JSON(apiError)
		}

		org := model.Organization{}
		org.SetTenantId(defaultTenantId)
		if payload.Code != nil {
			org.SetCode(*payload.Code)
		}
		if payload.Label != nil {
			org.SetLabel(*payload.Label)
		}
		if payload.Kind != nil {
			org.SetType(model.OrganizationType(*payload.Kind))
		}
		org.SetStatus(model.OrganizationStatus(payload.Status))
		id, err := orgSvc.Create(rdbmsUrl, defaultTenantId, &org)
		if err != nil {
			if err.Error() == commons.OrgAlreadyExistsByCode {
				ctx.SendStatus(fiber.StatusConflict)
				apiErr := contracts.ConvertToInternalError(err)
				return ctx.JSON(apiErr)
			} else {
				ctx.SendStatus(fiber.StatusInternalServerError)
				apiErr := contracts.ConvertToInternalError(err)
				return ctx.JSON(apiErr)
			}
		} else {
			ctx.SendStatus(fiber.StatusCreated)
			idResponse := contracts.IdResponse{ID: id}
			return ctx.JSON(idResponse)
		}

	}
}

func MakeOrgUpdateEndpoint(rdbmsUrl string, defaultTenantId int64, orgSvc api.OrganizationServiceInterface) func(ctx *fiber.Ctx) error {
	return func(ctx *fiber.Ctx) error {
		orgCode := ctx.Params("orgCode")
		payload := struct {
			Label string `json:"label"`
		}{}
		if err := ctx.BodyParser(&payload); err != nil {
			ctx.SendStatus(fiber.StatusInternalServerError)
			apiErr := contracts.ConvertToInternalError(err)
			return ctx.JSON(apiErr)
		}
		errUpdate := orgSvc.Update(rdbmsUrl, defaultTenantId, orgCode, payload.Label)

		if errUpdate != nil {
			ctx.SendStatus(fiber.StatusInternalServerError)
			apiErr := contracts.ConvertToInternalError(errUpdate)
			return ctx.JSON(apiErr)
		} else {
			ctx.SendStatus(fiber.StatusNoContent)
			return nil
		}

	}
}

func MakeOrgDeleteEndpoint(rdbmsUrl string, defaultTenantId int64, orgSvc api.OrganizationServiceInterface) func(ctx *fiber.Ctx) error {
	return func(ctx *fiber.Ctx) error {
		orgCode := ctx.Params("orgCode")
		_, errFind := orgSvc.FindByCode(rdbmsUrl, defaultTenantId, orgCode)
		if errFind != nil {
			if errFind.Error() == commons.OrgDoesNotExistByCode {
				ctx.SendStatus(fiber.StatusNotFound)
				apiErr := contracts.ConvertToFunctionalError(errFind, fiber.StatusNotFound)
				return ctx.JSON(apiErr)
			} else {
				ctx.SendStatus(fiber.StatusInternalServerError)
				apiErr := contracts.ConvertToInternalError(errFind)
				return ctx.JSON(apiErr)
			}
		} else {
			errUpdate := orgSvc.Delete(rdbmsUrl, defaultTenantId, orgCode)
			if errUpdate != nil {
				ctx.SendStatus(fiber.StatusInternalServerError)
				apiErr := contracts.ConvertToInternalError(errUpdate)
				return ctx.JSON(apiErr)
			} else {
				ctx.SendStatus(fiber.StatusNoContent)
				return nil
			}
		}
	}
}

func MakeOrgFindByCodeEndpoint(rdbmsUrl string, defaultTenantId int64, orgSvc api.OrganizationServiceInterface) func(ctx *fiber.Ctx) error {
	return func(ctx *fiber.Ctx) error {
		orgCode := ctx.Params("orgCode")
		org, errFind := orgSvc.FindByCode(rdbmsUrl, defaultTenantId, orgCode)
		if errFind != nil {
			if errFind.Error() == commons.OrgDoesNotExistByCode {
				ctx.SendStatus(fiber.StatusNotFound)
				apiErr := contracts.ConvertToFunctionalError(errFind, fiber.StatusNotFound)
				return ctx.JSON(apiErr)
			} else {
				ctx.SendStatus(fiber.StatusInternalServerError)
				apiErr := contracts.ConvertToInternalError(errFind)
				return ctx.JSON(apiErr)
			}
		} else {
			orgResponse := contracts.OrganizationResponse{
				Code:   org.GetCode(),
				Label:  org.GetLabel(),
				Status: int(org.GetStatus()),
				Kind:   string(org.GetType()),
			}
			ctx.GetRespHeader(commons.ContentTypeHeader, commons.ContentTypeJson)
			ctx.SendStatus(fiber.StatusOK)
			return ctx.JSON(orgResponse)
		}
	}
}

func MakeOrgFindAll(dbmsUrl string, defaultTenantId int64, orgSvc api.OrganizationServiceInterface) func(ctx *fiber.Ctx) error {
	return func(ctx *fiber.Ctx) error {
		orgs, errFindAll := orgSvc.FindAll(dbmsUrl, defaultTenantId)
		if errFindAll != nil {
			ctx.SendStatus(fiber.StatusInternalServerError)
			apiErr := contracts.ConvertToInternalError(errFindAll)
			return ctx.JSON(apiErr)
		} else {
			orgList := make([]contracts.OrganizationResponse, len(orgs), len(orgs))
			for inc, org := range orgs {
				orgResponse := contracts.OrganizationResponse{
					Code:   org.GetCode(),
					Label:  org.GetLabel(),
					Status: int(org.GetStatus()),
					Kind:   string(org.GetType()),
				}
				orgList[inc] = orgResponse
			}
			orgListResponse := contracts.OrganizationListResponse{
				Organizations: orgList,
			}
			ctx.GetRespHeader(commons.ContentTypeHeader, commons.ContentTypeJson)
			ctx.SendStatus(fiber.StatusOK)
			return ctx.JSON(orgListResponse)
		}
	}
}
