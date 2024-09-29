package api

import (
	"micro-fiber-test/pkg/model"
)

type UserServiceInterface interface {
	Create(defautTenantId int64, user model.User) (int64, error)
	Update(user model.User) error
	FindByCriteria(criteria model.UserFilterCriteria) (model.UserSearchResult, error)
	FindByCode(tenantId int64, orgId int64, externalId string) (model.User, error)
	Delete(externalId string) error
}
