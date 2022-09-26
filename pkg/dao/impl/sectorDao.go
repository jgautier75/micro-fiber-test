package impl

import (
	"context"
	"database/sql"
	"github.com/jackc/pgx/v4"
	"micro-fiber-test/pkg/dao/api"
	"micro-fiber-test/pkg/model"
)

type SectorDao struct {
}

func NewSectorDao() api.SectorDaoInterface {
	return &SectorDao{}
}

func (s SectorDao) CreateInTx(tx pgx.Tx, sector model.SectorInterface) (int64, error) {
	var id int64
	insertStmt := "insert into sectors(tenant_id,org_id,code,label,parent_id,has_parent,depth,status) values($1,$2,$3,$4,$5,$6,$7,$8) returning id"
	errQuery := tx.QueryRow(context.Background(), insertStmt, sector.GetTenantId(), sector.GetOrgId(), sector.GetCode(), sector.GetLabel(), sector.GetParentId(), sector.GetHasParent(), sector.GetDepth(), sector.GetSectorStatus()).Scan(&id)
	return id, errQuery
}

func (s SectorDao) Create(cnxParams string, sector model.SectorInterface) (int64, error) {
	conn, err := pgx.Connect(context.Background(), cnxParams)
	if err != nil {
		return -1, err
	}
	defer conn.Close(context.Background())
	if err != nil {
		return -1, err
	}
	var id int64
	insertStmt := "insert into sectors(tenant_id,org_id,code,label,parent_id,has_parent,depth,status) values($1,$2,$3,$4,$5,$6,$7,$8) returning id"
	errQuery := conn.QueryRow(context.Background(), insertStmt, sector.GetTenantId(), sector.GetOrgId(), sector.GetCode(), sector.GetLabel(), sector.GetParentId(), sector.GetHasParent(), sector.GetDepth(), sector.GetSectorStatus()).Scan(&id)
	return id, errQuery
}

func (s SectorDao) DeleteByOrgId(cnxParams string, orgId int64) error {
	conn, err := pgx.Connect(context.Background(), cnxParams)
	if err != nil {
		return err
	}
	defer conn.Close(context.Background())
	if err != nil {
		return err
	}
	deleteStmt := "delete from sectors where org_id=$1"
	_, errDelete := conn.Exec(context.Background(), deleteStmt, orgId)
	if errDelete != nil {
		return errDelete
	}
	return nil
}

func (s SectorDao) FindSectorsByTenantOrg(cnxParams string, tenantId int64, orgId int64) ([]model.SectorInterface, error) {
	conn, err := pgx.Connect(context.Background(), cnxParams)
	if err != nil {
		return nil, err
	}
	defer conn.Close(context.Background())
	if err != nil {
		return nil, err
	}
	selStmt := "select id,tenant_id,org_id,code,label,parent_id,has_parent,depth,status from sectors where tenant_id=$1 and org_id=$2"
	rows, e := conn.Query(context.Background(), selStmt, tenantId, orgId)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
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
		err = rows.Scan(&id, &tenantId, &orgId, &rsCode, &label, &parentId, &hasParent, &depth, &status)
		if e != nil {
			return nil, err
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

func (s SectorDao) FindByCode(cnxParams string, defaultTenantId int64, code string) (model.SectorInterface, error) {
	conn, err := pgx.Connect(context.Background(), cnxParams)
	if err != nil {
		return nil, err
	}
	defer conn.Close(context.Background())
	if err != nil {
		return nil, err
	}
	selStmt := "select id,tenant_id,org_id,code,label,parent_id,has_parent,depth,status from sectors where tenant_id=$1 and code=$2"
	rows, e := conn.Query(context.Background(), selStmt, defaultTenantId, code)
	if e != nil {
		return nil, e
	}
	defer rows.Close()
	if err != nil {
		return nil, err
	}
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
		err = rows.Scan(&id, &tenantId, &orgId, &rsCode, &label, &parentId, &hasParent, &depth, &status)
		if e != nil {
			return nil, err
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

func (s SectorDao) FindRootSector(cnxParams string, defaultTenantId int64, orgId int64) (int64, error) {
	conn, err := pgx.Connect(context.Background(), cnxParams)
	if err != nil {
		return 0, err
	}
	defer conn.Close(context.Background())
	if err != nil {
		return 0, err
	}
	selStmt := "select id from sectors where tenant_id=$1 and org_id=$2 and has_parent=$3"
	rows, e := conn.Query(context.Background(), selStmt, defaultTenantId, orgId, false)
	if e != nil {
		return 0, e
	}
	defer rows.Close()
	if err != nil {
		return 0, err
	}
	var sectorId int64
	for rows.Next() {
		var id int64
		err = rows.Scan(&id)
		if e != nil {
			return 0, err
		}
		sectorId = id
	}
	return sectorId, nil
}

func (s SectorDao) DeleteSector(cnxParams string, defaultTenantId int64, sectorId int64) error {
	conn, err := pgx.Connect(context.Background(), cnxParams)
	if err != nil {
		return err
	}
	defer conn.Close(context.Background())
	if err != nil {
		return err
	}
	deleteStmt := "delete from sectors where tenant_id=$1 and (id=$2 or parent_id=$3)"
	_, e := conn.Exec(context.Background(), deleteStmt, defaultTenantId, sectorId, sectorId)
	return e
}
