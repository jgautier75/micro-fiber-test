package impl

import (
	"context"
	"github.com/jackc/pgx/v4"
	"micro-fiber-test/pkg/dao/api"
	"micro-fiber-test/pkg/model"
)

type SectorDao struct {
}

func NewSectorDao() api.SectorDaoInterface {
	return &SectorDao{}
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
