package helpers

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"micro-fiber-test/pkg/dto/sectors"
	"micro-fiber-test/pkg/model"
	"testing"
)

func TestSectorHierarchy(t *testing.T) {
	var secList []sectors.SectorResponse
	root := sectors.SectorResponse{
		Id:     1,
		Status: model.SectorStatusActive,
		Depth:  0,
		Code:   "root",
		Label:  "root",
	}
	secList = append(secList, root)
	north := sectors.SectorResponse{
		Id:       2,
		ParentId: 1,
		Status:   model.SectorStatusActive,
		Depth:    1,
		Code:     "north",
		Label:    "north",
	}
	secList = append(secList, north)
	northEast := sectors.SectorResponse{
		Id:       3,
		ParentId: 2,
		Status:   model.SectorStatusActive,
		Depth:    1,
		Code:     "north-east",
		Label:    "north-east",
	}
	secList = append(secList, northEast)
	northWest := sectors.SectorResponse{
		Id:       4,
		ParentId: 2,
		Status:   model.SectorStatusActive,
		Depth:    1,
		Code:     "north-west",
		Label:    "north-west",
	}
	secList = append(secList, northWest)

	south := sectors.SectorResponse{
		Id:       5,
		ParentId: 1,
		Status:   model.SectorStatusActive,
		Depth:    1,
		Code:     "south",
		Label:    "south",
	}
	secList = append(secList, south)
	southEast := sectors.SectorResponse{
		Id:       6,
		ParentId: 5,
		Status:   model.SectorStatusActive,
		Depth:    1,
		Code:     "south-east",
		Label:    "south-east",
	}
	secList = append(secList, southEast)
	southWest := sectors.SectorResponse{
		Id:       7,
		ParentId: 5,
		Status:   model.SectorStatusActive,
		Depth:    1,
		Code:     "south-west",
		Label:    "south-west",
	}
	secList = append(secList, southWest)
	sectorHierarchy, err := BuildSectorsHierarchy(secList)
	if err != nil {
		fmt.Printf("Sector error [%v]", err)
	} else {
		fmt.Printf("Sector hierarchy [%v]", sectorHierarchy)
		assert.Truef(t, sectorHierarchy.Label == "root", "First sector [%s]", sectorHierarchy.Label)
		assert.Truef(t, len(sectorHierarchy.Children) == 2, "[%d] children", 2)
	}
}
