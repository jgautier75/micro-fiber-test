package orgs

import "micro-fiber-test/pkg/model"

type CreateSectorReq struct {
	Label      *string `json:"label" validate:"required,max=50"`
	ParentCode string  `json:"parentCode"`
	Status     int     `json:"status"`
}

func ConvertSectorReqToDaoModel(defaultTenantId int64, sectorReq CreateSectorReq) model.Sector {
	sect := model.Sector{}
	sect.TenantId = defaultTenantId
	if sectorReq.Label != nil {
		sect.Label = *sectorReq.Label
	}
	sect.Status = model.SectorStatus(sectorReq.Status)
	return sect
}
