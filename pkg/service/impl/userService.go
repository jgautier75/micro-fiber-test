package impl

import (
	"errors"
	"micro-fiber-test/pkg/commons"
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

	// Login is unique
	idUsr, _, errLogin := u.dao.IsLoginInUse(cnxParams, user.GetLogin())
	if errLogin != nil {
		return 0, errLogin
	}
	if idUsr > 0 {
		return 0, errors.New(commons.UserLoginAlreadyInUse)
	}

	// Email is unique
	idUsr, _, errEmail := u.dao.IsEmailInUse(cnxParams, user.GetEmail())
	if errEmail != nil {
		return 0, errEmail
	}
	if idUsr > 0 {
		return 0, errors.New(commons.UserEmailAlreadyInUse)
	}

	id, createErr := u.dao.Create(cnxParams, user)
	if createErr != nil {
		return 0, createErr
	} else {
		return id, nil
	}
}

func (u UserService) Update(cnxParams string, user model.UserInterface) error {

	// Login is unique
	_, extId, errLogin := u.dao.IsLoginInUse(cnxParams, user.GetLogin())
	if errLogin != nil {
		return errLogin
	}
	if extId != user.GetExternalId() {
		return errors.New(commons.UserLoginAlreadyInUse)
	}

	// Email is unique
	_, extId, errEmail := u.dao.IsEmailInUse(cnxParams, user.GetEmail())
	if errEmail != nil {
		return errEmail
	}
	if extId != user.GetExternalId() {
		return errors.New(commons.UserEmailAlreadyInUse)
	}

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
	return u.dao.FindByExternalId(cnxParams, tenantId, orgId, externalId)
}

func (u UserService) Delete(cnxParams string, externalId string) error {
	return u.dao.Delete(cnxParams, externalId)
}
