package users

import "micro-fiber-test/pkg/model"

type CreateUserReq struct {
	LastName   string  `json:"lastName" validate:"notblank,maxLength(50)"`
	FirstName  string  `json:"firstName" validate:"notblank,maxLength(50)"`
	MiddleName *string `json:"middleName"`
	Login      string  `json:"login" validate:"notblank,maxLength(50)"`
	Email      string  `json:"email" validate:"notblank,maxLength(50)"`
	Status     int     `json:"status"`
}

func ConvertUserReqToDaoModel(defaultTenantId int64, userReq CreateUserReq) model.User {
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
