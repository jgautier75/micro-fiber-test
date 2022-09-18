package impl

import (
	"errors"
	"micro-fiber-test/pkg/commons"
	daoApi "micro-fiber-test/pkg/dao/api"
	"micro-fiber-test/pkg/model"
	svcApi "micro-fiber-test/pkg/service/api"
)

type OrganizationService struct {
	dao daoApi.OrgDaoInterface
}

func (orgService *OrganizationService) Create(cnxParams string, defaultTenant int64, organization model.OrganizationInterface) (int64, error) {
	orgExists, err := orgService.dao.ExistsByCode(cnxParams, defaultTenant, organization.GetCode())
	if err != nil {
		return 0, err
	}
	if orgExists == false {
		id, createErr := orgService.dao.Create(cnxParams, organization)
		if createErr != nil {
			return 0, err
		} else {
			return id, nil
		}
	} else {
		return 0, errors.New(commons.OrgAlreadyExistsByCode)
	}
}

func (orgService *OrganizationService) Update(cnxParams string, defaultTenant int64, orgCode string, label string) error {
	orgExists, err := orgService.dao.ExistsByCode(cnxParams, defaultTenant, orgCode)
	if err != nil {
		return err
	}
	if orgExists == false {
		return errors.New(commons.OrgDoesNotExistByCode)
	}
	return orgService.dao.Update(cnxParams, orgCode, label)
}

func (orgService *OrganizationService) Delete(cnxParams string, defaultTenant int64, orgCode string) error {
	orgExists, err := orgService.dao.ExistsByCode(cnxParams, defaultTenant, orgCode)
	if err != nil {
		return err
	}
	if orgExists == false {
		return errors.New(commons.OrgDoesNotExistByCode)
	}
	return orgService.dao.Delete(cnxParams, orgCode)
}

func (orgService *OrganizationService) FindByCode(cnxParams string, defaultTenant int64, code string) (model.OrganizationInterface, error) {
	orgExists, err := orgService.dao.ExistsByCode(cnxParams, defaultTenant, code)
	if err != nil {
		return nil, err
	}
	if orgExists == false {
		return nil, errors.New(commons.OrgDoesNotExistByCode)
	}
	return orgService.dao.FindByCode(cnxParams, code)
}

func (orgService *OrganizationService) FindAll(cnxParams string, defaultTenant int64) ([]model.OrganizationInterface, error) {
	orgs, err := orgService.dao.FindAll(cnxParams, defaultTenant)
	if err != nil {
		return nil, err
	}
	return orgs, nil
}

func NewOrgService(daoP daoApi.OrgDaoInterface) svcApi.OrganizationServiceInterface {
	return &OrganizationService{dao: daoP}
}
