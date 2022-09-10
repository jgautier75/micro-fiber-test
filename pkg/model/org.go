package model

type Organization struct {
	id       int64
	tenantId int64
	code     string
	label    string
	orgType  OrganizationType
	status   OrganizationStatus
}

func (org *Organization) GetId() int64 {
	return org.id
}

func (org *Organization) SetId(id int64) {
	org.id = id
}

func (org *Organization) GetTenantId() int64 {
	return org.tenantId
}

func (org *Organization) SetTenantId(tenantId int64) {
	org.tenantId = tenantId
}

func (org *Organization) GetCode() string {
	return org.code
}

func (org *Organization) SetCode(code string) {
	org.code = code
}

func (org *Organization) GetLabel() string {
	return org.label
}

func (org *Organization) SetLabel(label string) {
	org.label = label
}

func (org *Organization) GetType() OrganizationType {
	return org.orgType
}

func (org *Organization) SetType(orgType OrganizationType) {
	org.orgType = orgType
}

func (org *Organization) GetStatus() OrganizationStatus {
	return org.status
}

func (org *Organization) SetStatus(status OrganizationStatus) {
	org.status = status
}
