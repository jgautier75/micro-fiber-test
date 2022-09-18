package api

import "micro-fiber-test/pkg/model"

type OrganizationServiceInterface interface {
	Create(cnxParams string, organization model.OrganizationInterface) (int64, error)
	Update(cnxParams string, orgCode string, label string) error
	Delete(cnxParams string, orgCode string) error
	FindByCode(cnxParams string, code string) (model.OrganizationInterface, error)
	FindAll(cnxParams string) ([]model.OrganizationInterface, error)
}
