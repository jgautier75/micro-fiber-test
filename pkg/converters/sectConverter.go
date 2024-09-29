package converters

import (
	"micro-fiber-test/pkg/dto/sectors"
	"micro-fiber-test/pkg/model"
)

func ConvertSectorModelToSectorResp(sect model.Sector) sectors.SectorResponse {
	sectorResponse := sectors.SectorResponse{
		Id:     sect.Id,
		Code:   sect.Code,
		Label:  sect.Label,
		Depth:  sect.Depth,
		Status: sect.Status,
	}
	if sect.HasParent {
		sectorResponse.ParentId = sect.ParentId.Int64
	}
	return sectorResponse
}
