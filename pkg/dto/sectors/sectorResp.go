package sectors

import (
	"micro-fiber-test/pkg/model"
)

type SectorResponse struct {
	Id       int64              `json:"-"`
	Code     string             `json:"code"`
	Label    string             `json:"label"`
	Depth    int                `json:"depth"`
	ParentId int64              `json:"-"`
	Status   model.SectorStatus `json:"status"`
	Children []SectorResponse   `json:"children,omitempty"`
}
