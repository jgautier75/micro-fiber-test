package impl

import (
	"context"
	"fmt"
	"micro-fiber-test/pkg/model"
	"micro-fiber-test/pkg/repository/api"
	"strconv"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/knadh/koanf"
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

func (u UserDao) Create(user model.User) (int64, error) {
	var id int64
	insertStmt := u.koanf.String("users.create")
	errQuery := u.dbPool.QueryRow(context.Background(), insertStmt, user.TenantId, user.OrgId, user.ExternalId, user.LastName, user.FirstName, user.MiddleName, user.Login, user.Email, user.Status).Scan(&id)
	return id, errQuery
}

func (u UserDao) Update(user model.User) error {
	updateStmt := u.koanf.String("users.update_by_external_id")
	_, errQuery := u.dbPool.Exec(context.Background(), updateStmt, user.LastName, user.FirstName, user.MiddleName, user.Login, user.Email, user.ExternalId)
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
	qryPrefix := u.koanf.String("users.find_by_query")
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

	users, errCollect := pgx.CollectRows(rows, pgx.RowToStructByName[model.User])
	if errCollect != nil {
		return searchResults, errCollect
	}

	searchResults.Users = users

	return searchResults, nil
}

func (u UserDao) FindByExternalId(tenantId int64, orgId int64, externalId string) (model.User, error) {
	qry := u.koanf.String("users.find_by_external_id")
	var nilUser model.User
	rows, errQuery := u.dbPool.Query(context.Background(), qry, tenantId, orgId, externalId)
	if errQuery != nil {
		return nilUser, errQuery
	}
	defer rows.Close()

	userInterface, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[model.User])
	if err != nil {
		return nilUser, err
	}
	return userInterface, nil
}

func (u UserDao) IsLoginInUse(login string) (int64, string, error) {
	selStmt := u.koanf.String("users.find_by_login")
	rows, errQuery := u.dbPool.Query(context.Background(), selStmt, login)
	if errQuery != nil {
		return 0, "", errQuery
	}
	defer rows.Close()

	var id int64
	var extId string
	if rows.Next() {
		errScan := rows.Scan(&id, &extId)
		if errScan != nil {
			return 0, "", errScan
		}
		return id, extId, nil
	}
	return 0, "", nil
}

func (u UserDao) IsEmailInUse(email string) (int64, string, error) {
	selStmt := u.koanf.String("users.email_in_user")
	rows, errQuery := u.dbPool.Query(context.Background(), selStmt, email)
	if errQuery != nil {
		return 0, "", errQuery
	}
	defer rows.Close()

	var id int64
	var extId string
	if rows.Next() {
		errScan := rows.Scan(&id, &extId)
		if errScan != nil {
			return 0, "", errScan
		}
		return id, extId, nil
	}
	return 0, "", nil
}

func (u UserDao) Delete(userExtId string) error {
	selStmt := u.koanf.String("users.delete_by_external_id")
	rows, errQuery := u.dbPool.Query(context.Background(), selStmt, userExtId)
	if errQuery != nil {
		return errQuery
	}
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
		buf.WriteString(whereLastName)
	}
	if criteria.FirstName != "" {
		values = append(values, "%"+criteria.FirstName+"%")
		nextInc, whereFirstName := addCriteria(WhereExprLike, "first_name", inc, LogicalOperatorAnd)
		inc = nextInc
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
