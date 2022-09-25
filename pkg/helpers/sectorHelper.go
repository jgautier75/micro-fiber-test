package helpers

import (
	"errors"
	"micro-fiber-test/pkg/commons"
	dtos "micro-fiber-test/pkg/dto/sectors"
)

func BuildSectorsHierarchy(sectors []dtos.SectorResponse) (dtos.SectorResponse, error) {
	var rootSector dtos.SectorResponse
	for _, sector := range sectors {
		if sector.Depth == 0 {
			rootSector = sector
			break
		}
	}
	if &rootSector == nil {
		return rootSector, errors.New(commons.SectorRootNotFound)
	}
	return fetchRecursively(&rootSector, sectors), nil
}

func fetchRecursively(parentSector *dtos.SectorResponse, sectors []dtos.SectorResponse) dtos.SectorResponse {
	var c = make([]dtos.SectorResponse, 0)
	for _, sector := range sectors {
		if sector.ParentId == parentSector.Id {
			c = append(c, fetchRecursively(&sector, sectors))
		}
	}
	return dtos.SectorResponse{
		Id:       parentSector.Id,
		Code:     parentSector.Code,
		Label:    parentSector.Label,
		ParentId: parentSector.ParentId,
		Depth:    parentSector.Depth,
		Status:   parentSector.Status,
		Children: c,
	}
}
