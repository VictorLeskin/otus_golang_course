package hw09structvalidator

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type UserRole string

// Test the function on different structures and other types.
type (
	User struct {
		ID     string          `json:"id" validate:"len:36"`
		Name   string          `validate:"len:48"`
		Age    int             `validate:"min:18|max:50"`
		Email  string          `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
		Role   UserRole        `validate:"in:admin,stuff"`
		Phones []string        `validate:"len:11"`
		meta   json.RawMessage //nolint:unused
	}

	App struct {
		Version string `validate:"len:5"`
	}

	Token struct {
		Header    []byte
		Payload   []byte
		Signature []byte
	}

	Response struct {
		Code int    `validate:"in:200,404,500"`
		Body string `json:"omitempty"`
	}
)

func Test_ValidationErrors_Error(t *testing.T) {
	v := ValidationErrors{}
	v = append(v, ValidationError{Field: "Age", Err: fmt.Errorf("qqqqq")})
	v = append(v, ValidationError{Field: "Mail", Err: fmt.Errorf("tttt")})

	res := v.Error()
	expected := "qqqqq\n" + "tttt\n"

	assert.Equal(t, expected, res)
}

func Test_CValidator_appendValidatingError(t *testing.T) {
	// simple member
	{
		type User = TUser[int]

		user := User{Age: 43}

		t0 := &CValidator{ //  CValidator with only neccessary fields
			rv: reflect.ValueOf(user),
			rt: reflect.TypeOf(user),
		}

		t0.appendValidatingError("min", "Age", -1)

		assert.Equal(t, 1, len(t0.vErrors))
		expected := "Validating error of member 'Age' of struct 'TUser[int]' by rule 'min'"
		assert.Equal(t, expected, t0.vErrors[0].Err.Error())
	}

	// slice
	{
		type User = TUser[int]

		user := User{Age: 43}

		t0 := &CValidator{ //  CValidator with only neccessary fields
			rv: reflect.ValueOf(user),
			rt: reflect.TypeOf(user),
		}

		t0.appendValidatingError("min", "Age", 4)

		assert.Equal(t, 1, len(t0.vErrors))
		expected := "Validating error of member 'Age[4]' of struct 'TUser[int]' by rule 'min'"
		assert.Equal(t, expected, t0.vErrors[0].Err.Error())
	}
}

func Test_CValidator_getRules(t *testing.T) {
	t0 := &CValidator{}

	assert.Equal(t, []string{"min:4", "len:33"}, t0.getRules("min:4|len:33"))
}

func Test_CValidator_createRules(t *testing.T) {
	t0 := &CValidator{}

	// get two rules
	{
		res, err := t0.createRules([]string{"min:4", "len:33"})

		assert.Nil(t, err)
		assert.Equal(t, 2, len(res))
		v0, b0 := res[0].(*MinValidator)
		v1, b1 := res[1].(*LenValidator)
		assert.True(t, b0)
		assert.True(t, b1)
		assert.Equal(t, 4, int(v0.limit))
		assert.Equal(t, 33, int(v1.limit))
	}

	// error
	{
		res, err := t0.createRules([]string{"min:4", "Len:33"})

		assert.NotNil(t, err)
		assert.Nil(t, res)
		assert.Equal(t, "Wrong rule 'Len'", err.Error())
	}
}

func Test_CValidator_ValidateStruct(t *testing.T) {
	type S0 struct {
		ID  string `validate:"len:9"`
		Age int    `validate:"min:18|max:50"`
	}

	// successful validation of a whole struct
	{
		user := S0{
			ID:  "X138-A234",
			Age: 43,
		}

		rt := reflect.TypeOf(user)
		rv := reflect.ValueOf(user)

		v := &CValidator{ //  CValidator with only neccessary fields
			rv: rv,
			rt: rt,
		}

		err := v.validateStruct()
		assert.Nil(t, err)
		assert.Equal(t, 0, len(v.vErrors))
	}

	{
		// failed validation of a field
		user := S0{
			ID:  "X138-A234 addition not allowed",
			Age: 43,
		}

		rt := reflect.TypeOf(user)
		rv := reflect.ValueOf(user)

		v := &CValidator{ //  CValidator with only neccessary fields
			rv: rv,
			rt: rt,
		}

		err := v.validateStruct()
		assert.Nil(t, err)
		assert.Equal(t, 1, len(v.vErrors))
		expected := "Validating error of member 'ID' of struct 'S0' by rule 'len'"
		assert.Equal(t, expected, v.vErrors[0].Err.Error())
	}

	// not a struct
	{
		user := 999

		rt := reflect.TypeOf(user)
		rv := reflect.ValueOf(user)

		v := &CValidator{ //  CValidator with only neccessary fields
			rv: rv,
			rt: rt,
		}

		err := v.validateStruct()
		assert.NotNil(t, err)
		expected := "argument is not a struct"
		assert.Equal(t, expected, err.Error())
	}

}

func Test_CValidator_validateStructField(t *testing.T) {

	type S0 struct {
		ID           string  `validate:"len:9"`
		Age          int     `validate:"min:18|max:50"`
		meta         string  // no validation
		Mail         string  `validate:"Len:9"`
		AverageScore float64 `validate:"len:9"`
	}
	user := S0{
		ID:   "X138-A234 extra symbols",
		Age:  43,
		meta: "don't do it!",
	}

	rt := reflect.TypeOf(user)
	rv := reflect.ValueOf(user)

	v := &CValidator{ //  CValidator with only neccessary fields
		rv: rv,
		rt: rt,
	}

	// Age: ok
	{ // chek age
		typeField := v.rt.Field(1)  // type info
		valueField := v.rv.Field(1) // value info

		err := v.validateStructField(typeField, valueField)
		assert.Nil(t, err)
		assert.Equal(t, 0, len(v.vErrors))
	}

	// meta: ok - no checking
	{ // chek age
		typeField := v.rt.Field(2)  // type info
		valueField := v.rv.Field(2) // value info

		err := v.validateStructField(typeField, valueField)
		assert.Nil(t, err)
		assert.Equal(t, 0, len(v.vErrors))
	}

	// ID: no errors, but wrong input
	{ // chek age
		typeField := v.rt.Field(0)  // type info
		valueField := v.rv.Field(0) // value info

		err := v.validateStructField(typeField, valueField)
		assert.Nil(t, err)
		assert.Equal(t, 1, len(v.vErrors))
	}

	// Mail: no such validator
	{ // chek age
		typeField := v.rt.Field(3)  // type info
		valueField := v.rv.Field(3) // value info

		err := v.validateStructField(typeField, valueField)
		assert.NotNil(t, err)
	}

	// averageScore: right validator applied to wrong type
	{ // chek age
		typeField := v.rt.Field(4)  // type info
		valueField := v.rv.Field(4) // value info

		err := v.validateStructField(typeField, valueField)
		assert.NotNil(t, err)
	}

}

func TestValidateInt(t *testing.T) {
	// try to validate not a struct
	i := 42
	var in interface{}
	in = i
	ret := Validate(in)

	require.NotNil(t, ret)
}

func TestValidateStruct(t *testing.T) {
	// try to validate a struct
	user := User{
		ID:     "X138-A234",
		Name:   "Vic",
		Age:    43,
		Email:  "vic@in.the.middle.of.nowhere",
		Role:   "student",
		Phones: []string{"211-99-22", "303-42-77"},
		meta:   []byte(""),
	}

	var in interface{}
	in = user
	ret := Validate(in)

	require.Nil(t, ret)
}

func TestValidate(t *testing.T) {
	tests := []struct {
		in          interface{}
		expectedErr error
	}{
		{
			// Place your code here.
		},
		// ...
		// Place your code here.
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			tt := tt
			t.Parallel()

			// Place your code here.
			_ = tt
		})
	}
}
