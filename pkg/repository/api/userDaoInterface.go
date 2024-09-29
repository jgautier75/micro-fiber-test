package api

import (
	"micro-fiber-test/pkg/model"
)

type UserDaoInterface interface {
	Create(user model.User) (int64, error)
	FindByExternalId(tenantId int64, orgId int64, externalId string) (model.User, error)
	FindByCriteria(criteria model.UserFilterCriteria) (model.UserSearchResult, error)
	CountByCriteria(criteria model.UserFilterCriteria) (int, error)
	Update(user model.User) error
	IsLoginInUse(login string) (int64, string, error)
	IsEmailInUse(email string) (int64, string, error)
	Delete(userExtId string) error
}
