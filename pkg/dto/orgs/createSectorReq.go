package orgs

import "micro-fiber-test/pkg/model"

type CreateSectorReq struct {
	Code       *string `json:"code" validate:"notblank,maxLength(50)"`
	Label      *string `json:"label" validate:"notblank,maxLength(50)"`
	ParentCode string  `json:"parentCode"`
}

func ConvertSectorReqToDaoModel(defaultTenantId int64, sectorReq CreateSectorReq) model.Sector {
	sect := model.Sector{}
	sect.SetTenantId(defaultTenantId)
	if sectorReq.Code != nil {
		sect.SetCode(*sectorReq.Code)
	}
	if sectorReq.Label != nil {
		sect.SetLabel(*sectorReq.Label)
	}
	return sect
}
