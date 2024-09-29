package converters

import (
	"micro-fiber-test/pkg/dto/users"
	"micro-fiber-test/pkg/model"
)

func ConvertUserReqToDaoModel(defaultTenantId int64, userReq users.CreateUserReq) model.User {
	usr := model.User{}
	usr.TenantId = defaultTenantId
	usr.LastName = userReq.LastName
	usr.FirstName = userReq.FirstName
	if userReq.MiddleName != nil {
		usr.MiddleName = *userReq.MiddleName
	}
	usr.Login = userReq.Login
	usr.Email = userReq.Email
	usr.Status = model.UserStatus(userReq.Status)
	return usr
}

func ConvertUserUpdateReqToDaoModel(defaultTenantId int64, userReq users.UpdateUserReq) model.User {
	usr := model.User{}
	usr.TenantId = defaultTenantId
	usr.LastName = userReq.LastName
	usr.FirstName = userReq.FirstName
	if userReq.MiddleName != nil {
		usr.MiddleName = *userReq.MiddleName
	}
	usr.Email = userReq.Email
	usr.Login = userReq.Login
	return usr
}

func ConvertFromDaoModelToUserResponse(userInterface model.User) users.UserResponse {
	usr := users.UserResponse{}
	usr.ExternalId = userInterface.ExternalId
	usr.Login = userInterface.Login
	usr.Email = userInterface.Email
	usr.LastName = userInterface.LastName
	usr.FirstName = userInterface.FirstName
	usr.MiddleName = userInterface.MiddleName
	usr.Status = int(userInterface.Status)
	return usr
}
