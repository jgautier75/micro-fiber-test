package endpoints

import (
	"database/sql"
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"micro-fiber-test/pkg/converters"
	dtos "micro-fiber-test/pkg/dto/commons"
	"micro-fiber-test/pkg/dto/orgs"
	"micro-fiber-test/pkg/dto/sectors"
	"micro-fiber-test/pkg/exceptions"
	"micro-fiber-test/pkg/helpers"
	"micro-fiber-test/pkg/service/api"
	"micro-fiber-test/pkg/validation"
)

func MakeSectorsFindByOrga(defaultTenantId int64, orgSvc api.OrganizationServiceInterface, sectSvc api.SectorServiceInterface) func(ctx *fiber.Ctx) error {
	return func(ctx *fiber.Ctx) error {
		orgCode := ctx.Params("orgCode")
		org, errFindOrga := orgSvc.FindByCode(defaultTenantId, orgCode)
		if errFindOrga != nil {
			if errFindOrga.Error() == dtos.OrgDoesNotExistByCode {
				_ = ctx.SendStatus(fiber.StatusNotFound)
				apiErr := exceptions.ConvertToFunctionalError(errFindOrga, fiber.StatusNotFound)
				return ctx.JSON(apiErr)
			} else {
				_ = ctx.SendStatus(fiber.StatusInternalServerError)
				apiErr := exceptions.ConvertToInternalError(errFindOrga)
				return ctx.JSON(apiErr)
			}
		}
		if org == nil {
			_ = ctx.SendStatus(fiber.StatusNotFound)
			apiErr := exceptions.ConvertToFunctionalError(errors.New(dtos.OrgNotFound), fiber.StatusNotFound)
			return ctx.JSON(apiErr)
		}

		sectorsList, errFindAll := sectSvc.FindSectorsByTenantOrg(defaultTenantId, org.GetId())
		if errFindAll != nil {
			_ = ctx.SendStatus(fiber.StatusInternalServerError)
			apiErr := exceptions.ConvertToInternalError(errFindAll)
			return ctx.JSON(apiErr)
		} else {
			sectorsResponseList := make([]sectors.SectorResponse, len(sectorsList), len(sectorsList))
			for inc, s := range sectorsList {
				sgResponse := converters.ConvertSectorModelToSectorResp(s)
				sectorsResponseList[inc] = sgResponse
			}
			s, errHierarchy := helpers.BuildSectorsHierarchy(sectorsResponseList)
			if errHierarchy != nil {
				return nil
			}
			sectListResponse := sectors.SectorListResponse{
				Sectors: s,
			}
			ctx.GetRespHeader(fiber.HeaderContentType, fiber.MIMEApplicationJavaScriptCharsetUTF8)
			_ = ctx.SendStatus(fiber.StatusOK)
			return ctx.JSON(sectListResponse)
		}
	}
}

func MakeSectorCreateEndpoint(defaultTenantId int64, orgSvc api.OrganizationServiceInterface, sectSvc api.SectorServiceInterface) func(ctx *fiber.Ctx) error {
	return func(ctx *fiber.Ctx) error {
		orgCode := ctx.Params("orgCode")

		// Ensure organization exists
		org, errFindOrga := orgSvc.FindByCode(defaultTenantId, orgCode)
		if errFindOrga != nil {
			_ = ctx.SendStatus(fiber.StatusInternalServerError)
			apiErr := exceptions.ConvertToInternalError(errFindOrga)
			return ctx.JSON(apiErr)
		}
		if org == nil {
			_ = ctx.SendStatus(fiber.StatusNotFound)
			apiErr := exceptions.ConvertToFunctionalError(errors.New(dtos.OrgNotFound), fiber.StatusNotFound)
			return ctx.JSON(apiErr)
		}

		// Deserialize request
		sectorReq := orgs.CreateSectorReq{}
		if err := ctx.BodyParser(&sectorReq); err != nil {
			_ = ctx.SendStatus(fiber.StatusInternalServerError)
			apiErr := exceptions.ConvertToInternalError(err)
			return ctx.JSON(apiErr)
		}

		// Validate payload
		errValid := validate.Struct(sectorReq)
		if errValid != nil {
			_ = ctx.SendStatus(fiber.StatusBadRequest)
			apiError := exceptions.ConvertValidationError(validation.ConvertValidationErrors(errValid))
			return ctx.JSON(apiError)
		}

		secModel := orgs.ConvertSectorReqToDaoModel(defaultTenantId, sectorReq)
		secModel.SetOrgId(org.GetId())
		secModel.SetHasParent(true)
		codeUUID := uuid.New().String()
		secModel.SetCode(codeUUID)
		if sectorReq.ParentCode != "" {
			// Find parent sector
			parentSector, errParent := sectSvc.FindByCode(defaultTenantId, sectorReq.ParentCode)
			if errParent != nil {
				return errParent
			}
			if parentSector != nil {
				nillableInt64 := sql.NullInt64{
					Int64: parentSector.GetId(),
					Valid: true,
				}
				secModel.SetParentId(nillableInt64)
				secModel.SetDepth(parentSector.GetDepth() + 1)
			}
		} else {
			// If parent sector not set, inherits from root
			rootSector, err := sectSvc.FindRootSectorId(defaultTenantId, org.GetId())
			if err != nil {
				return err
			}
			nillableInt64 := sql.NullInt64{
				Int64: rootSector,
				Valid: true,
			}
			secModel.SetParentId(nillableInt64)
			secModel.SetDepth(1)
		}

		_, errCreate := sectSvc.Create(defaultTenantId, &secModel)
		if errCreate != nil {
			if errCreate.Error() == dtos.SectorAlreadyExist {
				_ = ctx.SendStatus(fiber.StatusConflict)
				apiError := exceptions.ConvertToFunctionalError(errCreate, fiber.StatusConflict)
				return ctx.JSON(apiError)
			} else {
				_ = ctx.SendStatus(fiber.StatusInternalServerError)
				apiError := exceptions.ConvertToInternalError(errCreate)
				return ctx.JSON(apiError)
			}
		}

		idResponse := dtos.CodeResponse{Code: codeUUID}
		_ = ctx.SendStatus(fiber.StatusCreated)
		return ctx.JSON(idResponse)
	}
}

