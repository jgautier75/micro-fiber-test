package impl

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v4"
	"micro-fiber-test/pkg/commons"
	daoApi "micro-fiber-test/pkg/dao/api"
	"micro-fiber-test/pkg/model"
	svcApi "micro-fiber-test/pkg/service/api"
)

type OrganizationService struct {
	orgDao  daoApi.OrgDaoInterface
	sectDao daoApi.SectorDaoInterface
}

func (orgService *OrganizationService) Create(cnxParams string, defaultTenant int64, organization model.OrganizationInterface) (int64, error) {
	orgExists, err := orgService.orgDao.ExistsByCode(defaultTenant, organization.GetCode())
	if err != nil {
		return 0, err
	}
	if orgExists == false {
		conn, err := pgx.Connect(context.Background(), cnxParams)
		if err != nil {
			return -1, err
		}
		defer func(conn *pgx.Conn, ctx context.Context) {
			err := conn.Close(ctx)
			if err != nil {

			}
		}(conn, context.Background())

		tx, err := conn.BeginTx(context.Background(), pgx.TxOptions{AccessMode: pgx.ReadWrite, IsoLevel: pgx.RepeatableRead})
		if err != nil {
			return 0, err
		}
		defer func() {
			if err != nil {
				err := tx.Rollback(context.Background())
				if err != nil {
					return
				}
			} else {
				err := tx.Commit(context.Background())
				if err != nil {
					return
				}
			}
		}()

		id, err := orgService.orgDao.CreateInTx(tx, organization)
		if err != nil {
			return 0, err
		} else {
			sector := model.Sector{}
			sector.SetLabel(organization.GetLabel())
			sector.SetCode(organization.GetCode())
			sector.SetTenantId(defaultTenant)
			sector.SetSectorStatus(model.SectorStatusActive)
			sector.SetDepth(0)
			sector.SetHasParent(false)
			sector.SetOrgId(id)
			_, err := orgService.sectDao.CreateInTx(tx, &sector)
			if err != nil {
				return 0, err
			}
			return id, nil
		}
	} else {
		return 0, errors.New(commons.OrgAlreadyExistsByCode)
	}
}

func (orgService *OrganizationService) Update(defaultTenant int64, orgCode string, label string) error {
	orgExists, err := orgService.orgDao.ExistsByCode(defaultTenant, orgCode)
	if err != nil {
		return err
	}
	if orgExists == false {
		return errors.New(commons.OrgDoesNotExistByCode)
	}
	return orgService.orgDao.Update(orgCode, label)
}

func (orgService *OrganizationService) Delete(defaultTenant int64, orgCode string) error {
	orgExists, err := orgService.orgDao.ExistsByCode(defaultTenant, orgCode)
	if err != nil {
		return err
	}
	if orgExists == false {
		return errors.New(commons.OrgDoesNotExistByCode)
	}
	org, err := orgService.orgDao.FindByCode(orgCode)
	if err != nil {
		return err
	}
	errSector := orgService.sectDao.DeleteByOrgId(org.GetId())
	if errSector != nil {
		return errSector
	}
	return orgService.orgDao.Delete(orgCode)
}

func (orgService *OrganizationService) FindByCode(defaultTenant int64, code string) (model.OrganizationInterface, error) {
	orgExists, err := orgService.orgDao.ExistsByCode(defaultTenant, code)
	if err != nil {
		return nil, err
	}
	if orgExists == false {
		return nil, errors.New(commons.OrgDoesNotExistByCode)
	}
	return orgService.orgDao.FindByCode(code)
}

func (orgService *OrganizationService) FindAll(defaultTenant int64) ([]model.OrganizationInterface, error) {
	orgs, err := orgService.orgDao.FindAll(defaultTenant)
	if err != nil {
		return nil, err
	}
	return orgs, nil
}
func NewOrgService(orgDao daoApi.OrgDaoInterface, sectorDao daoApi.SectorDaoInterface) svcApi.OrganizationServiceInterface {
	return &OrganizationService{orgDao: orgDao, sectDao: sectorDao}
}
