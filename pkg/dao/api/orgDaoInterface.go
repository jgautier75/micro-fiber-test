package api

import (
	"github.com/jackc/pgx/v4"
	"micro-fiber-test/pkg/model"
)

type OrgDaoInterface interface {
	Create(cnxParams string, organization model.OrganizationInterface) (int64, error)
	Update(cnxParams string, orgCode string, label string) error
	Delete(cnxParams string, orgCode string) error
	FindByCode(cnxParams string, code string) (model.OrganizationInterface, error)
	FindAll(cnxParams string, tenantId int64) ([]model.OrganizationInterface, error)
	ExistsByCode(cnxParams string, tenantId int64, code string) (bool, error)
	CreateInTx(tx pgx.Tx, organization model.OrganizationInterface) (int64, error)
}