func MakeSectorDeleteEndpoint(defaultTenantId int64, orgSvc api.OrganizationServiceInterface, sectSvc api.SectorServiceInterface) func(ctx *fiber.Ctx) error {
	return func(ctx *fiber.Ctx) error {
		orgCode := ctx.Params("orgCode")

		// Ensure organization exists
		org, errFindOrga := orgSvc.FindByCode(defaultTenantId, orgCode)
		if errFindOrga != nil {
			_ = ctx.SendStatus(fiber.StatusInternalServerError)
			apiErr := exceptions.ConvertToInternalError(errFindOrga)
			return ctx.JSON(apiErr)
		}
		if org == nil {
			_ = ctx.SendStatus(fiber.StatusNotFound)
			apiErr := exceptions.ConvertToFunctionalError(errors.New(dtos.OrgNotFound), fiber.StatusNotFound)
			return ctx.JSON(apiErr)
		}

		// Ensure sector exists
		sectorCode := ctx.Params("sectorCode")
		sector, errSect := sectSvc.FindByCode(defaultTenantId, sectorCode)
		if errSect != nil {
			return errSect
		}
		if sector == nil || sector.GetId() <= 0 {
			_ = ctx.SendStatus(fiber.StatusNotFound)
			apiErr := exceptions.ConvertToFunctionalError(errors.New(dtos.SectorNotFound), fiber.StatusNotFound)
			return ctx.JSON(apiErr)
		}

		errDelete := sectSvc.DeleteSector(defaultTenantId, sector.GetId())
		if errDelete != nil {
			_ = ctx.SendStatus(fiber.StatusInternalServerError)
			apiErr := exceptions.ConvertToInternalError(errDelete)
			return ctx.JSON(apiErr)
		}

		_ = ctx.SendStatus(fiber.StatusNoContent)
		return nil
	}
}

func MakeSectorUpdateEndpoint(defaultTenantId int64, orgSvc api.OrganizationServiceInterface, sectSvc api.SectorServiceInterface) func(ctx *fiber.Ctx) error {
	return func(ctx *fiber.Ctx) error {
		orgCode := ctx.Params("orgCode")

		// Ensure organization exists
		org, errFindOrga := orgSvc.FindByCode(defaultTenantId, orgCode)
		if errFindOrga != nil {
			_ = ctx.SendStatus(fiber.StatusInternalServerError)
			apiErr := exceptions.ConvertToInternalError(errFindOrga)
			return ctx.JSON(apiErr)
		}
		if org == nil {
			_ = ctx.SendStatus(fiber.StatusNotFound)
			apiErr := exceptions.ConvertToFunctionalError(errors.New(dtos.OrgNotFound), fiber.StatusNotFound)
			return ctx.JSON(apiErr)
		}

		// Ensure sector exists
		sectorCode := ctx.Params("sectorCode")
		sector, errSect := sectSvc.FindByCode(defaultTenantId, sectorCode)
		if errSect != nil {
			return errSect
		}
		if sector == nil || sector.GetId() <= 0 {
			_ = ctx.SendStatus(fiber.StatusNotFound)
			apiErr := exceptions.ConvertToFunctionalError(errors.New(dtos.SectorNotFound), fiber.StatusNotFound)
			return ctx.JSON(apiErr)
		}

		payload := struct {
			Label string `json:"label"`
		}{}
		if err := ctx.BodyParser(&payload); err != nil {
			_ = ctx.SendStatus(fiber.StatusInternalServerError)
			apiErr := exceptions.ConvertToInternalError(err)
			return ctx.JSON(apiErr)
		}

		idSec, codeSec, errSec := sectSvc.FindByLabel(defaultTenantId, payload.Label)
		if errSec != nil {
			_ = ctx.SendStatus(fiber.StatusInternalServerError)
			apiErr := exceptions.ConvertToInternalError(errSec)
			return ctx.JSON(apiErr)
		}

		if idSec > 0 && codeSec != sector.GetCode() {
			_ = ctx.SendStatus(fiber.StatusConflict)
			apiErr := exceptions.ConvertToFunctionalError(errors.New(dtos.SectorAlreadyExist), fiber.StatusConflict)
			return ctx.JSON(apiErr)
		}

		errDelete := sectSvc.Update(defaultTenantId, sector.GetId(), payload.Label)
		if errDelete != nil {
			_ = ctx.SendStatus(fiber.StatusInternalServerError)
			apiErr := exceptions.ConvertToInternalError(errDelete)
			return ctx.JSON(apiErr)
		}

		_ = ctx.SendStatus(fiber.StatusNoContent)
		return nil
	}
}
