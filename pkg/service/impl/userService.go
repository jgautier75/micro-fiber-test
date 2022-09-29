package impl

import (
	"micro-fiber-test/pkg/dao/api"
	"micro-fiber-test/pkg/model"
	svcApi "micro-fiber-test/pkg/service/api"
)

type UserService struct {
	dao api.UserDaoInterface
}

func NewUserService(daoP api.UserDaoInterface) svcApi.UserServiceInterface {
	return &UserService{dao: daoP}
}

func (u UserService) Create(cnxParams string, defautTenantId int64, user model.UserInterface) (int64, error) {
	user.SetTenantId(defautTenantId)
	id, createErr := u.dao.Create(cnxParams, user)
	if createErr != nil {
		return 0, createErr
	} else {
		return id, nil
	}
}

func (u UserService) FindByCriteria(cnxParams string, criteria model.UserFilterCriteria) (model.UserSearchResult, error) {
	return u.dao.FindByCriteria(cnxParams, criteria)
}
