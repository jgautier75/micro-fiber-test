package api

import "micro-fiber-test/pkg/model"

type SectorServiceInterface interface {
	Create(defautTenantId int64, sector model.SectorInterface) (int64, error)
	FindSectorsByTenantOrg(defaultTenantId int64, orgId int64) ([]model.SectorInterface, error)
	FindByCode(defaultTenantId int64, code string) (model.SectorInterface, error)
	FindRootSectorId(defaultTenantId int64, orgId int64) (int64, error)
	DeleteSector(defaultTenantId int64, sectorId int64) error
}
