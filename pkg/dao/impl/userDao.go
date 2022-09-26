package impl

import (
	"context"
	"github.com/jackc/pgx/v4"
	"micro-fiber-test/pkg/dao/api"
	"micro-fiber-test/pkg/model"
)

type UserDao struct {
}

func NewUserDao() api.UserDaoInterface {
	return &UserDao{}
}

func (u UserDao) Create(cnxParams string, user model.UserInterface) (int64, error) {
	conn, err := pgx.Connect(context.Background(), cnxParams)
	if err != nil {
		return -1, err
	}
	defer conn.Close(context.Background())
	if err != nil {
		return -1, err
	}
	var id int64
	insertStmt := "insert into users(tenant_id,org_id,code,last_name,first_name,middle_name,login,email,status) values($1,$2,$3,$4,$5,$6,$7,$8,$9) returning id"
	errQuery := conn.QueryRow(context.Background(), insertStmt, user.GetTenantId(), user.GetOrgId(), user.GetCode(), user.GetLastName(), user.GetFirstName(), user.GetMiddleName(), user.GetLogin(), user.GetEmail(), user.GetStatus()).Scan(&id)
	return id, errQuery
}
