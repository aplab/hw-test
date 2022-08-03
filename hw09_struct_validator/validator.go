package hw09structvalidator

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/aplab/hw-test/hw09_struct_validator/validators"
)

const (
	TagName    = "validate"
	RuleNested = "nested"
)

var ErrValueIsNotAStruct = errors.New("value is not a struct")

type ValidationError struct {
	Field string
	Err   error
}

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	b := strings.Builder{}
	for _, e := range v {
		b.WriteString(fmt.Sprintf("field:%v,", e.Field))
		b.WriteString(fmt.Sprintf("error:%v,", e.Err))
	}
	return b.String()
}

func Validate(v interface{}) error {
	return validate(v)
}

func validate(v interface{}) error {
	ref := reflect.TypeOf(v)
	if ref.Kind() != reflect.Struct {
		return ErrValueIsNotAStruct
	}
	val := reflect.ValueOf(v)
	numField := ref.NumField()
	validationErrors := ValidationErrors{}
	for i := 0; i < numField; i++ {
		ve := ValidationErrors{}
		field := ref.Field(i)
		tag := field.Tag
		rule, ok := tag.Lookup(TagName)
		if !ok {
			continue
		}
		value := val.Field(i).Interface()
		err := validateValue(field.Name, rule, value)
		if errors.As(err, &ve) {
			validationErrors = append(validationErrors, ve...)
			continue
		}
		if err != nil {
			return err
		}
	}
	if len(validationErrors) > 0 {
		return validationErrors
	}
	return nil
}

func validateValue(name, rule string, value interface{}) error {
	validationErrors := ValidationErrors{}
	switch v := value.(type) {
	case int:
		rules, err := validators.IntGetRules(rule)
		if err != nil {
			return err
		}
		validationErrors = validateInt(name, rules, v, validationErrors)
	case []int:
		rules, err := validators.IntGetRules(rule)
		if err != nil {
			return err
		}
		validationErrors = validateSliceInt(name, rules, v, validationErrors)
	case string:
		rules, err := validators.StringGetRules(rule)
		if err != nil {
			return err
		}
		validationErrors = validateString(name, rules, v, validationErrors)
	case []string:
		rules, err := validators.StringGetRules(rule)
		if err != nil {
			return err
		}
		validationErrors = validateSliceString(name, rules, v, validationErrors)
	default:
		ref := reflect.TypeOf(value)
		if ref.Kind() == reflect.Struct && rule == RuleNested {
			ve := ValidationErrors{}
			err := validate(value)
			if errors.As(err, &ve) {
				validationErrors = append(validationErrors, ve...)
				break
			}
			if err != nil {
				return err
			}
		}
	}
	if len(validationErrors) > 0 {
		return validationErrors
	}
	return nil
}

func validateSliceString(name string, rules []func(string) error, v []string, e ValidationErrors) ValidationErrors {
	for _, rule := range rules {
		for _, vi := range v {
			err := rule(vi)
			if err != nil {
				e = append(e, ValidationError{
					Field: name,
					Err:   err,
				})
			}
		}
	}
	return e
}

func validateString(name string, rules []func(string) error, v string, e ValidationErrors) ValidationErrors {
	for _, rule := range rules {
		err := rule(v)
		if err != nil {
			e = append(e, ValidationError{
				Field: name,
				Err:   err,
			})
		}
	}
	return e
}

func validateSliceInt(name string, rules []func(int) error, v []int, e ValidationErrors) ValidationErrors {
	for _, rule := range rules {
		for _, vi := range v {
			err := rule(vi)
			if err != nil {
				e = append(e, ValidationError{
					Field: name,
					Err:   err,
				})
			}
		}
	}
	return e
}

func validateInt(name string, rules []func(int) error, v int, e ValidationErrors) ValidationErrors {
	for _, rule := range rules {
		err := rule(v)
		if err != nil {
			e = append(e, ValidationError{
				Field: name,
				Err:   err,
			})
		}
	}
	return e
}
