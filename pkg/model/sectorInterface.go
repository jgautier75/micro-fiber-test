package model

import "database/sql"

type SectorStatus int64

const (
	SectorStatusDraft    SectorStatus = 0
	SectorStatusActive   SectorStatus = 1
	SectorStatusInactive SectorStatus = 2
	SectorStatusDeleted  SectorStatus = 3
)

type SectorInterface interface {
	GetId() int64
	SetId(id int64)
	GetTenantId() int64
	SetTenantId(id int64)
	GetOrgId() int64
	SetOrgId(orgId int64)
	GetCode() string
	SetCode(code string)
	GetLabel() string
	SetLabel(label string)
	GetParentId() sql.NullInt64
	SetParentId(orgId sql.NullInt64)
	GetHasParent() bool
	SetHasParent(hasParent bool)
	GetDepth() int
	SetDepth(depth int)
	GetSectorStatus() SectorStatus
	SetSectorStatus(status SectorStatus)
}
