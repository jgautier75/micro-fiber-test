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
