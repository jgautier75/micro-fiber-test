package validation

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

type OrgTest struct {
	Label string `json:"label" validate:"notblank,maxLength(2)"`
	Code  string `json:"code" validate:"required"`
	Kind  int    `json:"int"`
}

func TestValidation(t *testing.T) {
	orgTest := OrgTest{
		Code:  "code_test",
		Label: "",
		Kind:  0,
	}
	errors := Validate(orgTest)
	assert.Truef(t, len(errors) == 1, "One error")
	assert.Truef(t, errors[0].Error.Error() == ValidErrorNotBlank && errors[0].Field == "Label", fmt.Sprintf("Constraint [%s] for field [%s]", ValidErrorNotBlank, "label"))
}
