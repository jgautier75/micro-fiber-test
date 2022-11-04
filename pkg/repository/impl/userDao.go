package impl

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/knadh/koanf"
	"micro-fiber-test/pkg/model"
	"micro-fiber-test/pkg/repository/api"
	"strconv"
	"strings"
)

const (
	LogicalOperatorAnd = "and"
	LogicalOperatorOr  = "or"
	WhereExprEq        = "%s=%s"
	WhereExprNotEq     = "%s!=%s"
	WhereExprIn        = "%s in(%s)"
	WhereExprNotIn     = "%s not in(%s)"
	WhereExprLike      = "%s like %s"
	WhereExprNotLike   = "%s not like %s"
)

type UserDao struct {
	dbPool *pgxpool.Pool
	koanf  *koanf.Koanf
}

func NewUserDao(pool *pgxpool.Pool, kSql *koanf.Koanf) api.UserDaoInterface {
	userDao := UserDao{}
	userDao.dbPool = pool
	userDao.koanf = kSql
	return &userDao
}

func (u UserDao) Create(user model.UserInterface) (int64, error) {
	var id int64
	insertStmt := u.koanf.String("users.create")
	errQuery := u.dbPool.QueryRow(context.Background(), insertStmt, user.GetTenantId(), user.GetOrgId(), user.GetExternalId(), user.GetLastName(), user.GetFirstName(), user.GetMiddleName(), user.GetLogin(), user.GetEmail(), user.GetStatus()).Scan(&id)
	return id, errQuery
}

func (u UserDao) Update(user model.UserInterface) error {
	updateStmt := "update users set last_name=$1,first_name=$2,middle_name=$3,login=$4,email=$5 where external_id=$6"
	_, errQuery := u.dbPool.Exec(context.Background(), updateStmt, user.GetLastName(), user.GetFirstName(), user.GetMiddleName(), user.GetLogin(), user.GetEmail(), user.GetExternalId())
	return errQuery
}

func (u UserDao) CountByCriteria(criteria model.UserFilterCriteria) (int, error) {
	var fullQry strings.Builder
	qryPrefix := "select count(1) from users"
	whereClause, vals := computeFindByCriteriaQuery(qryPrefix, criteria)
	fullQry.WriteString(whereClause)
	countRes, errCount := u.dbPool.Query(context.Background(), fullQry.String(), vals...)
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

func (u UserDao) FindByCriteria(criteria model.UserFilterCriteria) (model.UserSearchResult, error) {
	searchResults := model.UserSearchResult{}
	var fullQry strings.Builder
	qryPrefix := "select id,external_id,last_name,first_name,middle_name,login,email,status from users"
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

	rows, errQuery := u.dbPool.Query(context.Background(), query, vals...)
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
		errScan := rows.Scan(&id, &externalId, &lastName, &firstName, &middleName, &login, &email, &status)
		if errScan != nil {
			return searchResults, errScan
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

func (u UserDao) FindByExternalId(tenantId int64, orgId int64, externalId string) (model.UserInterface, error) {
	qry := "select id,external_id,last_name,first_name,middle_name,login,email,status from users where tenant_id=$1 and org_id=$2 and external_id=$3"
	rows, errQuery := u.dbPool.Query(context.Background(), qry, tenantId, orgId, externalId)
	if errQuery != nil {
		return nil, errQuery
	}
	defer rows.Close()
	///var userInterface model.User
	for rows.Next() {
		var id int64
		var extId string
		var lastName string
		var firstName string
		var middleName string
		var login string
		var email string
		var status int64
		errScan := rows.Scan(&id, &extId, &lastName, &firstName, &middleName, &login, &email, &status)
		if errScan != nil {
			return nil, errScan
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
	return nil, nil
}

func (u UserDao) IsLoginInUse(login string) (int64, string, error) {
	selStmt := "select id,external_id from users where login=$1"
	rows, errQuery := u.dbPool.Query(context.Background(), selStmt, login)
	defer rows.Close()
	if errQuery != nil {
		return 0, "", errQuery
	}
	for rows.Next() {
		var id int64
		var extId string
		errScan := rows.Scan(&id, &extId)
		if errScan != nil {
			return 0, "", errScan
		}
		return id, extId, nil
	}
	return 0, "", nil
}

func (u UserDao) IsEmailInUse(email string) (int64, string, error) {
	selStmt := "select id,external_id from users where email=$1"
	rows, errQuery := u.dbPool.Query(context.Background(), selStmt, email)
	defer rows.Close()
	if errQuery != nil {
		return 0, "", errQuery
	}
	for rows.Next() {
		var id int64
		var extId string
		errScan := rows.Scan(&id, &extId)
		if errScan != nil {
			return 0, "", errScan
		}
		return id, extId, nil
	}
	return 0, "", nil
}

func (u UserDao) Delete(userExtId string) error {
	selStmt := "delete from users where external_id=$1"
	rows, errQuery := u.dbPool.Query(context.Background(), selStmt, userExtId)
	defer rows.Close()
	return errQuery
}

func computeFindByCriteriaQuery(qryPrefix string, criteria model.UserFilterCriteria) (string string, params []interface{}) {

	var values []interface{}
	var buf strings.Builder
	inc := 1

	buf.WriteString(qryPrefix)
	buf.WriteString(" where ")

	values = append(values, criteria.TenantId)
	inc, whereTenant := addCriteria(WhereExprEq, "tenant_id", inc, LogicalOperatorAnd)
	buf.WriteString(whereTenant)

	values = append(values, criteria.OrgId)
	inc, whereOrg := addCriteria(WhereExprEq, "org_id", inc, LogicalOperatorAnd)
	buf.WriteString(whereOrg)

	if criteria.Login != "" {
		values = append(values, "%"+criteria.Login+"%")
		nextInc, whereLogin := addCriteria(WhereExprLike, "login", inc, LogicalOperatorAnd)
		inc = nextInc
		buf.WriteString(whereLogin)
	}
	if criteria.Email != "" {
		values = append(values, "%"+criteria.Email+"%")
		nextInc, whereEmail := addCriteria(WhereExprLike, "email", inc, LogicalOperatorAnd)
		inc = nextInc
		buf.WriteString(whereEmail)
	}
	if criteria.LastName != "" {
		values = append(values, "%"+criteria.LastName+"%")
		nextInc, whereLastName := addCriteria(WhereExprLike, "last_name", inc, LogicalOperatorAnd)
		inc = nextInc
		inc = inc + 1
		buf.WriteString(whereLastName)
	}
	if criteria.FirstName != "" {
		values = append(values, "%"+criteria.FirstName+"%")
		nextInc, whereFirstName := addCriteria(WhereExprLike, "first_name", inc, LogicalOperatorAnd)
		inc = nextInc
		inc = inc + 1
		buf.WriteString(whereFirstName)
	}
	fullQry := buf.String()
	return fullQry, values
}

func addCriteria(expression string, criteriaName string, inc int, logicalOperator string) (nextInc int, qry string) {
	var whereClause string
	var next int
	if inc == 0 {
		inc = 1
		next = 1
	} else {
		next = inc + 1
	}
	if inc > 1 {
		whereClause = " " + logicalOperator + " "
	} else {
		whereClause = ""
	}
	whereClause = whereClause + fmt.Sprintf(expression, criteriaName, "$"+strconv.Itoa(inc))
	return next, whereClause
}
