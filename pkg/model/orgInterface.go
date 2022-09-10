package model

type OrganizationStatus int64

type OrganizationType string

const (
	OrgStatusDraft    OrganizationStatus = 0
	OrgStatusActive   OrganizationStatus = 1
	OrgStatusInactive OrganizationStatus = 2
	OrgStatusDeleted  OrganizationStatus = 3
)

const (
	OrgTypeLxsi       OrganizationType = "lxsi"
	OrgTypeBu         OrganizationType = "bu"
	OrgTypeCommunity  OrganizationType = "community"
	OrgTypeEnterprise OrganizationType = "enterprise"
)

type OrganizationInterface interface {
	GetId() int64
	SetId(id int64)
	GetTenantId() int64
	SetTenantId(id int64)
	GetCode() string
	SetCode(code string)
	GetLabel() string
	SetLabel(label string)
	GetType() OrganizationType
	SetType(aType OrganizationType)
	GetStatus() OrganizationStatus
	SetStatus(status OrganizationStatus)
}
