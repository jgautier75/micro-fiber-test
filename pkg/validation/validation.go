package validation

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

type ErrorType string

const (
	RuleRangeLower        string = "("
	RuleRangeUpper        string = ")"
	RulesSeparator        string = ","
	ValidRulePattern      string = "validate"
	ValidRuleNotBlank     string = "notblank"
	ValidateRuleMaxLength string = "maxLength"
	ValidateRuleMinLength string = "minLength"
	ValidErrorNotBlank    string = "error_not_blank"
	ValidErrorMaxLength   string = "error_max_length"
)

type ErrorValidation struct {
	Field string
	Error error
}

func Validate(obj any) []ErrorValidation {
	fields := reflect.TypeOf(obj)
	nbFields := fields.NumField()
	var errorsList []ErrorValidation
	for inc := 0; inc < nbFields; inc++ {
		if val, ok := fields.Field(inc).Tag.Lookup(ValidRulePattern); ok && val != "" {
			if fields.Field(inc).Type.String() == "string" {
				valRules := strings.Split(fields.Field(inc).Tag.Get(ValidRulePattern), RulesSeparator)
				for _, va := range valRules {
					if va == ValidRuleNotBlank {
						val := reflect.ValueOf(obj)
						f := reflect.Indirect(val).FieldByName(fields.Field(inc).Name)
						if len(strings.TrimSpace(f.String())) <= 0 {
							validationError := ErrorValidation{
								Field: fields.Field(inc).Name,
								Error: errors.New(ValidErrorNotBlank),
							}
							errorsList = append(errorsList, validationError)
						}
					} else if strings.Contains(va, ValidateRuleMaxLength) {
						posLeftParenthesis := strings.Index(va, RuleRangeLower)
						posRightParenthesis := strings.Index(va, RuleRangeUpper)
						if posLeftParenthesis > 0 && posRightParenthesis > 0 {
							maxLengthStr := va[posLeftParenthesis+1 : posRightParenthesis]
							maxLength, errConvert := strconv.Atoi(maxLengthStr)
							if errConvert == nil {
								val := reflect.ValueOf(obj)
								f := reflect.Indirect(val).FieldByName(fields.Field(inc).Name)
								if len(strings.TrimSpace(f.String())) > maxLength {
									validationError := ErrorValidation{
										Field: fields.Field(inc).Name,
										Error: errors.New(ValidErrorMaxLength),
									}
									errorsList = append(errorsList, validationError)
								}
							}
						} else {
							fmt.Printf("Invalid validation expression [%s]", va)
						}
					}
				}
			}
		}
	}
	return errorsList
}
