package api

import "micro-fiber-test/pkg/model"

type SectorDaoInterface interface {
	Create(cnxParams string, sector model.SectorInterface) (int64, error)
	DeleteByOrgId(cnxParams string, orgId int64) error
}
