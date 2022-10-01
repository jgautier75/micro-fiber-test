package impl

import (
	"context"
	"database/sql"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"micro-fiber-test/pkg/dao/api"
	"micro-fiber-test/pkg/model"
)

type SectorDao struct {
	dbPool *pgxpool.Pool
}

func NewSectorDao(pool *pgxpool.Pool) api.SectorDaoInterface {
	sectorDao := SectorDao{}
	sectorDao.dbPool = pool
	return &sectorDao
}

func (s SectorDao) CreateInTx(tx pgx.Tx, sector model.SectorInterface) (int64, error) {
	var id int64
	insertStmt := "insert into sectors(tenant_id,org_id,code,label,parent_id,has_parent,depth,status) values($1,$2,$3,$4,$5,$6,$7,$8) returning id"
	errQuery := tx.QueryRow(context.Background(), insertStmt, sector.GetTenantId(), sector.GetOrgId(), sector.GetCode(), sector.GetLabel(), sector.GetParentId(), sector.GetHasParent(), sector.GetDepth(), sector.GetSectorStatus()).Scan(&id)
	return id, errQuery
}

func (s SectorDao) Create(sector model.SectorInterface) (int64, error) {
	var id int64
	insertStmt := "insert into sectors(tenant_id,org_id,code,label,parent_id,has_parent,depth,status) values($1,$2,$3,$4,$5,$6,$7,$8) returning id"
	errQuery := s.dbPool.QueryRow(context.Background(), insertStmt, sector.GetTenantId(), sector.GetOrgId(), sector.GetCode(), sector.GetLabel(), sector.GetParentId(), sector.GetHasParent(), sector.GetDepth(), sector.GetSectorStatus()).Scan(&id)
	return id, errQuery
}

func (s SectorDao) DeleteByOrgId(orgId int64) error {
	deleteStmt := "delete from sectors where org_id=$1"
	_, errDelete := s.dbPool.Exec(context.Background(), deleteStmt, orgId)
	if errDelete != nil {
		return errDelete
	}
	return nil
}

func (s SectorDao) FindSectorsByTenantOrg(tenantId int64, orgId int64) ([]model.SectorInterface, error) {
	selStmt := "select id,tenant_id,org_id,code,label,parent_id,has_parent,depth,status from sectors where tenant_id=$1 and org_id=$2"
	rows, e := s.dbPool.Query(context.Background(), selStmt, tenantId, orgId)
	defer rows.Close()
	if e != nil {
		return nil, e
	}

	var sectors []model.SectorInterface
	for rows.Next() {
		var id int64
		var tenantId int64
		var orgId int64
		var rsCode string
		var label string
		var parentId sql.NullInt64
		var hasParent bool
		var depth int
		var status model.SectorStatus
		errScan := rows.Scan(&id, &tenantId, &orgId, &rsCode, &label, &parentId, &hasParent, &depth, &status)
		if errScan != nil {
			return nil, errScan
		}
		sector := model.Sector{}
		sector.SetId(id)
		sector.SetTenantId(tenantId)
		sector.SetOrgId(orgId)
		sector.SetCode(rsCode)
		sector.SetLabel(label)
		sector.SetParentId(parentId)
		sector.SetHasParent(hasParent)
		sector.SetDepth(depth)
		sector.SetSectorStatus(status)
		sectors = append(sectors, &sector)
	}
	return sectors, nil
}

func (s SectorDao) FindByLabel(defaultTenantId int64, label string) (int64, string, error) {
	selStmt := "select id,code from sectors where tenant_id=$1 and label=$2"
	rows, errQry := s.dbPool.Query(context.Background(), selStmt, defaultTenantId, label)
	defer rows.Close()
	if errQry != nil {
		return 0, "", errQry
	}
	for rows.Next() {
		var rsCode string
		var rsId int64
		errScan := rows.Scan(&rsId, &rsCode)
		if errScan != nil {
			return 0, "", errScan
		}
		return rsId, rsCode, nil
	}
	return 0, "", nil
}

func (s SectorDao) FindByCode(defaultTenantId int64, code string) (model.SectorInterface, error) {
	selStmt := "select id,tenant_id,org_id,code,label,parent_id,has_parent,depth,status from sectors where tenant_id=$1 and code=$2"
	rows, errQry := s.dbPool.Query(context.Background(), selStmt, defaultTenantId, code)
	if errQry != nil {
		return nil, errQry
	}
	defer rows.Close()
	for rows.Next() {
		var id int64
		var tenantId int64
		var orgId int64
		var rsCode string
		var label string
		var parentId sql.NullInt64
		var hasParent bool
		var depth int
		var status model.SectorStatus
		errScan := rows.Scan(&id, &tenantId, &orgId, &rsCode, &label, &parentId, &hasParent, &depth, &status)
		if errScan != nil {
			return nil, errScan
		}
		sector := model.Sector{}
		sector.SetId(id)
		sector.SetTenantId(tenantId)
		sector.SetOrgId(orgId)
		sector.SetCode(rsCode)
		sector.SetLabel(label)
		sector.SetParentId(parentId)
		sector.SetHasParent(hasParent)
		sector.SetDepth(depth)
		sector.SetSectorStatus(status)
		return &sector, nil
	}
	return nil, nil
}

func (s SectorDao) FindRootSector(defaultTenantId int64, orgId int64) (int64, error) {
	selStmt := "select id from sectors where tenant_id=$1 and org_id=$2 and has_parent=$3"
	rows, errQry := s.dbPool.Query(context.Background(), selStmt, defaultTenantId, orgId, false)
	if errQry != nil {
		return 0, errQry
	}
	defer rows.Close()
	var sectorId int64
	for rows.Next() {
		var id int64
		errScan := rows.Scan(&id)
		if errScan != nil {
			return 0, errScan
		}
		sectorId = id
	}
	return sectorId, nil
}

func (s SectorDao) DeleteSector(defaultTenantId int64, sectorId int64) error {
	deleteStmt := "delete from sectors where tenant_id=$1 and (id=$2 or parent_id=$3)"
	_, e := s.dbPool.Exec(context.Background(), deleteStmt, defaultTenantId, sectorId, sectorId)
	return e
}

func (s SectorDao) Update(defaultTenantId int64, id int64, label string) error {
	updateStmt := "update sectors set label=$1 where id=$2 and tenant_id=$3"
	_, errQuery := s.dbPool.Exec(context.Background(), updateStmt, label, id, defaultTenantId)
	if errQuery != nil {
		return errQuery
	}
	return nil
}
