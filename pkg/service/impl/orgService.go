package impl

import (
	"context"
	"errors"
	"fmt"
	"micro-fiber-test/pkg/dto/commons"
	"micro-fiber-test/pkg/model"
	daoApi "micro-fiber-test/pkg/repository/api"
	svcApi "micro-fiber-test/pkg/service/api"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type OrganizationService struct {
	orgDao  daoApi.OrgDaoInterface
	sectDao daoApi.SectorDaoInterface
	dbPool  *pgxpool.Pool
}

func (orgService *OrganizationService) Create(cnxParams string, defaultTenant int64, organization model.Organization) (int64, error) {
	orgExists, err := orgService.orgDao.ExistsByCode(defaultTenant, organization.Code)
	if err != nil {
		return 0, err
	}
	if orgExists {
		return 0, errors.New(commons.OrgAlreadyExistsByCode)
	}

	exists, errOrg := orgService.orgDao.ExistsByLabel(defaultTenant, organization.Label)
	if errOrg != nil {
		return 0, errOrg
	}
	if exists {
		return 0, errors.New(commons.OrgAlreadyExistsByLabel)
	}

	conn, errConnect := pgx.Connect(context.Background(), cnxParams)
	if errConnect != nil {
		return -1, errConnect
	}
	defer func(conn *pgx.Conn, ctx context.Context) {
		errClose := conn.Close(ctx)
		if errClose != nil {
			fullErr := fmt.Errorf("error closing connection [%w]", errClose)
			fmt.Printf("Commit error [%s]", fullErr)
		}
	}(conn, context.Background())

	tx, errTx := conn.BeginTx(context.Background(), pgx.TxOptions{AccessMode: pgx.ReadWrite, IsoLevel: pgx.RepeatableRead})
	if errTx != nil {
		return 0, errTx
	}
	defer func() {
		if errTx != nil {
			errRbk := tx.Rollback(context.Background())
			if errRbk != nil {
				fullErr := fmt.Errorf("error rolling back connection [%w]", errRbk)
				fmt.Printf("Rollback error [%s]", fullErr)
			}
		} else {
			errCmt := tx.Commit(context.Background())
			if errCmt != nil {
				fullErr := fmt.Errorf("error commiting connection [%w]", errCmt)
				fmt.Printf("Commit error [%s]", fullErr)
			}
		}
	}()

	id, errOrgCreateTx := orgService.orgDao.CreateInTx(tx, organization)
	if errOrgCreateTx != nil {
		return 0, errOrgCreateTx
	}
	sector := model.Sector{}
	sector.Label = organization.Label
	sector.Code = organization.Code
	sector.TenantId = defaultTenant
	sector.Status = model.SectorStatusActive
	sector.Depth = 0
	sector.HasParent = false
	sector.OrgId = id
	_, errSect := orgService.sectDao.CreateInTx(tx, sector)
	if errSect != nil {
		return 0, errSect
	}
	return id, nil
}

func (orgService *OrganizationService) Update(defaultTenant int64, orgCode string, label string) error {
	orgExists, err := orgService.orgDao.ExistsByCode(defaultTenant, orgCode)
	if err != nil {
		return err
	}
	if !orgExists {
		return errors.New(commons.OrgDoesNotExistByCode)
	}
	return orgService.orgDao.Update(orgCode, label)
}

func (orgService *OrganizationService) Delete(defaultTenant int64, orgCode string) error {
	orgExists, errExists := orgService.orgDao.ExistsByCode(defaultTenant, orgCode)
	if errExists != nil {
		return errExists
	}
	if !orgExists {
		return errors.New(commons.OrgDoesNotExistByCode)
	}
	org, errFind := orgService.orgDao.FindByCode(orgCode)
	if errFind != nil {
		return errFind
	}
	errSector := orgService.sectDao.DeleteByOrgId(org.Id)
	if errSector != nil {
		return errSector
	}
	return orgService.orgDao.Delete(orgCode)
}

func (orgService *OrganizationService) FindByCode(defaultTenant int64, code string) (model.Organization, error) {
	var nilOrg model.Organization
	orgExists, errExists := orgService.orgDao.ExistsByCode(defaultTenant, code)
	if errExists != nil {
		return nilOrg, errExists
	}
	if !orgExists {
		return nilOrg, errors.New(commons.OrgDoesNotExistByCode)
	}
	return orgService.orgDao.FindByCode(code)
}

func (orgService *OrganizationService) FindAll(defaultTenant int64) ([]model.Organization, error) {
	orgs, err := orgService.orgDao.FindAll(defaultTenant)
	if err != nil {
		return nil, err
	}
	return orgs, nil
}
func NewOrgService(pool *pgxpool.Pool, orgDao daoApi.OrgDaoInterface, sectorDao daoApi.SectorDaoInterface) svcApi.OrganizationServiceInterface {
	return &OrganizationService{orgDao: orgDao, sectDao: sectorDao, dbPool: pool}
}
