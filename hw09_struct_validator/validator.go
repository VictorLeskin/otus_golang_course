package hw09structvalidator

import (
	"fmt"
	"reflect"
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

func processField(field reflect.StructField) error {
	// get validate tag.
	validateTag := field.Tag.Get("validate")
	if validateTag != "" {
		fmt.Printf("A validate tag of the field %s : %s\n", field.Name, validateTag)
	} else {
		fmt.Printf("The field %s hasn't a validate tag\n", field.Name)
	}
	return nil
}

func Validate(v interface{}) error {
	rv := reflect.ValueOf(v)
	rt := rv.Type()

	if rt.Kind() == reflect.Struct {
		for i := 0; i < rt.NumField(); i++ {
			processField(rt.Field(i))
		}
	} else {
		return ErrArgumentNotStructure
	}

	return nil
}
