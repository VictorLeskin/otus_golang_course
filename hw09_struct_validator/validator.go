package hw09structvalidator

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

type ValidationError struct {
	Field string
	Err   error
}

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() (ret string) {
	for _, s := range v {
		ret = ret + s.Err.Error() + "\n"
	}
	return ret
}

var ErrArgumentNotStructure = fmt.Errorf("argument is not a struct")

var (
	ErrExecution  = errors.New("execution error")
	ErrValidation = errors.New("validation failed")
)

type CValidator struct {
	in interface{}
	rv reflect.Value // intial struct value and type
	rt reflect.Type

	vErrors ValidationErrors
}

func (parent *CValidator) appendValidatingError(ruleName string, fieldName string, index int) {
	var ve error
	if index == -1 {
		ve = fmt.Errorf("Validating error of member '%s' of struct '%s' by rule '%s'", fieldName, parent.rt.Name(), ruleName)
	} else {
		ve = fmt.Errorf("Validating error of member '%s[%d]' of struct '%s' by rule '%s'", fieldName, index, parent.rt.Name(), ruleName)
	}

	parent.vErrors = append(parent.vErrors, ValidationError{Field: fieldName, Err: ve})
}

func (v *CValidator) getRules(tag string) []string {
	return strings.Split(tag, "|")
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

func (v *CValidator) validateStruct() error {
	if v.rt.Kind() == reflect.Struct {
		for i := 0; i < v.rt.NumField(); i++ {
			typeField := v.rt.Field(i)  // type info
			valueField := v.rv.Field(i) // value info
			if err := v.validateStructField(typeField, valueField); err != nil {
				return nil
			}
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
		// fmt.Printf("A validate tag of the field %s : %s\n", tf.Name, validateTag)
		tags := v.getRules(validateTag)
		rules, err := v.createRules(tags)
		if err != nil {
			return err
		}
		for _, r := range rules {
			err = ValidateValue(r, v, tf, vf)
			if err != nil {
				return err
			}
		}
	} else {
		//fmt.Printf("The field %s hasn't a validate tag\n", tf.Name)
	}
	return nil
}

func Validate(i interface{}) error {
	v := CValidator{in: i}

	v.rv = reflect.ValueOf(v.in)
	v.rt = v.rv.Type()

	if err := v.validateStruct(); err != nil {
		return fmt.Errorf("%w: %v", ErrExecution, ErrExecution)
	}

	if 0 == len(v.vErrors) {
		return nil
	}

	return fmt.Errorf("%w: %v", ErrValidation, v.vErrors.Error())
}
