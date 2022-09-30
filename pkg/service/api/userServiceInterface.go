package api

import (
	"micro-fiber-test/pkg/model"
)

type UserServiceInterface interface {
	Create(cnxParams string, defautTenantId int64, user model.UserInterface) (int64, error)
	FindByCriteria(cnxParams string, criteria model.UserFilterCriteria) (model.UserSearchResult, error)
	FindByCode(cnxParams string, tenantId int64, orgId int64, externalId string) (model.UserInterface, error)
}
