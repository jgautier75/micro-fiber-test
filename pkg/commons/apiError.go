package commons

const (
	ContentTypeHeader      = "Content-Type"
	ContentTypeJson        = "application/json; charset=utf-8"
	OrgAlreadyExistsByCode = "org_already_exists"
	OrgDoesNotExistByCode  = "org_does_not_exist"
)

type ApiErrorType string

const (
	ErrorTypeFunctional ApiErrorType = "functional"
	ErrorTypeTechnical  ApiErrorType = "technical"
)

type ApiError struct {
	Code    int
	Kind    string
	Message string
}
