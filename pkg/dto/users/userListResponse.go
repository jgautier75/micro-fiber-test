package users

import "micro-fiber-test/pkg/dto/commons"

type UserListResponse struct {
	Pagination commons.Pagination `json:"pagination,omitempty"`
	Users      []UserResponse     `json:"users,omitempty"`
}
