package impl

import (
	"errors"
	"micro-fiber-test/pkg/commons"
	daoApi "micro-fiber-test/pkg/dao/api"
	"micro-fiber-test/pkg/model"
	svcApi "micro-fiber-test/pkg/service/api"
)

type OrganizationService struct {
	orgDao  daoApi.OrgDaoInterface
	sectDao daoApi.SectorDaoInterface
}

func (orgService *OrganizationService) Create(cnxParams string, defaultTenant int64, organization model.OrganizationInterface) (int64, error) {
	orgExists, err := orgService.orgDao.ExistsByCode(cnxParams, defaultTenant, organization.GetCode())
	if err != nil {
		return 0, err
	}
	if orgExists == false {
		id, createErr := orgService.orgDao.Create(cnxParams, organization)
		if createErr != nil {
			return 0, err
		} else {
			sector := model.Sector{}
			sector.SetLabel(organization.GetLabel())
			sector.SetCode(organization.GetCode())
			sector.SetTenantId(defaultTenant)
			sector.SetSectorStatus(model.SectorStatusActive)
			sector.SetDepth(0)
			sector.SetHasParent(false)
			sector.SetOrgId(id)
			_, errSector := orgService.sectDao.Create(cnxParams, &sector)
			if errSector != nil {
				return 0, err
			}
			return id, nil
		}
	} else {
		return 0, errors.New(commons.OrgAlreadyExistsByCode)
	}
}

func (orgService *OrganizationService) Update(cnxParams string, defaultTenant int64, orgCode string, label string) error {
	orgExists, err := orgService.orgDao.ExistsByCode(cnxParams, defaultTenant, orgCode)
	if err != nil {
		return err
	}
	if orgExists == false {
		return errors.New(commons.OrgDoesNotExistByCode)
	}
	return orgService.orgDao.Update(cnxParams, orgCode, label)
}

func (orgService *OrganizationService) Delete(cnxParams string, defaultTenant int64, orgCode string) error {
	orgExists, err := orgService.orgDao.ExistsByCode(cnxParams, defaultTenant, orgCode)
	if err != nil {
		return err
	}
	if orgExists == false {
		return errors.New(commons.OrgDoesNotExistByCode)
	}
	org, err := orgService.orgDao.FindByCode(cnxParams, orgCode)
	if err != nil {
		return err
	}
	errSector := orgService.sectDao.DeleteByOrgId(cnxParams, org.GetId())
	if errSector != nil {
		return errSector
	}
	return orgService.orgDao.Delete(cnxParams, orgCode)
}

func (orgService *OrganizationService) FindByCode(cnxParams string, defaultTenant int64, code string) (model.OrganizationInterface, error) {
	orgExists, err := orgService.orgDao.ExistsByCode(cnxParams, defaultTenant, code)
	if err != nil {
		return nil, err
	}
	if orgExists == false {
		return nil, errors.New(commons.OrgDoesNotExistByCode)
	}
	return orgService.orgDao.FindByCode(cnxParams, code)
}

func (orgService *OrganizationService) FindAll(cnxParams string, defaultTenant int64) ([]model.OrganizationInterface, error) {
	orgs, err := orgService.orgDao.FindAll(cnxParams, defaultTenant)
	if err != nil {
		return nil, err
	}
	return orgs, nil
}
func NewOrgService(orgDao daoApi.OrgDaoInterface, sectorDao daoApi.SectorDaoInterface) svcApi.OrganizationServiceInterface {
	return &OrganizationService{orgDao: orgDao, sectDao: sectorDao}
}
