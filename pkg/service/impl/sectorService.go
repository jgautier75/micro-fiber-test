package impl

import (
	"micro-fiber-test/pkg/dao/api"
	"micro-fiber-test/pkg/model"
	svcApi "micro-fiber-test/pkg/service/api"
)

type SectorService struct {
	dao api.SectorDaoInterface
}

func NewSectorService(daoP api.SectorDaoInterface) svcApi.SectorServiceInterface {
	return &SectorService{dao: daoP}
}

func (sectorSvc SectorService) Create(cnxParams string, defautTenantId int64, sector model.SectorInterface) (int64, error) {
	sector.SetTenantId(defautTenantId)
	id, createErr := sectorSvc.dao.Create(cnxParams, sector)
	if createErr != nil {
		return 0, createErr
	} else {
		return id, nil
	}
}
