package validation

import (
	"fmt"
	"github.com/go-playground/validator"
	"testing"
)

type OrgPlay struct {
	Code string `validate:"required,min=3,max=32"`
}

func TestValidation(t *testing.T) {
	orgPlay := OrgPlay{
		Code: "12",
	}

	var validate = validator.New()
	err := validate.Struct(orgPlay)
	var errors []*ErrorResponse
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			var element ErrorResponse
			element.FailedField = err.StructNamespace()
			element.Tag = err.Tag()
			element.Value = err.Param()
			errors = append(errors, &element)
		}
	}

	for _, e := range errors {
		fmt.Printf("Field: [%s], Tag: [%s], Value: [%s]", e.FailedField, e.Tag, e.Value)
	}

}
