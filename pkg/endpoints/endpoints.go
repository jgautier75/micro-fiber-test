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

func MakeOrgUpdateEndpoint(rdbmsUrl string, orgSvc api.OrganizationServiceInterface) func(ctx *fiber.Ctx) error {
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
		errUpdate := orgSvc.Update(rdbmsUrl, orgCode, payload.Label)

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

func MakeOrgDeleteEndpoint(rdbmsUrl string, orgSvc api.OrganizationServiceInterface) func(ctx *fiber.Ctx) error {
	return func(ctx *fiber.Ctx) error {
		orgCode := ctx.Params("orgCode")
		errUpdate := orgSvc.Delete(rdbmsUrl, orgCode)

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

func MakeOrgFindByCodeEndpoint(rdbmsUrl string, orgSvc api.OrganizationServiceInterface) func(ctx *fiber.Ctx) error {
	return func(ctx *fiber.Ctx) error {
		orgCode := ctx.Params("orgCode")
		org, errUpdate := orgSvc.FindByCode(rdbmsUrl, orgCode)

		if errUpdate != nil {
			ctx.SendStatus(fiber.StatusInternalServerError)
			apiErr := contracts.ConvertToInternalError(errUpdate)
			return ctx.JSON(apiErr)
		} else {
			payload := struct {
				Code   string `json:"code"`
				Label  string `json:"label"`
				Kind   string `json:"type"`
				Status int    `json:"status"`
			}{}
			payload.Code = org.GetCode()
			payload.Kind = string(org.GetType())
			payload.Label = org.GetLabel()
			payload.Status = int(org.GetStatus())
			ctx.GetRespHeader(commons.ContentTypeHeader, commons.ContentTypeJson)
			ctx.SendStatus(fiber.StatusOK)
			return ctx.JSON(payload)
		}
	}
}

func MakeOrgFindAll(dbmsUrl string, orgSvc api.OrganizationServiceInterface) func(ctx *fiber.Ctx) error {
	return func(ctx *fiber.Ctx) error {
		orgs, errFindAll := orgSvc.FindAll(dbmsUrl)
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