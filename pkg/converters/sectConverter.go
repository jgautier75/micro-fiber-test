package converters

import (
	"micro-fiber-test/pkg/dto/sectors"
	"micro-fiber-test/pkg/model"
)

func ConvertSectorModelToSectorResp(sect model.SectorInterface) sectors.SectorResponse {
	sectorResponse := sectors.SectorResponse{
		Id:     sect.GetId(),
		Code:   sect.GetCode(),
		Label:  sect.GetLabel(),
		Depth:  sect.GetDepth(),
		Status: sect.GetSectorStatus(),
	}
	if sect.GetHasParent() {
		sectorResponse.ParentId = sect.GetParentId().Int64
	}
	return sectorResponse
}
