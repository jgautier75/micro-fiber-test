package commons

type Pagination struct {
	Page  int `json:"page"`
	Total int `json:"total"`
	Count int `json:"count"`
}
