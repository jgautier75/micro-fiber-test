package endpoints

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"micro-fiber-test/pkg/commons"
	"micro-fiber-test/pkg/contracts"
	dtos "micro-fiber-test/pkg/dto/commons"
	"micro-fiber-test/pkg/dto/users"
	"micro-fiber-test/pkg/service/api"
	"micro-fiber-test/pkg/validation"
)

func MakeUserCreateEndpoint(dbmsUrl string, defaultTenantId int64, userSvc api.UserServiceInterface, orgSvc api.OrganizationServiceInterface) func(ctx *fiber.Ctx) error {
	return func(ctx *fiber.Ctx) error {
		orgCode := ctx.Params("orgCode")

		// Ensure organization exists
		org, errFindOrga := orgSvc.FindByCode(dbmsUrl, defaultTenantId, orgCode)
		if errFindOrga != nil {
			ctx.SendStatus(fiber.StatusInternalServerError)
			apiErr := contracts.ConvertToInternalError(errFindOrga)
			return ctx.JSON(apiErr)
		}
		if org == nil {
			ctx.SendStatus(fiber.StatusNotFound)
			apiErr := contracts.ConvertToFunctionalError(errors.New(commons.OrgNotFound), fiber.StatusNotFound)
			return ctx.JSON(apiErr)
		}

		// Deserialize request
		userReq := users.CreateUserReq{}
		if err := ctx.BodyParser(&userReq); err != nil {
			ctx.SendStatus(fiber.StatusInternalServerError)
			apiErr := contracts.ConvertToInternalError(err)
			return ctx.JSON(apiErr)
		}

		// Validate payload
		validErr := validation.Validate(userReq)
		if validErr != nil && len(validErr) > 0 {
			ctx.SendStatus(fiber.StatusBadRequest)
			apiError := contracts.ConvertValidationError(validErr)
			return ctx.JSON(apiError)
		}

		usrModel := users.ConvertUserReqToDaoModel(defaultTenantId, userReq)
		usrModel.SetOrgId(org.GetId())
		codeUUID := uuid.New().String()
		usrModel.SetCode(codeUUID)
		_, errCreate := userSvc.Create(dbmsUrl, defaultTenantId, &usrModel)
		if errCreate != nil {
			return errCreate
		}
		idResponse := dtos.CodeResponse{Code: codeUUID}
		ctx.SendStatus(fiber.StatusCreated)
		return ctx.JSON(idResponse)
	}
}
