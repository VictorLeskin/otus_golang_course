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
	rv reflect.Value // intial struct value and type
	rt reflect.Type

	vErrors []validatingError
}

func (parent *CValidator) appendValidatingError(ruleName string, fieldName string, index int) {
	var ve error
	if index == -1 {
		ve = fmt.Errorf("Validating error of member '%s' of struct '%s' by rule '%s'", fieldName, parent.rt.Name(), ruleName)
	} else {
		ve = fmt.Errorf("Validating error of member '%s[%d]' of struct '%s' by rule '%s'", fieldName, index, parent.rt.Name(), ruleName)
	}

	parent.vErrors = append(parent.vErrors, ve)
}

func (v *CValidator) Validate0() error {
	v.rv = reflect.ValueOf(v.in)
	v.rt = v.rv.Type()

	return v.validateStruct()
}

func (v *CValidator) validateStruct() error {
	if v.rt.Kind() == reflect.Struct {
		for i := 0; i < v.rt.NumField(); i++ {
			typeField := v.rt.Field(i)  // type info
			valueField := v.rv.Field(i) // value info
			v.validateStructField(typeField, valueField)
		}
	} else {
		return ErrArgumentNotStructure
	}

	return nil
}

func (v *CValidator) validateStructField(tf reflect.StructField, vf reflect.Value) error {
	// get validate tag.
	validateTag := tf.Tag.Get("validate")
	if validateTag != "" {
		fmt.Printf("A validate tag of the field %s : %s\n", tf.Name, validateTag)
		tags := getRules(validateTag)
		rules, err := v.createRules(tags)
		if err != nil {
			return err
		}
		for _, r := range rules {
			err = r.ValidateValue(v, tf, vf)
			if err != nil {
				return err
			}
		}
	} else {
		fmt.Printf("The field %s hasn't a validate tag\n", tf.Name)
	}
	return nil
}

func (v *CValidator) createRules(tags []string) (ret []RuleValidator, err error) {
	for _, t := range tags {
		s := strings.Split(t, ":")
		if rv, err := CreateRule(s[0], s[1]); err == nil {
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
