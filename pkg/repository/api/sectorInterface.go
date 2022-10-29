package api

import (
	"github.com/jackc/pgx/v5"
	"micro-fiber-test/pkg/model"
)

type SectorDaoInterface interface {
	Create(sector model.SectorInterface) (int64, error)
	DeleteByOrgId(orgId int64) error
	Update(defaultTenantId int64, id int64, label string) error
	CreateInTx(tx pgx.Tx, sector model.SectorInterface) (int64, error)
	FindSectorsByTenantOrg(defaultTenantId int64, orgId int64) ([]model.SectorInterface, error)
	FindByCode(defaultTenantId int64, code string) (model.SectorInterface, error)
	FindRootSector(defaultTenantId int64, orgId int64) (int64, error)
	DeleteSector(defaultTenantId int64, sectorId int64) error
	FindByLabel(defaultTenantId int64, label string) (int64, string, error)
}
