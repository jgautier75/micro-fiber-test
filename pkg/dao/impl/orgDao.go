package impl

import (
	"context"
	pgx2 "github.com/jackc/pgx"
	"github.com/jackc/pgx/v4"
	"micro-fiber-test/pkg/dao/api"
	"micro-fiber-test/pkg/model"
)

type OrgDao struct {
}

func NewOrgDao() api.OrgDaoInterface {
	return &OrgDao{}
}

func (orgRepo *OrgDao) CreateInTx(tx pgx.Tx, org model.OrganizationInterface) (int64, error) {
	var id int64
	insertStmt := "insert into organizations(tenant_id,code,label,type,status) values($1,$2,$3,$4,$5) returning id"
	errQuery := tx.QueryRow(context.Background(), insertStmt, org.GetTenantId(), org.GetCode(), org.GetLabel(), org.GetType(), org.GetStatus()).Scan(&id)
	return id, errQuery
}

func (orgRepo *OrgDao) Create(cnxParams string, org model.OrganizationInterface) (int64, error) {
	conn, err := pgx.Connect(context.Background(), cnxParams)
	if err != nil {
		return -1, err
	}
	defer func(conn *pgx2.Conn, ctx context.Context) {
		err := conn.Close(ctx)
		if err != nil {

		}
	}(conn, context.Background())
	if err != nil {
		return -1, err
	}
	var id int64
	insertStmt := "insert into organizations(tenant_id,code,label,type,status) values($1,$2,$3,$4,$5) returning id"
	errQuery := conn.QueryRow(context.Background(), insertStmt, org.GetTenantId(), org.GetCode(), org.GetLabel(), org.GetType(), org.GetStatus()).Scan(&id)
	return id, errQuery
}

func (orgRepo *OrgDao) Update(cnxParams string, orgCode string, label string) error {
	conn, err := pgx.Connect(context.Background(), cnxParams)
	if err != nil {
		return err
	}
	defer func(conn *pgx2.Conn, ctx context.Context) {
		err := conn.Close(ctx)
		if err != nil {

		}
	}(conn, context.Background())
	if err != nil {
		return err
	}
	updateStmt := "update organizations set label=$1 where code=$2"
	_, errQuery := conn.Exec(context.Background(), updateStmt, label, orgCode)
	return errQuery
}

func (orgRepo *OrgDao) Delete(cnxParams string, orgCode string) error {
	conn, err := pgx.Connect(context.Background(), cnxParams)
	if err != nil {
		return err
	}
	defer func(conn *pgx2.Conn, ctx context.Context) {
		err := conn.Close(ctx)
		if err != nil {

		}
	}(conn, context.Background())
	deleteStmt := "delete from organizations where code=$1"
	_, errQuery := conn.Exec(context.Background(), deleteStmt, orgCode)
	return errQuery
}

func (orgRepo *OrgDao) FindByCode(cnxParams string, code string) (model.OrganizationInterface, error) {
	conn, err := pgx.Connect(context.Background(), cnxParams)
	if err != nil {
		return nil, err
	}
	defer func(conn *pgx2.Conn, ctx context.Context) {
		err := conn.Close(ctx)
		if err != nil {

		}
	}(conn, context.Background())
	selStmt := "select id,tenant_id,code,label,type,status from organizations where code=$1"
	rows, e := conn.Query(context.Background(), selStmt, code)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var id int64
		var tenantId int64
		var rsType model.OrganizationType
		var rsCode string
		var label string
		var status model.OrganizationStatus
		err = rows.Scan(&id, &tenantId, &rsCode, &label, &rsType, &status)
		if e != nil {
			return nil, err
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

func (orgRepo *OrgDao) FindAll(cnxParams string, tenantId int64) ([]model.OrganizationInterface, error) {
	conn, err := pgx.Connect(context.Background(), cnxParams)
	if err != nil {
		return nil, err
	}
	defer func(conn *pgx2.Conn, ctx context.Context) {
		err := conn.Close(ctx)
		if err != nil {

		}
	}(conn, context.Background())
	selStmt := "select id,tenant_id,code,label,type,status from organizations where tenant_id=$1"
	rows, e := conn.Query(context.Background(), selStmt, tenantId)
	if err != nil {
		return nil, err
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
		err = rows.Scan(&id, &tenantId, &rsCode, &label, &kind, &status)
		if e != nil {
			return nil, err
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

func (orgRepo *OrgDao) ExistsByCode(cnxParams string, tenantId int64, code string) (bool, error) {
	conn, err := pgx.Connect(context.Background(), cnxParams)
	if err != nil {
		return false, err
	}
	defer func(conn *pgx2.Conn, ctx context.Context) {
		err := conn.Close(ctx)
		if err != nil {

		}
	}(conn, context.Background())
	selStmt := "select count(1) from organizations where tenant_id=$1 and code=$2"
	rows, e := conn.Query(context.Background(), selStmt, tenantId, code)
	defer rows.Close()
	if e != nil {
		return false, err
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
