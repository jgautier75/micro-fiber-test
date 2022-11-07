package impl

import (
	"context"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/knadh/koanf"
	"micro-fiber-test/pkg/model"
	"micro-fiber-test/pkg/repository/api"
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

func (orgRepo *OrgDao) CreateInTx(tx pgx.Tx, org model.OrganizationInterface) (int64, error) {
	var id int64
	insertStmt := orgRepo.koanf.String("organizations.create")
	errQuery := tx.QueryRow(context.Background(), insertStmt, org.GetTenantId(), org.GetCode(), org.GetLabel(), org.GetType(), org.GetStatus()).Scan(&id)
	return id, errQuery
}

func (orgRepo *OrgDao) Create(org model.OrganizationInterface) (int64, error) {
	var id int64
	insertStmt := orgRepo.koanf.String("organizations.create")
	errQuery := orgRepo.dbPool.QueryRow(context.Background(), insertStmt, org.GetTenantId(), org.GetCode(), org.GetLabel(), org.GetType(), org.GetStatus()).Scan(&id)
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

func (orgRepo *OrgDao) FindByCode(code string) (model.OrganizationInterface, error) {
	selStmt := orgRepo.koanf.String("organizations.findbycode")
	rows, e := orgRepo.dbPool.Query(context.Background(), selStmt, code)
	if e != nil {
		return nil, e
	}
	defer rows.Close()

	for rows.Next() {
		var id int64
		var tenantId int64
		var rsType model.OrganizationType
		var rsCode string
		var label string
		var status model.OrganizationStatus
		errScan := rows.Scan(&id, &tenantId, &rsCode, &label, &rsType, &status)
		if e != nil {
			return nil, errScan
		}
		org := model.Organization{}
		org.SetId(id)
		org.SetTenantId(tenantId)
		org.SetCode(rsCode)
		org.SetLabel(label)
		org.SetType(rsType)
		org.SetStatus(status)
		return &org, nil
	}
	return nil, nil
}

func (orgRepo *OrgDao) FindAll(tenantId int64) ([]model.OrganizationInterface, error) {
	selStmt := orgRepo.koanf.String("organizations.findall")
	rows, errQry := orgRepo.dbPool.Query(context.Background(), selStmt, tenantId)
	if errQry != nil {
		return nil, errQry
	}

	defer rows.Close()
	var orgs []model.OrganizationInterface
	for rows.Next() {
		var id int64
		var tenantId int64
		var rsCode string
		var label string
		var kind string
		var status int64
		errScan := rows.Scan(&id, &tenantId, &rsCode, &label, &kind, &status)
		if errScan != nil {
			return nil, errScan
		}
		org := model.Organization{}
		org.SetId(id)
		org.SetTenantId(tenantId)
		org.SetCode(rsCode)
		org.SetLabel(label)
		org.SetStatus(model.OrganizationStatus(status))
		org.SetType(model.OrganizationType(kind))
		orgs = append(orgs, &org)
	}
	return orgs, nil
}

func (orgRepo *OrgDao) ExistsByCode(tenantId int64, code string) (bool, error) {
	selStmt := orgRepo.koanf.String("organizations.existsbycode")
	rows, e := orgRepo.dbPool.Query(context.Background(), selStmt, tenantId, code)
	defer rows.Close()
	if e != nil {
		return false, e
	}
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
	defer rows.Close()
	if errQry != nil {
		return false, errQry
	}
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
