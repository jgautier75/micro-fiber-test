package orgs

type CreateOrgRequest struct {
	Label  *string `json:"label" validate:"notblank,maxLength(50)"`
	Kind   *string `json:"type" validate:"notblank"`
	Status int     `json:"status"`
}
