package model

type UserStatus int64

const (
	UserStatusDraft    UserStatus = 0
	UserStatusActive   UserStatus = 1
	UserStatusInactive UserStatus = 2
	UserStatusDeleted  UserStatus = 3
)

type UserInterface interface {
	GetId() int64
	SetId(id int64)
	GetTenantId() int64
	SetTenantId(id int64)
	GetOrgId() int64
	SetOrgId(orgId int64)
	GetExternalId() string
	SetExternalId(code string)
	GetLastName() string
	SetLastName(lastName string)
	GetFirstName() string
	SetFirstName(firstName string)
	GetMiddleName() string
	SetMiddleName(middleName string)
	GetLogin() string
	SetLogin(login string)
	GetEmail() string
	SetEmail(email string)
	GetStatus() UserStatus
	SetStatus(status UserStatus)
}
