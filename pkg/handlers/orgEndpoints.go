package endpoints

import (
	"fmt"
	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	jsoniter "github.com/json-iterator/go"
	"micro-fiber-test/pkg/converters"
	dtos "micro-fiber-test/pkg/dto/commons"
	"micro-fiber-test/pkg/dto/orgs"
	"micro-fiber-test/pkg/exceptions"
	"micro-fiber-test/pkg/middlewares"
	"micro-fiber-test/pkg/model"
	"micro-fiber-test/pkg/service/api"
	"micro-fiber-test/pkg/validation"
)

var validate = validator.New()

func MakeOrgCreateEndpoint(rdbmsUrl string, defaultTenantId int64, orgSvc api.OrganizationServiceInterface) func(ctx *fiber.Ctx) error {
	return func(ctx *fiber.Ctx) error {

		var bsTxId interface{}
		bsTxId = ctx.Locals(middlewares.BsTxId)
		fmt.Printf("Bs Transaction id: [%s]", bsTxId)

		orgReq := orgs.CreateOrgRequest{}
		var json = jsoniter.ConfigCompatibleWithStandardLibrary
		if err := json.Unmarshal(ctx.Body(), &orgReq); err != nil {
			_ = ctx.SendStatus(fiber.StatusInternalServerError)
			apiErr := exceptions.ConvertToInternalError(err)
			return ctx.JSON(apiErr)
		}

		errValid := validate.Struct(orgReq)
		if errValid != nil {
			_ = ctx.SendStatus(fiber.StatusBadRequest)
			apiError := exceptions.ConvertValidationError(validation.ConvertValidationErrors(errValid))
			return ctx.JSON(apiError)
		}

		org := converters.ConvertOrgReqToDaoModel(defaultTenantId, orgReq)
		codeUUID := uuid.New().String()
		org.SetCode(codeUUID)
		switch org.GetStatus() {
		case model.OrgStatusDraft, model.OrgStatusActive, model.OrgStatusInactive, model.OrgStatusDeleted:
		default:
			_ = ctx.SendStatus(fiber.StatusBadRequest)
			apiErr := dtos.ApiError{
				Code:    fiber.StatusBadRequest,
				Kind:    string(dtos.ErrorTypeFunctional),
				Message: fmt.Sprintf("Invalid org status [%d]", orgReq.Status),
			}
			return ctx.JSON(apiErr)
		}
		_, err := orgSvc.Create(rdbmsUrl, defaultTenantId, &org)
		if err != nil {
			if err.Error() == dtos.OrgAlreadyExistsByCode {
				_ = ctx.SendStatus(fiber.StatusConflict)
				apiErr := exceptions.ConvertToFunctionalError(err, fiber.StatusConflict)
				return ctx.JSON(apiErr)
			} else {
				_ = ctx.SendStatus(fiber.StatusInternalServerError)
				apiErr := exceptions.ConvertToInternalError(err)
				return ctx.JSON(apiErr)
			}
		} else {
			_ = ctx.SendStatus(fiber.StatusCreated)
			idResponse := dtos.CodeResponse{Code: codeUUID}
			return ctx.JSON(idResponse)
		}

	}
}

func MakeOrgUpdateEndpoint(defaultTenantId int64, orgSvc api.OrganizationServiceInterface) func(ctx *fiber.Ctx) error {
	return func(ctx *fiber.Ctx) error {
		orgCode := ctx.Params("orgCode")
		payload := struct {
			Label string `json:"label"`
		}{}
		if err := ctx.BodyParser(&payload); err != nil {
			_ = ctx.SendStatus(fiber.StatusInternalServerError)
			apiErr := exceptions.ConvertToInternalError(err)
			return ctx.JSON(apiErr)
		}
		errUpdate := orgSvc.Update(defaultTenantId, orgCode, payload.Label)

		if errUpdate != nil {
			_ = ctx.SendStatus(fiber.StatusInternalServerError)
			apiErr := exceptions.ConvertToInternalError(errUpdate)
			return ctx.JSON(apiErr)
		} else {
			_ = ctx.SendStatus(fiber.StatusNoContent)
			return nil
		}

	}
}

func MakeOrgDeleteEndpoint(defaultTenantId int64, orgSvc api.OrganizationServiceInterface) func(ctx *fiber.Ctx) error {
	return func(ctx *fiber.Ctx) error {
		orgCode := ctx.Params("orgCode")
		_, errFind := orgSvc.FindByCode(defaultTenantId, orgCode)
		if errFind != nil {
			if errFind.Error() == dtos.OrgDoesNotExistByCode {
				_ = ctx.SendStatus(fiber.StatusNotFound)
				apiErr := exceptions.ConvertToFunctionalError(errFind, fiber.StatusNotFound)
				return ctx.JSON(apiErr)
			} else {
				_ = ctx.SendStatus(fiber.StatusInternalServerError)
				apiErr := exceptions.ConvertToInternalError(errFind)
				return ctx.JSON(apiErr)
			}
		} else {
			errUpdate := orgSvc.Delete(defaultTenantId, orgCode)
			if errUpdate != nil {
				_ = ctx.SendStatus(fiber.StatusInternalServerError)
				apiErr := exceptions.ConvertToInternalError(errUpdate)
				return ctx.JSON(apiErr)
			} else {
				_ = ctx.SendStatus(fiber.StatusNoContent)
				return nil
			}
		}
	}
}

func MakeOrgFindByCodeEndpoint(defaultTenantId int64, orgSvc api.OrganizationServiceInterface) func(ctx *fiber.Ctx) error {
	return func(ctx *fiber.Ctx) error {
		orgCode := ctx.Params("orgCode")
		org, errFind := orgSvc.FindByCode(defaultTenantId, orgCode)
		if errFind != nil {
			if errFind.Error() == dtos.OrgDoesNotExistByCode {
				_ = ctx.SendStatus(fiber.StatusNotFound)
				apiErr := exceptions.ConvertToFunctionalError(errFind, fiber.StatusNotFound)
				return ctx.JSON(apiErr)
			} else {
				_ = ctx.SendStatus(fiber.StatusInternalServerError)
				apiErr := exceptions.ConvertToInternalError(errFind)
				return ctx.JSON(apiErr)
			}
		} else {
			orgResponse := converters.ConvertOrgModelToOrgResp(org)
			ctx.GetRespHeader(fiber.HeaderContentType, fiber.MIMEApplicationJavaScriptCharsetUTF8)
			_ = ctx.SendStatus(fiber.StatusOK)
			return ctx.JSON(orgResponse)
		}
	}
}

func MakeOrgFindAll(defaultTenantId int64, orgSvc api.OrganizationServiceInterface) func(ctx *fiber.Ctx) error {
	return func(ctx *fiber.Ctx) error {
		orgsList, errFindAll := orgSvc.FindAll(defaultTenantId)
		if errFindAll != nil {
			_ = ctx.SendStatus(fiber.StatusInternalServerError)
			apiErr := exceptions.ConvertToInternalError(errFindAll)
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
			ctx.GetRespHeader(fiber.HeaderContentType, fiber.MIMEApplicationJavaScriptCharsetUTF8)
			_ = ctx.SendStatus(fiber.StatusOK)
			return ctx.JSON(orgListResponse)
		}
	}
}
