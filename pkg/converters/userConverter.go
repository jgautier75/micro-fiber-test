package converters

import (
	"micro-fiber-test/pkg/dto/users"
	"micro-fiber-test/pkg/model"
)

func ConvertUserReqToDaoModel(defaultTenantId int64, userReq users.CreateUserReq) model.User {
	usr := model.User{}
	usr.SetTenantId(defaultTenantId)
	usr.SetLastName(userReq.LastName)
	usr.SetFirstName(userReq.FirstName)
	if userReq.MiddleName != nil {
		usr.SetMiddleName(*userReq.MiddleName)
	}
	usr.SetLogin(userReq.Login)
	usr.SetEmail(userReq.Email)
	usr.SetStatus(model.UserStatus(userReq.Status))
	return usr
}

func ConvertUserUpdateReqToDaoModel(defaultTenantId int64, userReq users.UpdateUserReq) model.User {
	usr := model.User{}
	usr.SetTenantId(defaultTenantId)
	usr.SetLastName(userReq.LastName)
	usr.SetFirstName(userReq.FirstName)
	if userReq.MiddleName != nil {
		usr.SetMiddleName(*userReq.MiddleName)
	}
	return usr
}

func ConvertFromDaoModelToUserResponse(userInterface model.UserInterface) users.UserResponse {
	usr := users.UserResponse{}
	usr.ExternalId = userInterface.GetExternalId()
	usr.Login = userInterface.GetLogin()
	usr.Email = userInterface.GetEmail()
	usr.LastName = userInterface.GetLastName()
	usr.FirstName = userInterface.GetFirstName()
	usr.MiddleName = userInterface.GetMiddleName()
	usr.Status = int(userInterface.GetStatus())
	return usr
}
