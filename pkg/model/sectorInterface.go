package model

type SectorStatus int64

const (
	SectorStatusDraft    SectorStatus = 0
	SectorStatusActive   SectorStatus = 1
	SectorStatusInactive SectorStatus = 2
	SectorStatusDeleted  SectorStatus = 3
)
