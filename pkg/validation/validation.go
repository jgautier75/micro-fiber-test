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
	Size  int
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
				val := reflect.ValueOf(obj)
				f := reflect.Indirect(val).FieldByName(fields.Field(inc).Name)
				errorsValidation := validateString(valRules, f.String(), fields.Field(inc).Name)
				if len(errorsValidation) > 0 {
					for _, e := range errorsValidation {
						errorsList = append(errorsList, e)
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
				} else {
					errorsValidation := validateString(valRules, f.Elem().String(), fields.Field(inc).Name)
					if len(errorsValidation) > 0 {
						for _, e := range errorsValidation {
							errorsList = append(errorsList, e)
						}
					}
				}
			}
		}
	}
	return errorsList
}

func validateString(valRules []string, fieldValue string, fieldName string) []ErrorValidation {
	var errorsList []ErrorValidation
	for _, va := range valRules {
		if va == ValidRuleNotBlank {
			if len(strings.TrimSpace(fieldValue)) <= 0 {
				validationError := ErrorValidation{
					Field: fieldName,
					Error: errors.New(ValidErrorNotBlank),
				}
				errorsList = append(errorsList, validationError)
			}
		} else if strings.Contains(va, ValidateRuleMaxLength) {
			l, err := extractLengthFromRule(va)
			if err == nil {
				if len(strings.TrimSpace(fieldValue)) > l {
					validationError := ErrorValidation{
						Field: fieldName,
						Error: errors.New(ValidErrorMaxLength),
						Size:  l,
					}
					errorsList = append(errorsList, validationError)
				}
			}
		} else if strings.Contains(va, ValidateRuleMinLength) {
			l, err := extractLengthFromRule(va)
			if err == nil {
				if len(strings.TrimSpace(fieldValue)) < l {
					validationError := ErrorValidation{
						Field: fieldName,
						Error: errors.New(ValidErrorMinLength),
						Size:  l,
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
