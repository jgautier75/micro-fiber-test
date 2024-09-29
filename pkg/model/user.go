package model

type User struct {
	Id         int64      `db:"id"`
	TenantId   int64      `db:"tenant_id"`
	OrgId      int64      `db:"org_id"`
	ExternalId string     `db:"external_id"`
	LastName   string     `db:"last_name"`
	FirstName  string     `db:"first_name"`
	MiddleName string     `db:"middle_name"`
	Login      string     `db:"login"`
	Email      string     `db:"email"`
	Status     UserStatus `db:"status"`
}
