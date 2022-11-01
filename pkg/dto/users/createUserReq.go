package users

type CreateUserReq struct {
	LastName   string  `json:"lastName" validate:"required,max=50"`
	FirstName  string  `json:"firstName" validate:"required,max=50"`
	MiddleName *string `json:"middleName"`
	Login      string  `json:"login" validate:"required,max=50"`
	Email      string  `json:"email" validate:"required,max=50"`
	Status     int     `json:"status"`
}

type UpdateUserReq struct {
	LastName   string  `json:"lastName" validate:"required,max=50"`
	FirstName  string  `json:"firstName" validate:"required,max=50"`
	MiddleName *string `json:"middleName"`
	Login      string  `json:"login" validate:"required,max=50"`
	Email      string  `json:"email" validate:"required,max=50"`
}
