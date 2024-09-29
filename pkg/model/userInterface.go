package model

type UserStatus int64

const (
	UserStatusDraft    UserStatus = 0
	UserStatusActive   UserStatus = 1
	UserStatusInactive UserStatus = 2
	UserStatusDeleted  UserStatus = 3
)
