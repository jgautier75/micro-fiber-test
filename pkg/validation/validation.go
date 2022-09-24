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
	ValidErrorMinLength   string = "error_min_length"
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
			valRules := strings.Split(fields.Field(inc).Tag.Get(ValidRulePattern), RulesSeparator)
			switch fields.Field(inc).Type.Kind() {
			case reflect.String:
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
						l, err := extractLengthFromRule(va)
						if err == nil {
							val := reflect.ValueOf(obj)
							f := reflect.Indirect(val).FieldByName(fields.Field(inc).Name)
							if len(strings.TrimSpace(f.String())) > l {
								validationError := ErrorValidation{
									Field: fields.Field(inc).Name,
									Error: errors.New(ValidErrorMaxLength),
								}
								errorsList = append(errorsList, validationError)
							}
						}
					} else if strings.Contains(va, ValidateRuleMinLength) {
						l, err := extractLengthFromRule(va)
						if err == nil {
							val := reflect.ValueOf(obj)
							f := reflect.Indirect(val).FieldByName(fields.Field(inc).Name)
							if len(strings.TrimSpace(f.String())) < l {
								validationError := ErrorValidation{
									Field: fields.Field(inc).Name,
									Error: errors.New(ValidErrorMinLength),
								}
								errorsList = append(errorsList, validationError)
							}
						}
					}
				}
			case reflect.Pointer:
				val := reflect.ValueOf(obj)
				f := reflect.Indirect(val).FieldByName(fields.Field(inc).Name)
				if f.IsNil() {
					validationError := ErrorValidation{
						Field: fields.Field(inc).Name,
						Error: errors.New(ValidErrorNotBlank),
					}
					errorsList = append(errorsList, validationError)
				}
			}
		}
	}
	return errorsList
}

func extractLengthFromRule(validationRule string) (int, error) {
	posLeftParenthesis := strings.Index(validationRule, RuleRangeLower)
	posRightParenthesis := strings.Index(validationRule, RuleRangeUpper)
	if posLeftParenthesis > 0 && posRightParenthesis > 0 {
		maxLengthStr := validationRule[posLeftParenthesis+1 : posRightParenthesis]
		return strconv.Atoi(maxLengthStr)
	} else {
		return 0, errors.New(fmt.Sprintf("unable to extract length from validation rule [%s]", validationRule))
	}
}
