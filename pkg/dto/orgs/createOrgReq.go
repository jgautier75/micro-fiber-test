package orgs

type CreateOrgRequest struct {
	Label  *string `json:"label" validate:"required,max=50"`
	Kind   *string `json:"type" validate:"required"`
	Status int     `json:"status" validate:"required"`
}
