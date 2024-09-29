package impl

import (
	"context"
	"micro-fiber-test/pkg/model"
	"micro-fiber-test/pkg/repository/api"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/knadh/koanf"
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

func (s SectorDao) CreateInTx(tx pgx.Tx, sector model.Sector) (int64, error) {
	var id int64
	insertStmt := s.koanf.String("sectors.create")
	errQuery := tx.QueryRow(context.Background(), insertStmt, sector.TenantId, sector.OrgId, sector.Code, sector.Label, sector.ParentId, sector.HasParent, sector.Depth, sector.Status).Scan(&id)
	return id, errQuery
}

func (s SectorDao) Create(sector model.Sector) (int64, error) {
	var id int64
	insertStmt := s.koanf.String("sectors.create")
	errQuery := s.dbPool.QueryRow(context.Background(), insertStmt, sector.TenantId, sector.OrgId, sector.Code, sector.Label, sector.ParentId, sector.HasParent, sector.Depth, sector.Status).Scan(&id)
	return id, errQuery
}

func (s SectorDao) DeleteByOrgId(orgId int64) error {
	deleteStmt := s.koanf.String("sectors.deletebyorgid")
	_, errDelete := s.dbPool.Exec(context.Background(), deleteStmt, orgId)
	return errDelete
}

func (s SectorDao) FindSectorsByTenantOrg(tenantId int64, orgId int64) ([]model.Sector, error) {
	selStmt := s.koanf.String("sectors.findbytenantorg")
	rows, e := s.dbPool.Query(context.Background(), selStmt, tenantId, orgId)
	if e != nil {
		return nil, e
	}
	defer rows.Close()

	sectors, errCollect := pgx.CollectRows(rows, pgx.RowToStructByName[model.Sector])
	if errCollect != nil {
		return nil, errCollect
	}
	return sectors, nil
}

func (s SectorDao) FindByLabel(defaultTenantId int64, label string) (int64, string, error) {
	selStmt := s.koanf.String("sectors.findbylabel")
	rows, errQry := s.dbPool.Query(context.Background(), selStmt, defaultTenantId, label)
	if errQry != nil {
		return 0, "", errQry
	}
	defer rows.Close()

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

func (s SectorDao) FindByCode(defaultTenantId int64, code string) (model.Sector, error) {
	var nilSector model.Sector
	selStmt := s.koanf.String("sectors.findbycode")
	rows, errQry := s.dbPool.Query(context.Background(), selStmt, defaultTenantId, code)
	if errQry != nil {
		return nilSector, errQry
	}
	defer rows.Close()

	sector, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[model.Sector])
	if err != nil {
		return nilSector, err
	}

	return sector, nil
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
	if rows.Err() != nil {
		return 0, rows.Err()
	}
	return sectorId, nil
}

func (s SectorDao) DeleteSector(defaultTenantId int64, sectorId int64) error {
	deleteStmt := s.koanf.String("sectors.delete")
	_, e := s.dbPool.Exec(context.Background(), deleteStmt, defaultTenantId, sectorId, sectorId)
	if e != nil {
		return e
	}
	return e
}

func (s SectorDao) Update(defaultTenantId int64, id int64, label string) error {
	updateStmt := s.koanf.String("sectors.update")
	_, errQuery := s.dbPool.Exec(context.Background(), updateStmt, label, id, defaultTenantId)
	return errQuery
}
