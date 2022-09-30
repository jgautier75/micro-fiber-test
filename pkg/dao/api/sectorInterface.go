package api

import (
	"github.com/jackc/pgx/v4"
	"micro-fiber-test/pkg/model"
)

type SectorDaoInterface interface {
	Create(cnxParams string, sector model.SectorInterface) (int64, error)
	DeleteByOrgId(cnxParams string, orgId int64) error
	CreateInTx(tx pgx.Tx, sector model.SectorInterface) (int64, error)
	FindSectorsByTenantOrg(cnxParams string, defaultTenantId int64, orgId int64) ([]model.SectorInterface, error)
	FindByCode(cnxParams string, defaultTenantId int64, code string) (model.SectorInterface, error)
	FindRootSector(cnxParams string, defaultTenantId int64, orgId int64) (int64, error)
	DeleteSector(cnxParams string, defaultTenantId int64, sectorId int64) error
	FindByLabel(cnxParams string, defaultTenantId int64, label string) (int64, string, error)
}
