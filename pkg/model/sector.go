package model

import "database/sql"

type Sector struct {
	Id        int64         `db:"id"`
	TenantId  int64         `db:"tenant_id"`
	OrgId     int64         `db:"org_id"`
	Code      string        `db:"code"`
	Label     string        `db:"label"`
	ParentId  sql.NullInt64 `db:"parent_id"`
	HasParent bool          `db:"has_parent"`
	Depth     int           `db:"depth"`
	Status    SectorStatus  `db:"status"`
}
