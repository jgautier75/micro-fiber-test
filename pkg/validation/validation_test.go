package validation

import (
	"fmt"
	"testing"
)

type OrgTest struct {
	Label string `json:"label" validate:"notblank,maxLength(2)"`
	Code  string `json:"code" validate:"required"`
	Kind  int    `json:"int"`
}

func TestDao(t *testing.T) {
	orgTest := OrgTest{
		Code:  "code_test",
		Label: "",
		Kind:  0,
	}
	errors := Validate(orgTest)
	for _, errValid := range errors {
		fmt.Printf("Validation error for field [%s]:  [%v]", errValid.Field, errValid.Error)
	}
	fmt.Println("Finished")
}
