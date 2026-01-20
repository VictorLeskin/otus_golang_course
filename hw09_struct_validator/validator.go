package hw09structvalidator

import (
	"fmt"
	"reflect"
	"strings"
)

type ValidationError struct {
	Field string
	Err   error
}

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	panic("implement me")
}

var ErrArgumentNotStructure = fmt.Errorf("argument is not a struct")

type CValidator struct {
	in interface{}
	rv reflect.Value
	rt reflect.Type
}

func (v *CValidator) Validate0() error {
	v.rv = reflect.ValueOf(v.in)
	v.rt = v.rv.Type()

	return v.validateStruct()
}

func (v *CValidator) validateStruct() error {
	if v.rt.Kind() == reflect.Struct {
		for i := 0; i < v.rt.NumField(); i++ {
			v.validateStructField(v.rt.Field(i))
		}
	} else {
		return ErrArgumentNotStructure
	}

	return nil
}

func (v *CValidator) validateStructField(field reflect.StructField) error {
	// get validate tag.
	validateTag := field.Tag.Get("validate")
	if validateTag != "" {
		fmt.Printf("A validate tag of the field %s : %s\n", field.Name, validateTag)
		tags := getRules(validateTag)
		rules, _ := v.createRules(tags)
		_ = tags
		_ = rules
	} else {
		fmt.Printf("The field %s hasn't a validate tag\n", field.Name)
	}
	return nil
}

func (v *CValidator) createRules(tags []string) (ret []RuleValidator, err error) {
	for _, t := range tags {
		s := strings.Split(t, ":")
		if rv, err := CreateRule(s[0], s[1]); err != nil {
			ret = append(ret, rv)
		} else {
			return nil, err
		}
	}
	return ret, err
}

func getRules(tag string) []string {
	return strings.Split(tag, "|")
}

func Validate(v interface{}) error {
	v0 := CValidator{in: v}
	return v0.Validate0()
}
