package commons

const (
	ContentTypeHeader      = "Content-Type"
	ContentTypeJson        = "application/json; charset=utf-8"
	OrgAlreadyExistsByCode = "org_already_exists"
	OrgDoesNotExistByCode  = "org_does_not_exist"
	OrgNotFound            = "org_not_found"
	SectorAlreadyExist     = "sector_already_exists"
	SectorRootNotFound     = "sector_root_not_found"
	SectorNotFound         = "sector_not_found"
)

type ApiErrorType string

const (
	ErrorTypeFunctional ApiErrorType = "functional"
	ErrorTypeTechnical  ApiErrorType = "technical"
)

type ApiError struct {
	Code    int               `json:"code"`
	Kind    string            `json:"kind"`
	Message string            `json:"message"`
	Details []ApiErrorDetails `json:"details,omitempty"`
}

type ApiErrorDetails struct {
	Field  string `json:"field"`
	Detail string `json:"detail"`
}
