package endpoints

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"micro-fiber-test/pkg/commons"
	"micro-fiber-test/pkg/contracts"
	"micro-fiber-test/pkg/converters"
	commonsDto "micro-fiber-test/pkg/dto/commons"
	dtos "micro-fiber-test/pkg/dto/commons"
	"micro-fiber-test/pkg/dto/users"
	usersResponses "micro-fiber-test/pkg/dto/users"
	"micro-fiber-test/pkg/model"
	"micro-fiber-test/pkg/service/api"
	"micro-fiber-test/pkg/validation"
	"strconv"
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

		usrModel := converters.ConvertUserReqToDaoModel(defaultTenantId, userReq)
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

func MakeUserSearchFilter(dbmsUrl string, defaultTenantId int64, userSvc api.UserServiceInterface, orgSvc api.OrganizationServiceInterface) func(ctx *fiber.Ctx) error {
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

		userFilterCriteria, errCriteria := buildCriteria(org, ctx)
		if errCriteria != nil {
			return errCriteria
		}

		usersCriteria, errFind := userSvc.FindByCriteria(dbmsUrl, userFilterCriteria)
		if errFind != nil {
			return errFind
		}

		usersArray := make([]usersResponses.UserResponse, len(usersCriteria.Users), len(usersCriteria.Users))
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

		userListReponse := usersResponses.UserListResponse{
			Users:      usersArray,
			Pagination: pageResp,
		}

		ctx.GetRespHeader(commons.ContentTypeHeader, commons.ContentTypeJson)
		ctx.SendStatus(fiber.StatusOK)
		return ctx.JSON(userListReponse)
	}
}

func buildCriteria(org model.OrganizationInterface, ctx *fiber.Ctx) (model.UserFilterCriteria, error) {
	userFilterCriteria := model.UserFilterCriteria{}
	userFilterCriteria.OrgId = org.GetId()
	userFilterCriteria.TenantId = org.GetTenantId()

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
