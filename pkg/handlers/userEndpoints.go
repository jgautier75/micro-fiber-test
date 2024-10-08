package endpoints

import (
	"errors"
	"micro-fiber-test/pkg/converters"
	commonsDto "micro-fiber-test/pkg/dto/commons"
	"micro-fiber-test/pkg/dto/users"
	"micro-fiber-test/pkg/exceptions"
	"micro-fiber-test/pkg/model"
	"micro-fiber-test/pkg/service/api"
	"micro-fiber-test/pkg/validation"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func MakeUserCreateEndpoint(defaultTenantId int64, userSvc api.UserServiceInterface, orgSvc api.OrganizationServiceInterface) func(ctx *fiber.Ctx) error {
	return func(ctx *fiber.Ctx) error {

		var nilOrg model.Organization

		orgCode := ctx.Params("orgCode")

		// Ensure organization exists
		org, errFindOrga := orgSvc.FindByCode(defaultTenantId, orgCode)
		if errFindOrga != nil {
			_ = ctx.SendStatus(fiber.StatusInternalServerError)
			apiErr := exceptions.ConvertToInternalError(errFindOrga)
			return ctx.JSON(apiErr)
		}
		if org == nilOrg {
			_ = ctx.SendStatus(fiber.StatusNotFound)
			apiErr := exceptions.ConvertToFunctionalError(errors.New(commonsDto.OrgNotFound), fiber.StatusNotFound)
			return ctx.JSON(apiErr)
		}

		// Deserialize request
		userReq := users.CreateUserReq{}
		if err := ctx.BodyParser(&userReq); err != nil {
			_ = ctx.SendStatus(fiber.StatusInternalServerError)
			apiErr := exceptions.ConvertToInternalError(err)
			return ctx.JSON(apiErr)
		}

		// Validate payload
		errValid := validate.Struct(userReq)
		if errValid != nil {
			_ = ctx.SendStatus(fiber.StatusBadRequest)
			apiError := exceptions.ConvertValidationError(validation.ConvertValidationErrors(errValid))
			return ctx.JSON(apiError)
		}

		usrModel := converters.ConvertUserReqToDaoModel(defaultTenantId, userReq)
		usrModel.OrgId = org.Id
		extUUID := uuid.New().String()
		usrModel.ExternalId = extUUID
		_, errCreate := userSvc.Create(defaultTenantId, usrModel)
		if errCreate != nil {
			if errCreate.Error() == commonsDto.UserLoginAlreadyInUse || errCreate.Error() == commonsDto.UserEmailAlreadyInUse {
				apiError := exceptions.ConvertToFunctionalError(errCreate, fiber.StatusConflict)
				_ = ctx.SendStatus(fiber.StatusConflict)
				return ctx.JSON(apiError)
			} else {
				apiError := exceptions.ConvertToInternalError(errCreate)
				_ = ctx.SendStatus(fiber.StatusInternalServerError)
				return ctx.JSON(apiError)
			}
		}
		idResponse := commonsDto.ExternalIdResponse{ID: extUUID}
		_ = ctx.SendStatus(fiber.StatusCreated)
		return ctx.JSON(idResponse)
	}
}

func MakeUserSearchFilter(defaultTenantId int64, userSvc api.UserServiceInterface, orgSvc api.OrganizationServiceInterface) func(ctx *fiber.Ctx) error {
	return func(ctx *fiber.Ctx) error {
		orgCode := ctx.Params("orgCode")

		// Ensure organization exists
		org, errFindOrga := orgSvc.FindByCode(defaultTenantId, orgCode)
		if errFindOrga != nil {
			if errFindOrga.Error() == commonsDto.OrgDoesNotExistByCode {
				_ = ctx.SendStatus(fiber.StatusNotFound)
				apiErr := exceptions.ConvertToFunctionalError(errors.New(commonsDto.OrgNotFound), fiber.StatusNotFound)
				return ctx.JSON(apiErr)
			} else {
				_ = ctx.SendStatus(fiber.StatusInternalServerError)
				apiErr := exceptions.ConvertToInternalError(errFindOrga)
				return ctx.JSON(apiErr)
			}
		}

		userFilterCriteria, errCriteria := buildCriteria(org, ctx)
		if errCriteria != nil {
			return errCriteria
		}

		usersCriteria, errFind := userSvc.FindByCriteria(userFilterCriteria)
		if errFind != nil {
			return errFind
		}

		usersArray := make([]users.UserResponse, len(usersCriteria.Users))
		for inc, u := range usersCriteria.Users {
			usrResponse := converters.ConvertFromDaoModelToUserResponse(u)
			usersArray[inc] = usrResponse
		}

		pageResp := commonsDto.Pagination{
			Page:       userFilterCriteria.Page,
			TotalCount: usersCriteria.NbResults,
		}
		totalPages := 1
		if usersCriteria.NbResults >= userFilterCriteria.RowsPerPage {
			totalPages = usersCriteria.NbResults/userFilterCriteria.RowsPerPage + 1
		}
		pageResp.NbPages = totalPages

		userListReponse := users.UserListResponse{
			Users:      usersArray,
			Pagination: pageResp,
		}

		_ = ctx.SendStatus(fiber.StatusOK)
		return ctx.JSON(userListReponse)
	}
}

func MakeUserFindByCode(defaultTenantId int64, userSvc api.UserServiceInterface, orgSvc api.OrganizationServiceInterface) func(ctx *fiber.Ctx) error {
	return func(ctx *fiber.Ctx) error {
		var nilOrg model.Organization
		var nilUser model.User

		orgCode := ctx.Params("orgCode")

		// Ensure organization exists
		org, errFindOrga := orgSvc.FindByCode(defaultTenantId, orgCode)
		if errFindOrga != nil {
			_ = ctx.SendStatus(fiber.StatusInternalServerError)
			apiErr := exceptions.ConvertToInternalError(errFindOrga)
			return ctx.JSON(apiErr)
		}
		if org == nilOrg {
			_ = ctx.SendStatus(fiber.StatusNotFound)
			apiErr := exceptions.ConvertToFunctionalError(errors.New(commonsDto.OrgNotFound), fiber.StatusNotFound)
			return ctx.JSON(apiErr)
		}

		usrId := ctx.Params("userId")
		u, errFind := userSvc.FindByCode(defaultTenantId, org.Id, usrId)
		if errFind != nil {
			return errFind
		}
		if u != nilUser {
			_ = ctx.SendStatus(fiber.StatusOK)
			return ctx.JSON(converters.ConvertFromDaoModelToUserResponse(u))
		} else {
			apiError := exceptions.ConvertToFunctionalError(errors.New(commonsDto.UserNotFound), fiber.StatusNotFound)
			return ctx.JSON(apiError)
		}
	}
}

func MakeUserDelete(defaultTenantId int64, userSvc api.UserServiceInterface, orgSvc api.OrganizationServiceInterface) func(ctx *fiber.Ctx) error {
	return func(ctx *fiber.Ctx) error {

		var nilOrg model.Organization
		var nilUser model.User

		orgCode := ctx.Params("orgCode")

		// Ensure organization exists
		org, errFindOrga := orgSvc.FindByCode(defaultTenantId, orgCode)
		if errFindOrga != nil {
			_ = ctx.SendStatus(fiber.StatusInternalServerError)
			apiErr := exceptions.ConvertToInternalError(errFindOrga)
			return ctx.JSON(apiErr)
		}
		if org == nilOrg {
			_ = ctx.SendStatus(fiber.StatusNotFound)
			apiErr := exceptions.ConvertToFunctionalError(errors.New(commonsDto.OrgNotFound), fiber.StatusNotFound)
			return ctx.JSON(apiErr)
		}

		usrId := ctx.Params("userId")
		u, errFind := userSvc.FindByCode(defaultTenantId, org.Id, usrId)
		if errFind != nil {
			return errFind
		}
		if u == nilUser {
			_ = ctx.SendStatus(fiber.StatusNotFound)
			apiErr := exceptions.ConvertToFunctionalError(errors.New(commonsDto.UserNotFound), fiber.StatusNotFound)
			return ctx.JSON(apiErr)
		}

		errDel := userSvc.Delete(usrId)
		if errDel != nil {
			_ = ctx.SendStatus(fiber.StatusInternalServerError)
			apiErr := exceptions.ConvertToInternalError(errFindOrga)
			return ctx.JSON(apiErr)
		}
		_ = ctx.SendStatus(fiber.StatusNoContent)
		return nil
	}
}

func MakeUserUpdate(defaultTenantId int64, userSvc api.UserServiceInterface, orgSvc api.OrganizationServiceInterface) func(ctx *fiber.Ctx) error {
	return func(ctx *fiber.Ctx) error {

		var nilOrg model.Organization
		var nilUser model.User

		orgCode := ctx.Params("orgCode")

		// Ensure organization exists
		org, errFindOrga := orgSvc.FindByCode(defaultTenantId, orgCode)
		if errFindOrga != nil {
			_ = ctx.SendStatus(fiber.StatusInternalServerError)
			apiErr := exceptions.ConvertToInternalError(errFindOrga)
			return ctx.JSON(apiErr)
		}
		if org == nilOrg {
			_ = ctx.SendStatus(fiber.StatusNotFound)
			apiErr := exceptions.ConvertToFunctionalError(errors.New(commonsDto.OrgNotFound), fiber.StatusNotFound)
			return ctx.JSON(apiErr)
		}

		usrId := ctx.Params("userId")
		u, errFind := userSvc.FindByCode(defaultTenantId, org.Id, usrId)
		if errFind != nil {
			return errFind
		}

		if u == nilUser {
			apiError := exceptions.ConvertToFunctionalError(errors.New(commonsDto.UserNotFound), fiber.StatusNotFound)
			_ = ctx.SendStatus(fiber.StatusNotFound)
			return ctx.JSON(apiError)
		}

		// Deserialize request
		userReq := users.UpdateUserReq{}
		if err := ctx.BodyParser(&userReq); err != nil {
			_ = ctx.SendStatus(fiber.StatusInternalServerError)
			apiErr := exceptions.ConvertToInternalError(err)
			return ctx.JSON(apiErr)
		}

		// Validate payload
		errValid := validate.Struct(userReq)
		if errValid != nil {
			_ = ctx.SendStatus(fiber.StatusBadRequest)
			apiError := exceptions.ConvertValidationError(validation.ConvertValidationErrors(errValid))
			return ctx.JSON(apiError)
		}

		usrModel := converters.ConvertUserUpdateReqToDaoModel(defaultTenantId, userReq)
		usrModel.OrgId = org.Id
		usrModel.ExternalId = usrId

		errUpdate := userSvc.Update(usrModel)
		if errUpdate != nil {
			if errUpdate.Error() == commonsDto.UserLoginAlreadyInUse || errUpdate.Error() == commonsDto.UserEmailAlreadyInUse {
				apiError := exceptions.ConvertToFunctionalError(errUpdate, fiber.StatusConflict)
				_ = ctx.SendStatus(fiber.StatusConflict)
				return ctx.JSON(apiError)
			} else {
				_ = ctx.SendStatus(fiber.StatusInternalServerError)
				apiErr := exceptions.ConvertToInternalError(errUpdate)
				return ctx.JSON(apiErr)
			}
		}

		return ctx.SendStatus(fiber.StatusNoContent)
	}
}

func buildCriteria(org model.Organization, ctx *fiber.Ctx) (model.UserFilterCriteria, error) {
	userFilterCriteria := model.UserFilterCriteria{}
	userFilterCriteria.OrgId = org.Id
	userFilterCriteria.TenantId = org.TenantId

	firstName := ctx.Query("firstname", "")
	userFilterCriteria.FirstName = firstName
	lastName := ctx.Query("lastname", "")
	userFilterCriteria.LastName = lastName
	email := ctx.Query("email", "")
	userFilterCriteria.Email = email
	login := ctx.Query("login", "")
	userFilterCriteria.Login = login
	rowsPerPageStr := ctx.Query("rows", "5")
	rowsPerPage, errConvert := strconv.Atoi(rowsPerPageStr)
	if errConvert != nil {
		return userFilterCriteria, errConvert
	}
	userFilterCriteria.RowsPerPage = rowsPerPage
	pageStr := ctx.Query("page", "1")
	curPage, errPage := strconv.Atoi(pageStr)
	if errPage != nil {
		return userFilterCriteria, errPage
	}
	userFilterCriteria.Page = curPage
	return userFilterCriteria, nil
}
