package api

import "micro-fiber-test/pkg/model"

type OrganizationServiceInterface interface {
	Create(cnxParams string, defautTenantId int64, organization model.OrganizationInterface) (int64, error)
	Update(cnxParams string, defautTenantId int64, orgCode string, label string) error
	Delete(cnxParams string, defautTenantId int64, orgCode string) error
	FindByCode(cnxParams string, defautTenantId int64, code string) (model.OrganizationInterface, error)
	FindAll(cnxParams string, defautTenantId int64) ([]model.OrganizationInterface, error)
}
