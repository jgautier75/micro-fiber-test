package endpoints

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	jsoniter "github.com/json-iterator/go"
	"micro-fiber-test/pkg/commons"
	"micro-fiber-test/pkg/contracts"
	"micro-fiber-test/pkg/converters"
	dtos "micro-fiber-test/pkg/dto/commons"
	"micro-fiber-test/pkg/dto/orgs"
	"micro-fiber-test/pkg/model"
	"micro-fiber-test/pkg/service/api"
	"micro-fiber-test/pkg/validation"
)

func MakeOrgCreateEndpoint(rdbmsUrl string, defaultTenantId int64, orgSvc api.OrganizationServiceInterface) func(ctx *fiber.Ctx) error {
	return func(ctx *fiber.Ctx) error {
		orgReq := orgs.CreateOrgRequest{}
		var json = jsoniter.ConfigCompatibleWithStandardLibrary
		if err := json.Unmarshal(ctx.Body(), &orgReq); err != nil {
			_ = ctx.SendStatus(fiber.StatusInternalServerError)
			apiErr := contracts.ConvertToInternalError(err)
			return ctx.JSON(apiErr)
		}

		validErr := validation.Validate(orgReq)
		if validErr != nil && len(validErr) > 0 {
			_ = ctx.SendStatus(fiber.StatusBadRequest)
			apiError := contracts.ConvertValidationError(validErr)
			return ctx.JSON(apiError)
		}

		org := converters.ConvertOrgReqToDaoModel(defaultTenantId, orgReq)
		codeUUID := uuid.New().String()
		org.SetCode(codeUUID)
		switch org.GetStatus() {
		case model.OrgStatusDraft, model.OrgStatusActive, model.OrgStatusInactive, model.OrgStatusDeleted:
		default:
			_ = ctx.SendStatus(fiber.StatusBadRequest)
			apiErr := commons.ApiError{
				Code:    fiber.StatusBadRequest,
				Kind:    string(commons.ErrorTypeFunctional),
				Message: fmt.Sprintf("Invalid org status [%d]", orgReq.Status),
			}
			return ctx.JSON(apiErr)
		}
		_, err := orgSvc.Create(rdbmsUrl, defaultTenantId, &org)
		if err != nil {
			if err.Error() == commons.OrgAlreadyExistsByCode {
				_ = ctx.SendStatus(fiber.StatusConflict)
				apiErr := contracts.ConvertToFunctionalError(err, fiber.StatusConflict)
				return ctx.JSON(apiErr)
			} else {
				_ = ctx.SendStatus(fiber.StatusInternalServerError)
				apiErr := contracts.ConvertToInternalError(err)
				return ctx.JSON(apiErr)
			}
		} else {
			_ = ctx.SendStatus(fiber.StatusCreated)
			idResponse := dtos.CodeResponse{Code: codeUUID}
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
			_ = ctx.SendStatus(fiber.StatusInternalServerError)
			apiErr := contracts.ConvertToInternalError(err)
			return ctx.JSON(apiErr)
		}
		errUpdate := orgSvc.Update(rdbmsUrl, defaultTenantId, orgCode, payload.Label)

		if errUpdate != nil {
			_ = ctx.SendStatus(fiber.StatusInternalServerError)
			apiErr := contracts.ConvertToInternalError(errUpdate)
			return ctx.JSON(apiErr)
		} else {
			_ = ctx.SendStatus(fiber.StatusNoContent)
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
				_ = ctx.SendStatus(fiber.StatusNotFound)
				apiErr := contracts.ConvertToFunctionalError(errFind, fiber.StatusNotFound)
				return ctx.JSON(apiErr)
			} else {
				_ = ctx.SendStatus(fiber.StatusInternalServerError)
				apiErr := contracts.ConvertToInternalError(errFind)
				return ctx.JSON(apiErr)
			}
		} else {
			errUpdate := orgSvc.Delete(rdbmsUrl, defaultTenantId, orgCode)
			if errUpdate != nil {
				_ = ctx.SendStatus(fiber.StatusInternalServerError)
				apiErr := contracts.ConvertToInternalError(errUpdate)
				return ctx.JSON(apiErr)
			} else {
				_ = ctx.SendStatus(fiber.StatusNoContent)
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
				_ = ctx.SendStatus(fiber.StatusNotFound)
				apiErr := contracts.ConvertToFunctionalError(errFind, fiber.StatusNotFound)
				return ctx.JSON(apiErr)
			} else {
				_ = ctx.SendStatus(fiber.StatusInternalServerError)
				apiErr := contracts.ConvertToInternalError(errFind)
				return ctx.JSON(apiErr)
			}
		} else {
			orgResponse := converters.ConvertOrgModelToOrgResp(org)
			ctx.GetRespHeader(commons.ContentTypeHeader, commons.ContentTypeJson)
			_ = ctx.SendStatus(fiber.StatusOK)
			return ctx.JSON(orgResponse)
		}
	}
}

func MakeOrgFindAll(dbmsUrl string, defaultTenantId int64, orgSvc api.OrganizationServiceInterface) func(ctx *fiber.Ctx) error {
	return func(ctx *fiber.Ctx) error {
		orgsList, errFindAll := orgSvc.FindAll(dbmsUrl, defaultTenantId)
		if errFindAll != nil {
			_ = ctx.SendStatus(fiber.StatusInternalServerError)
			apiErr := contracts.ConvertToInternalError(errFindAll)
			return ctx.JSON(apiErr)
		} else {
			orgResponseList := make([]orgs.OrganizationResponse, len(orgsList), len(orgsList))
			for inc, org := range orgsList {
				orgResponse := converters.ConvertOrgModelToOrgResp(org)
				orgResponseList[inc] = orgResponse
			}
			orgListResponse := orgs.OrganizationListResponse{
				Organizations: orgResponseList,
			}
			ctx.GetRespHeader(commons.ContentTypeHeader, commons.ContentTypeJson)
			_ = ctx.SendStatus(fiber.StatusOK)
			return ctx.JSON(orgListResponse)
		}
	}
}
