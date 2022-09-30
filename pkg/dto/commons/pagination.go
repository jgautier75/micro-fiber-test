package commons

type Pagination struct {
	Page       int `json:"page"`
	NbPages    int `json:"nbPages"`
	TotalCount int `json:"totalCount"`
}
