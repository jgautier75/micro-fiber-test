package validation

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

type OrgTest struct {
	Code   string  `json:"code" validate:"notblank,maxLength(50)"`
	Label  *string `json:"label" validate:"notblank,maxLength(50)"`
	Kind   string  `json:"type" validate:"notblank"`
	Status int     `json:"status"`
}

func TestValidation(t *testing.T) {
	orgTest := OrgTest{
		Code: "code_test",
		Kind: "community",
	}
	errors := Validate(orgTest)
	assert.Truef(t, len(errors) == 1, "One error")
	assert.Truef(t, errors[0].Error.Error() == ValidErrorNotBlank && errors[0].Field == "Label", fmt.Sprintf("Constraint [%s] for field [%s]", ValidErrorNotBlank, "label"))
}
