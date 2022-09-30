package users

type UserResponse struct {
	ExternalId string `json:"id"`
	LastName   string `json:"lastName"`
	FirstName  string `json:"firstName"`
	MiddleName string `json:"middleName,omitempty"`
	Login      string `json:"login"`
	Email      string `json:"email"`
	Status     int    `json:"status"`
}
