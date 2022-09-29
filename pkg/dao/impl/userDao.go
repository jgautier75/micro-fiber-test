package impl

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4"
	"micro-fiber-test/pkg/dao/api"
	"micro-fiber-test/pkg/model"
	"strconv"
	"strings"
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

func (u UserDao) FindByCriteria(cnxParams string, criteria model.UserFilterCriteria) ([]model.UserInterface, error) {
	conn, err := pgx.Connect(context.Background(), cnxParams)
	if err != nil {
		return nil, err
	}
	defer conn.Close(context.Background())
	if err != nil {
		return nil, err
	}

	var values []interface{}
	var strBuf strings.Builder

	values = append(values, criteria.TenantId)
	values = append(values, criteria.OrgId)

	strBuf.WriteString("select id,code,last_name,first_name,middle_name,login,email,status from users where tenant_id=$1 and org_id=$2")
	inc := 3
	if criteria.Login != "" {
		values = append(values, criteria.Login)
		strBuf.WriteString(fmt.Sprintf(" and login=%s", "%"+strconv.Itoa(inc)))
		inc = inc + 1
	}
	if criteria.Email != "" {
		values = append(values, criteria.Email)
		strBuf.WriteString(fmt.Sprintf(" and email=%s", "%"+strconv.Itoa(inc)))
		inc = inc + 1
	}
	if criteria.LastName != "" {
		values = append(values, criteria.LastName)
		strBuf.WriteString(fmt.Sprintf(" and last_name=%s", "%"+strconv.Itoa(inc)))
		inc = inc + 1
	}
	if criteria.FirstName != "" {
		values = append(values, criteria.FirstName)
		strBuf.WriteString(fmt.Sprintf(" and first_name=%s", "%"+strconv.Itoa(inc)))
	}

	strBuf.WriteString(" order by last_name,first_name")

	rows, errQuery := conn.Query(context.Background(), strBuf.String(), values...)
	if errQuery != nil {
		return nil, errQuery
	}
	defer rows.Close()
	var users []model.UserInterface
	for rows.Next() {
		var id int64
		var code string
		var lastName string
		var firstName string
		var middleName string
		var login string
		var email string
		var status int64
		err = rows.Scan(&id, &code, &lastName, &firstName, &middleName, &login, &email, &status)
		if err != nil {
			return nil, err
		}
		userInterface := model.User{}
		userInterface.SetId(id)
		userInterface.SetCode(code)
		userInterface.SetLastName(lastName)
		userInterface.SetFirstName(firstName)
		userInterface.SetMiddleName(middleName)
		userInterface.SetLogin(login)
		userInterface.SetEmail(email)
		userInterface.SetStatus(model.UserStatus(status))
		users = append(users, &userInterface)
	}

	return users, nil
}
