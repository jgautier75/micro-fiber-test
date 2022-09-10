package impl

import (
	"errors"
	"micro-fiber-test/pkg/commons"
	daoApi "micro-fiber-test/pkg/dao/api"
	"micro-fiber-test/pkg/model"
	svcApi "micro-fiber-test/pkg/service/api"
)

const defaultTenant int64 = 1

type OrganizationService struct {
	dao daoApi.OrgDaoInterface
}

func (orgService *OrganizationService) Create(cnxParams string, organization model.OrganizationInterface) (int64, error) {
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

func (orgService *OrganizationService) Update(cnxParams string, org model.OrganizationInterface) error {
	return orgService.dao.Update(cnxParams, org)
}

func (orgService *OrganizationService) Delete(cnxParams string, id int64) error {
	return orgService.dao.Delete(cnxParams, id)
}

func (orgService *OrganizationService) FindByCode(cnxParams string, code string) (model.OrganizationInterface, error) {
	return orgService.dao.FindByCode(cnxParams, code)
}

func (orgService *OrganizationService) FindAll(cnxParams string) ([]model.OrganizationInterface, error) {
	return orgService.dao.FindAll(cnxParams, defaultTenant)
}

func NewOrgService(daoP daoApi.OrgDaoInterface) svcApi.OrganizationServiceInterface {
	return &OrganizationService{dao: daoP}
}
