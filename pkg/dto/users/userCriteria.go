package users

type UserCriteria struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	Login     string `json:"login"`
}
