package users

type CreateUserReq struct {
	LastName   string  `json:"lastName" validate:"notblank,maxLength(50)"`
	FirstName  string  `json:"firstName" validate:"notblank,maxLength(50)"`
	MiddleName *string `json:"middleName"`
	Login      string  `json:"login" validate:"notblank,maxLength(50)"`
	Email      string  `json:"email" validate:"notblank,maxLength(50)"`
	Status     int     `json:"status"`
}
