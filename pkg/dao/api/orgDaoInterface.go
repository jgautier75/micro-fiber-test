package api

import (
	"github.com/jackc/pgx/v4"
	"micro-fiber-test/pkg/model"
)

type OrgDaoInterface interface {
	Create(organization model.OrganizationInterface) (int64, error)
	Update(orgCode string, label string) error
	Delete(orgCode string) error
	FindByCode(code string) (model.OrganizationInterface, error)
	FindAll(tenantId int64) ([]model.OrganizationInterface, error)
	ExistsByCode(tenantId int64, code string) (bool, error)
	CreateInTx(tx pgx.Tx, organization model.OrganizationInterface) (int64, error)
}
