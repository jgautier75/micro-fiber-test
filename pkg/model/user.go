package model

type User struct {
	id         int64
	tenantId   int64
	orgId      int64
	code       string
	lastName   string
	firstName  string
	middleName string
	login      string
	email      string
	status     UserStatus
}

func (usr *User) GetId() int64 {
	return usr.id
}

func (usr *User) SetId(fid int64) {
	usr.id = fid
}

func (usr *User) GetTenantId() int64 {
	return usr.tenantId
}

func (usr *User) SetTenantId(ftenantId int64) {
	usr.tenantId = ftenantId
}

func (usr *User) GetOrgId() int64 {
	return usr.orgId
}

func (usr *User) SetOrgId(forgId int64) {
	usr.orgId = forgId
}

func (usr *User) GetCode() string {
	return usr.code
}

func (usr *User) SetCode(fcode string) {
	usr.code = fcode
}

func (usr *User) GetLastName() string {
	return usr.lastName
}

func (usr *User) SetLastName(lname string) {
	usr.lastName = lname
}

func (usr *User) GetFirstName() string {
	return usr.firstName
}

func (usr *User) SetFirstName(fname string) {
	usr.firstName = fname
}

func (usr *User) GetMiddleName() string {
	return usr.middleName
}

func (usr *User) SetMiddleName(mname string) {
	usr.middleName = mname
}

func (usr *User) GetLogin() string {
	return usr.login
}

func (usr *User) SetLogin(flogin string) {
	usr.login = flogin
}

func (usr *User) GetEmail() string {
	return usr.email
}

func (usr *User) SetEmail(femail string) {
	usr.email = femail
}

func (usr *User) GetStatus() UserStatus {
	return usr.status
}

func (usr *User) SetStatus(fstatus UserStatus) {
	usr.status = fstatus
}
