package commons

const (
	ContentTypeHeader      = "Content-Type"
	ContentTypeJson        = "application/json; charset=utf-8"
	TechnicalError         = "technical_error"
	OrgAlreadyExistsByCode = "org_already_exists"
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
