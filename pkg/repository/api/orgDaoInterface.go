package api

import (
	"micro-fiber-test/pkg/model"

	"github.com/jackc/pgx/v5"
)

type OrgDaoInterface interface {
	Create(organization model.Organization) (int64, error)
	Update(orgCode string, label string) error
	Delete(orgCode string) error
	FindByCode(code string) (model.Organization, error)
	FindAll(tenantId int64) ([]model.Organization, error)
	ExistsByCode(tenantId int64, code string) (bool, error)
	ExistsByLabel(tenantId int64, label string) (bool, error)
	CreateInTx(tx pgx.Tx, organization model.Organization) (int64, error)
}
