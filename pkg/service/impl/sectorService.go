package impl

import (
	"errors"
	"micro-fiber-test/pkg/dto/commons"
	"micro-fiber-test/pkg/model"
	"micro-fiber-test/pkg/repository/api"
	svcApi "micro-fiber-test/pkg/service/api"
)

type SectorService struct {
	dao api.SectorDaoInterface
}

func NewSectorService(daoP api.SectorDaoInterface) svcApi.SectorServiceInterface {
	return &SectorService{dao: daoP}
}

func (sectorSvc SectorService) Create(defautTenantId int64, sector model.Sector) (int64, error) {
	sector.TenantId = defautTenantId
	id, _, err := sectorSvc.dao.FindByLabel(defautTenantId, sector.Label)
	if err != nil {
		return 0, err
	}
	if id > 0 {
		return 0, errors.New(commons.SectorAlreadyExist)
	}

	id, createErr := sectorSvc.dao.Create(sector)
	if createErr != nil {
		return 0, createErr
	} else {
		return id, nil
	}
}

func (sectorSvc SectorService) FindSectorsByTenantOrg(defaultTenantId int64, orgId int64) ([]model.Sector, error) {
	sectors, err := sectorSvc.dao.FindSectorsByTenantOrg(defaultTenantId, orgId)
	if err != nil {
		return nil, err
	}
	return sectors, nil
}

func (sectorSvc SectorService) FindByCode(defaultTenantId int64, code string) (model.Sector, error) {
	return sectorSvc.dao.FindByCode(defaultTenantId, code)
}

func (sectorSvc SectorService) FindRootSectorId(defaultTenantId int64, orgId int64) (int64, error) {
	return sectorSvc.dao.FindRootSector(defaultTenantId, orgId)
}

func (sectorSvc SectorService) DeleteSector(defaultTenantId int64, sectorId int64) error {
	return sectorSvc.dao.DeleteSector(defaultTenantId, sectorId)
}

func (sectorSvc SectorService) Update(defaultTenantId int64, id int64, label string) error {
	return sectorSvc.dao.Update(defaultTenantId, id, label)
}

func (sectorSvc SectorService) FindByLabel(defaultTenantId int64, label string) (int64, string, error) {
	return sectorSvc.dao.FindByLabel(defaultTenantId, label)
}
