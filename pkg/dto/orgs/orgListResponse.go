package orgs

type OrganizationListResponse struct {
	Organizations []OrganizationResponse `json:"organizations,omitempty"`
}

type OrganizationResponse struct {
	Code   string `json:"code"`
	Label  string `json:"label"`
	Kind   string `json:"type"`
	Status int    `json:"status"`
}
