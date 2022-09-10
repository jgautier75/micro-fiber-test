package api

import "micro-fiber-test/pkg/model"

type OrgDaoInterface interface {
	Create(cnxParams string, organization model.OrganizationInterface) (int64, error)
	Update(cnxParams string, organization model.OrganizationInterface) error
	Delete(cnxParams string, id int64) error
	FindByCode(cnxParams string, code string) (model.OrganizationInterface, error)
	FindAll(cnxParams string, tenantId int64) ([]model.OrganizationInterface, error)
	ExistsByCode(cnxParams string, tenantId int64, code string) (bool, error)
}
