package api

import (
	"micro-fiber-test/pkg/model"
)

type UserDaoInterface interface {
	Create(cnxParams string, user model.UserInterface) (int64, error)
	FindByExternalId(cnxParams string, tenantId int64, orgId int64, externalId string) (model.UserInterface, error)
	FindByCriteria(cnxParams string, criteria model.UserFilterCriteria) (model.UserSearchResult, error)
	CountByCriteria(cnxParams string, criteria model.UserFilterCriteria) (int, error)
	Update(cnxParams string, user model.UserInterface) error
	IsLoginInUse(cnxParams string, login string) (int64, string, error)
	IsEmailInUse(cnxParams string, email string) (int64, string, error)
	Delete(cnxParams string, userExtId string) error
}
