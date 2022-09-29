package users

type UserListResponse struct {
	Users []UserResponse `json:"users,omitempty"`
}
