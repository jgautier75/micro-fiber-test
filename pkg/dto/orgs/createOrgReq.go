package orgs

type CreateOrgRequest struct {
	Code   *string `json:"code" validate:"notblank,maxLength(50)"`
	Label  *string `json:"label" validate:"notblank,maxLength(50)"`
	Kind   *string `json:"type" validate:"notblank"`
	Status int     `json:"status"`
}
