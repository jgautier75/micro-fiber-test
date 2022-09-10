package api

import "micro-fiber-test/pkg/model"

type OrganizationServiceInterface interface {
	Create(cnxParams string, organization model.OrganizationInterface) (int64, error)
	Update(cnxParams string, org model.OrganizationInterface) error
	Delete(cnxParams string, id int64) error
	FindByCode(cnxParams string, code string) (model.OrganizationInterface, error)
	FindAll(cnxParams string) ([]model.OrganizationInterface, error)
}
