package model

import "database/sql"

type Sector struct {
	id        int64
	tenantId  int64
	orgId     int64
	code      string
	label     string
	parentId  sql.NullInt64
	hasParent bool
	depth     int
	status    SectorStatus
}

func (sector *Sector) GetId() int64 {
	return sector.id
}

func (sector *Sector) SetId(id int64) {
	sector.id = id
}

func (sector *Sector) GetTenantId() int64 {
	return sector.tenantId
}

func (sector *Sector) SetTenantId(tenantId int64) {
	sector.tenantId = tenantId
}

func (sector *Sector) GetOrgId() int64 {
	return sector.orgId
}

func (sector *Sector) SetOrgId(orgId int64) {
	sector.orgId = orgId
}

func (sector *Sector) GetCode() string {
	return sector.code
}

func (sector *Sector) SetCode(code string) {
	sector.code = code
}

func (sector *Sector) GetLabel() string {
	return sector.label
}

func (sector *Sector) SetLabel(label string) {
	sector.label = label
}

func (sector *Sector) GetParentId() sql.NullInt64 {
	return sector.parentId
}

func (sector *Sector) SetParentId(parentId sql.NullInt64) {
	sector.parentId = parentId
}

func (sector *Sector) GetHasParent() bool {
	return sector.hasParent
}

func (sector *Sector) SetHasParent(hasParent bool) {
	sector.hasParent = hasParent
}

func (sector *Sector) GetDepth() int {
	return sector.depth
}

func (sector *Sector) SetDepth(depth int) {
	sector.depth = depth
}

func (sector *Sector) GetSectorStatus() SectorStatus {
	return sector.status
}

func (sector *Sector) SetSectorStatus(status SectorStatus) {
	sector.status = status
}
