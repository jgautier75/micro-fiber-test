package api

import "micro-fiber-test/pkg/model"

type SectorServiceInterface interface {
	Create(cnxParams string, defautTenantId int64, sector model.SectorInterface) (int64, error)
	FindSectorsByTenantOrg(cnxParams string, defaultTenantId int64, orgId int64) ([]model.SectorInterface, error)
}
