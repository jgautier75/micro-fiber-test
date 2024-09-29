package model

type Organization struct {
	Id       int64              `db:"id"`
	TenantId int64              `db:"tenant_id"`
	Code     string             `db:"code"`
	Label    string             `db:"label"`
	Type     OrganizationType   `db:"type"`
	Status   OrganizationStatus `db:"status"`
}
