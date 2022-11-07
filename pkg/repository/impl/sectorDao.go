package impl

import (
	"context"
	"database/sql"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/knadh/koanf"
	"micro-fiber-test/pkg/model"
	"micro-fiber-test/pkg/repository/api"
)

type SectorDao struct {
	dbPool *pgxpool.Pool
	koanf  *koanf.Koanf
}

func NewSectorDao(pool *pgxpool.Pool, kSql *koanf.Koanf) api.SectorDaoInterface {
	sectorDao := SectorDao{}
	sectorDao.dbPool = pool
	sectorDao.koanf = kSql
	return &sectorDao
}

func (s SectorDao) CreateInTx(tx pgx.Tx, sector model.SectorInterface) (int64, error) {
	var id int64
	insertStmt := s.koanf.String("sectors.create")
	errQuery := tx.QueryRow(context.Background(), insertStmt, sector.GetTenantId(), sector.GetOrgId(), sector.GetCode(), sector.GetLabel(), sector.GetParentId(), sector.GetHasParent(), sector.GetDepth(), sector.GetSectorStatus()).Scan(&id)
	return id, errQuery
}

func (s SectorDao) Create(sector model.SectorInterface) (int64, error) {
	var id int64
	insertStmt := s.koanf.String("sectors.create")
	errQuery := s.dbPool.QueryRow(context.Background(), insertStmt, sector.GetTenantId(), sector.GetOrgId(), sector.GetCode(), sector.GetLabel(), sector.GetParentId(), sector.GetHasParent(), sector.GetDepth(), sector.GetSectorStatus()).Scan(&id)
	return id, errQuery
}

func (s SectorDao) DeleteByOrgId(orgId int64) error {
	deleteStmt := s.koanf.String("sectors.deletebyorgid")
	_, errDelete := s.dbPool.Exec(context.Background(), deleteStmt, orgId)
	return errDelete
}

func (s SectorDao) FindSectorsByTenantOrg(tenantId int64, orgId int64) ([]model.SectorInterface, error) {
	selStmt := s.koanf.String("sectors.findbytenantorg")
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
	selStmt := s.koanf.String("sectors.findbylabel")
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
	selStmt := s.koanf.String("sectors.findbycode")
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
	selStmt := s.koanf.String("sectors.findroot")
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
	deleteStmt := s.koanf.String("sectors.delete")
	_, e := s.dbPool.Exec(context.Background(), deleteStmt, defaultTenantId, sectorId, sectorId)
	return e
}

func (s SectorDao) Update(defaultTenantId int64, id int64, label string) error {
	updateStmt := s.koanf.String("sectors.update")
	_, errQuery := s.dbPool.Exec(context.Background(), updateStmt, label, id, defaultTenantId)
	return errQuery
}
