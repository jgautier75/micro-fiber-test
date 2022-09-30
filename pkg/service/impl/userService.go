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

func (u UserService) Update(cnxParams string, user model.UserInterface) error {
	return u.dao.Update(cnxParams, user)
}

func (u UserService) FindByCriteria(cnxParams string, criteria model.UserFilterCriteria) (model.UserSearchResult, error) {
	userSearchResult, err := u.dao.FindByCriteria(cnxParams, criteria)
	if err != nil {
		return userSearchResult, err
	}
	cnt, errCount := u.dao.CountByCriteria(cnxParams, criteria)
	if errCount != nil {
		return userSearchResult, err
	}
	userSearchResult.NbResults = cnt
	return userSearchResult, nil
}

func (u UserService) FindByCode(cnxParams string, tenantId int64, orgId int64, externalId string) (model.UserInterface, error) {
	return u.dao.FindByCode(cnxParams, tenantId, orgId, externalId)
}
