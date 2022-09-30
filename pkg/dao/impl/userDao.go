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
	insertStmt := "insert into users(tenant_id,org_id,external_id,last_name,first_name,middle_name,login,email,status) values($1,$2,$3,$4,$5,$6,$7,$8,$9) returning id"
	errQuery := conn.QueryRow(context.Background(), insertStmt, user.GetTenantId(), user.GetOrgId(), user.GetExternalId(), user.GetLastName(), user.GetFirstName(), user.GetMiddleName(), user.GetLogin(), user.GetEmail(), user.GetStatus()).Scan(&id)
	return id, errQuery
}

func (u UserDao) Update(cnxParams string, user model.UserInterface) error {
	conn, err := pgx.Connect(context.Background(), cnxParams)
	if err != nil {
		return err
	}
	defer conn.Close(context.Background())
	if err != nil {
		return err
	}
	updateStmt := "update users set last_name=$1,first_name=$2,middle_name=$3 where external_id=$4"
	_, errQuery := conn.Exec(context.Background(), updateStmt, user.GetLastName(), user.GetFirstName(), user.GetMiddleName(), user.GetExternalId())
	if errQuery != nil {
		return errQuery
	}
	return nil
}

func (u UserDao) CountByCriteria(cnxParams string, criteria model.UserFilterCriteria) (int, error) {
	conn, err := pgx.Connect(context.Background(), cnxParams)
	if err != nil {
		return 0, err
	}
	defer conn.Close(context.Background())
	if err != nil {
		return 0, err
	}
	var fullQry strings.Builder
	qryPrefix := "select count(1) from users where tenant_id=$1 and org_id=$2"
	whereClause, vals := computeFindByCriteriaQuery(qryPrefix, criteria)
	fullQry.WriteString(whereClause)
	countRes, errCount := conn.Query(context.Background(), fullQry.String(), vals...)
	if errCount != nil {
		return 0, errCount
	}
	cnt := 0
	for countRes.Next() {
		err := countRes.Scan(&cnt)
		if err != nil {
			return 0, err
		}
	}
	return cnt, nil
}

func (u UserDao) FindByCriteria(cnxParams string, criteria model.UserFilterCriteria) (model.UserSearchResult, error) {
	conn, err := pgx.Connect(context.Background(), cnxParams)
	searchResults := model.UserSearchResult{}
	if err != nil {
		return searchResults, err
	}
	defer conn.Close(context.Background())
	if err != nil {
		return searchResults, err
	}

	var fullQry strings.Builder
	qryPrefix := "select id,external_id,last_name,first_name,middle_name,login,email,status from users where tenant_id=$1 and org_id=$2"
	whereClause, vals := computeFindByCriteriaQuery(qryPrefix, criteria)
	fullQry.WriteString(whereClause)
	fullQry.WriteString(" order by  last_name,first_name asc")

	if criteria.Page > 1 {
		startPg := (criteria.Page - 1) * criteria.RowsPerPage
		fullQry.WriteString(" offset ")
		fullQry.WriteString(strconv.Itoa(startPg))
	}

	fullQry.WriteString(" limit ")
	fullQry.WriteString(strconv.Itoa(criteria.RowsPerPage))

	query := fullQry.String()

	rows, errQuery := conn.Query(context.Background(), query, vals...)
	if errQuery != nil {
		return searchResults, errQuery
	}
	defer rows.Close()
	var users []model.UserInterface
	for rows.Next() {
		var id int64
		var externalId string
		var lastName string
		var firstName string
		var middleName string
		var login string
		var email string
		var status int64
		err = rows.Scan(&id, &externalId, &lastName, &firstName, &middleName, &login, &email, &status)
		if err != nil {
			return searchResults, err
		}
		userInterface := model.User{}
		userInterface.SetId(id)
		userInterface.SetExternalId(externalId)
		userInterface.SetLastName(lastName)
		userInterface.SetFirstName(firstName)
		userInterface.SetMiddleName(middleName)
		userInterface.SetLogin(login)
		userInterface.SetEmail(email)
		userInterface.SetStatus(model.UserStatus(status))
		users = append(users, &userInterface)
	}

	searchResults.Users = users

	return searchResults, nil
}

func (u UserDao) FindByCode(cnxParams string, tenantId int64, orgId int64, externalId string) (model.UserInterface, error) {
	conn, err := pgx.Connect(context.Background(), cnxParams)
	if err != nil {
		return nil, err
	}
	defer conn.Close(context.Background())
	if err != nil {
		return nil, err
	}
	qry := "select id,external_id,last_name,first_name,middle_name,login,email,status from users where tenant_id=$1 and org_id=$2 and external_id=$3"
	rows, errQuery := conn.Query(context.Background(), qry, tenantId, orgId, externalId)
	if errQuery != nil {
		return nil, errQuery
	}
	defer rows.Close()
	userInterface := model.User{}
	for rows.Next() {
		var id int64
		var extId string
		var lastName string
		var firstName string
		var middleName string
		var login string
		var email string
		var status int64
		err = rows.Scan(&id, &extId, &lastName, &firstName, &middleName, &login, &email, &status)
		if err != nil {
			return nil, err
		}
		userInterface := model.User{}
		userInterface.SetId(id)
		userInterface.SetExternalId(extId)
		userInterface.SetLastName(lastName)
		userInterface.SetFirstName(firstName)
		userInterface.SetMiddleName(middleName)
		userInterface.SetLogin(login)
		userInterface.SetEmail(email)
		userInterface.SetStatus(model.UserStatus(status))
		return &userInterface, nil
	}
	return &userInterface, nil
}

func computeFindByCriteriaQuery(qryPrefix string, criteria model.UserFilterCriteria) (string string, params []interface{}) {

	var values []interface{}
	var buf strings.Builder

	values = append(values, criteria.TenantId)
	values = append(values, criteria.OrgId)

	buf.WriteString(qryPrefix)

	inc := 3
	if criteria.Login != "" {
		values = append(values, criteria.Login)
		buf.WriteString(fmt.Sprintf(" and login=%s", "%"+strconv.Itoa(inc)))
		inc = inc + 1
	}
	if criteria.Email != "" {
		values = append(values, criteria.Email)
		buf.WriteString(fmt.Sprintf(" and email=%s", "%"+strconv.Itoa(inc)))
		inc = inc + 1
	}
	if criteria.LastName != "" {
		values = append(values, criteria.LastName)
		buf.WriteString(fmt.Sprintf(" and last_name=%s", "%"+strconv.Itoa(inc)))
		inc = inc + 1
	}
	if criteria.FirstName != "" {
		values = append(values, criteria.FirstName)
		buf.WriteString(fmt.Sprintf(" and first_name=%s", "%"+strconv.Itoa(inc)))
	}
	fullQry := buf.String()
	return fullQry, values
}
