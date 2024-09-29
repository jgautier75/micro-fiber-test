package impl

import (
	"context"
	"micro-fiber-test/pkg/model"
	"micro-fiber-test/pkg/repository/api"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/knadh/koanf"
)

type OrgDao struct {
	dbPool *pgxpool.Pool
	koanf  *koanf.Koanf
}

func NewOrgDao(pool *pgxpool.Pool, kSql *koanf.Koanf) api.OrgDaoInterface {
	orgDao := OrgDao{}
	orgDao.dbPool = pool
	orgDao.koanf = kSql
	return &orgDao
}

func (orgRepo *OrgDao) CreateInTx(tx pgx.Tx, org model.Organization) (int64, error) {
	var id int64
	insertStmt := orgRepo.koanf.String("organizations.create")
	errQuery := tx.QueryRow(context.Background(), insertStmt, org.TenantId, org.Code, org.Label, org.Type, org.Status).Scan(&id)
	return id, errQuery
}

func (orgRepo *OrgDao) Create(org model.Organization) (int64, error) {
	var id int64
	insertStmt := orgRepo.koanf.String("organizations.create")
	errQuery := orgRepo.dbPool.QueryRow(context.Background(), insertStmt, org.TenantId, org.Code, org.Label, org.Type, org.Status).Scan(&id)
	return id, errQuery
}

func (orgRepo *OrgDao) Update(orgCode string, label string) error {
	updateStmt := orgRepo.koanf.String("organizations.update")
	_, errQuery := orgRepo.dbPool.Exec(context.Background(), updateStmt, label, orgCode)
	return errQuery
}

func (orgRepo *OrgDao) Delete(orgCode string) error {
	deleteStmt := orgRepo.koanf.String("organizations.update")
	_, errQuery := orgRepo.dbPool.Exec(context.Background(), deleteStmt, orgCode)
	return errQuery
}

func (orgRepo *OrgDao) FindByCode(code string) (model.Organization, error) {
	var nilOrg model.Organization
	selStmt := orgRepo.koanf.String("organizations.findbycode")
	rows, e := orgRepo.dbPool.Query(context.Background(), selStmt, code)
	if e != nil {
		return nilOrg, e
	}
	defer rows.Close()
	org, errCollect := pgx.CollectOneRow(rows, pgx.RowToStructByName[model.Organization])
	if errCollect != nil {
		return nilOrg, errCollect
	}

	return org, nil
}

func (orgRepo *OrgDao) FindAll(tenantId int64) ([]model.Organization, error) {
	var nilOrg []model.Organization
	selStmt := orgRepo.koanf.String("organizations.findall")
	rows, errQry := orgRepo.dbPool.Query(context.Background(), selStmt, tenantId)
	if errQry != nil {
		return nil, errQry
	}
	defer rows.Close()
	orgs, errCollect := pgx.CollectRows(rows, pgx.RowToStructByName[model.Organization])
	if errCollect != nil {
		return nilOrg, errCollect
	}
	return orgs, nil
}

func (orgRepo *OrgDao) ExistsByCode(tenantId int64, code string) (bool, error) {
	selStmt := orgRepo.koanf.String("organizations.existsbycode")
	rows, e := orgRepo.dbPool.Query(context.Background(), selStmt, tenantId, code)
	if e != nil {
		return false, e
	}
	defer rows.Close()
	cnt := 0
	for rows.Next() {
		err := rows.Scan(&cnt)
		if err != nil {
			return false, err
		}
	}

	var exists = false
	if cnt > 0 {
		exists = true
	}
	return exists, nil
}

func (orgRepo *OrgDao) ExistsByLabel(tenantId int64, label string) (bool, error) {
	selStmt := orgRepo.koanf.String("organizations.findbylabel")
	rows, errQry := orgRepo.dbPool.Query(context.Background(), selStmt, tenantId, label)
	if errQry != nil {
		return false, errQry
	}
	defer rows.Close()
	cnt := 0
	for rows.Next() {
		err := rows.Scan(&cnt)
		if err != nil {
			return false, err
		}
	}
	var exists = false
	if cnt > 0 {
		exists = true
	}
	return exists, nil
}
