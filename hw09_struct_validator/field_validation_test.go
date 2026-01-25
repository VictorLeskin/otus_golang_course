package hw09structvalidator

import (
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type LimitValidatableNumericalTypes interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 |
		~float32 | ~float64
}

type TUser[T LimitValidatableNumericalTypes] struct {
	Age         T
	Scores      []T
	Unsupported interface{}
	Name        string
}

func Test_RuleValidator_ValidateValue(t *testing.T) {
	t0, _ := createRuleMin("42")

	type User = TUser[int]

	// successful validating of single int
	{
		user := User{Age: 43}
		rt := reflect.TypeOf(user)
		rv := reflect.ValueOf(user)
		field, _ := rt.FieldByName("Age")

		validator := &CValidator{ //  CValidator with only neccessary fields
			rv: rv,
			rt: rt,
		}

		err1 := ValidateValue(t0, validator, field, rv.FieldByName("Age"))
		assert.NoError(t, err1)
		assert.Equal(t, 0, len(validator.vErrors))
	}

	// failed validated of single int
	{
		user := User{Age: 41}
		rt := reflect.TypeOf(user)
		rv := reflect.ValueOf(user)
		field, _ := rt.FieldByName("Age")

		validator := &CValidator{ //  CValidator with only neccessary fields
			rv: rv,
			rt: rt,
		}

		err1 := ValidateValue(t0, validator, field, rv.FieldByName("Age"))
		assert.NoError(t, err1)
		assert.Equal(t, 1, len(validator.vErrors))
		expectedText := "Validating error of member 'Age' of struct 'TUser[int]' by rule 'min'"
		assert.Equal(t, expectedText, validator.vErrors[0].Err.Error())
	}

	// successful validating of slice int
	{
		user := User{Scores: []int{43, 45}}
		rt := reflect.TypeOf(user)
		rv := reflect.ValueOf(user)
		field, _ := rt.FieldByName("Scores")

		validator := &CValidator{ //  CValidator with only neccessary fields
			rv: rv,
			rt: rt,
		}

		err1 := ValidateValue(t0, validator, field, rv.FieldByName("Scores"))
		assert.NoError(t, err1)
		assert.Equal(t, 0, len(validator.vErrors))
	}

	// failed validated of slice int
	{
		user := User{Scores: []int{41, 40}}
		rt := reflect.TypeOf(user)
		rv := reflect.ValueOf(user)
		field, _ := rt.FieldByName("Scores")

		validator := &CValidator{ //  CValidator with only neccessary fields
			rv: rv,
			rt: rt,
		}

		err1 := ValidateValue(t0, validator, field, rv.FieldByName("Scores"))
		assert.NoError(t, err1)
		assert.Equal(t, 2, len(validator.vErrors))
		expectedText0 := "Validating error of member 'Scores[0]' of struct 'TUser[int]' by rule 'min'"
		expectedText1 := "Validating error of member 'Scores[1]' of struct 'TUser[int]' by rule 'min'"
		assert.Equal(t, expectedText0, validator.vErrors[0].Err.Error())
		assert.Equal(t, expectedText1, validator.vErrors[1].Err.Error())
	}

	// wrong type
	{
		user := User{Age: 43}
		rt := reflect.TypeOf(user)
		rv := reflect.ValueOf(user)
		field, _ := rt.FieldByName("Unsupported")

		validator := &CValidator{ //  CValidator with only neccessary fields
			rv: rv,
			rt: rt,
		}

		err1 := ValidateValue(t0, validator, field, rv.FieldByName("Unsupported"))
		assert.Error(t, err1)
		expectedText := "Non unsupported type 'interface' by rule 'min'"
		assert.Equal(t, expectedText, err1.Error())
	}
}

/************ LenValidator ************/
func Test_createRuleLen(t *testing.T) {
	{
		res, err := createRuleLen("44")
		require.NoError(t, err)
		require.NotNil(t, res)

		res1 := res.(*LenValidator)
		assert.Equal(t, 44, res1.limit)
	}

	{
		_, err := createRuleLen("A")
		expectedText := "An invalid value in the rule 'len'"
		assert.Equal(t, true, strings.Contains(err.Error(), expectedText))
	}
}

func Test_LenValidator_ValidateValue0(t *testing.T) {
	t0, _ := createRuleLen("6")
	type User = TUser[int]

	// successful validating
	{
		user := User{Name: "012345"}

		rt := reflect.TypeOf(user)
		rv := reflect.ValueOf(user.Name)

		validator := &CValidator{ //  CValidator with only neccessary fields
			rv: rv,
			rt: rt,
		}

		err1 := t0.ValidateValue0(validator, "Name", rv.Type().Kind(), rv, -1)
		assert.NoError(t, err1)
		assert.Equal(t, 0, len(validator.vErrors))
	}

	// unsuccessful validating
	{
		user := User{Name: "very long name"}

		rt := reflect.TypeOf(user)
		rv := reflect.ValueOf(user.Name)

		validator := &CValidator{ //  CValidator with only neccessary fields
			rv: rv,
			rt: rt,
		}

		err1 := t0.ValidateValue0(validator, "Name", rv.Type().Kind(), rv, -1)
		assert.NoError(t, err1)
		assert.Equal(t, 1, len(validator.vErrors))
		expectedText := fmt.Sprintf("Validating error of member 'Name' of struct '%s' by rule 'len'", rt.Name())
		assert.Equal(t, expectedText, validator.vErrors[0].Err.Error())
	}

	// unsupported type
	{
		user := User{Name: "012345"}

		rt := reflect.TypeOf(user)
		rv := reflect.ValueOf(user)

		validator := &CValidator{ //  CValidator with only neccessary fields
			rv: rv,
			rt: rt,
		}

		err1 := t0.ValidateValue0(validator, "Name", rv.Type().Kind(), rv, -1)
		assert.Error(t, err1)
		expectedText := "Non unsupported type 'struct' by rule 'len'"
		assert.Equal(t, expectedText, err1.Error())
	}
}

/************ MinValidator ************/
func Test_createRuleMin(t *testing.T) {
	{
		res, err := createRuleMin("44")
		require.NoError(t, err)
		require.NotNil(t, res)

		res1 := res.(*MinValidator)
		assert.Equal(t, int64(44), res1.limit)
	}

	{
		_, err := createRuleMin("A")
		expectedText := "An invalid value in the rule 'min'"
		assert.Equal(t, true, strings.Contains(err.Error(), expectedText))
	}
}

func TTest_MinValidator_ValidateValue0[T LimitValidatableNumericalTypes](t *testing.T) {

	t0, _ := createRuleMin("42")
	type User = TUser[T]

	// successful validating of single int
	{
		user := User{Age: 43}

		rt := reflect.TypeOf(user)
		rv := reflect.ValueOf(user.Age)

		validator := &CValidator{ //  CValidator with only neccessary fields
			rv: rv,
			rt: rt,
		}

		err1 := t0.ValidateValue0(validator, "Age", rv.Type().Kind(), rv, -1)
		assert.NoError(t, err1)
		assert.Equal(t, 0, len(validator.vErrors))
	}

	// unsuccessful validating of single int
	{
		user := User{Age: 41}

		rt := reflect.TypeOf(user)
		rv := reflect.ValueOf(user.Age)

		validator := &CValidator{ //  CValidator with only neccessary fields
			rv: rv,
			rt: rt,
		}

		err1 := t0.ValidateValue0(validator, "Age", rv.Type().Kind(), rv, -1)
		assert.NoError(t, err1)
		assert.Equal(t, 1, len(validator.vErrors))
		expectedText := fmt.Sprintf("Validating error of member 'Age' of struct '%s' by rule 'min'", rt.Name())
		assert.Equal(t, expectedText, validator.vErrors[0].Err.Error())
	}

	// unsupported type
	{
		user := User{Age: 41}

		rt := reflect.TypeOf(user)
		rv := reflect.ValueOf(user)

		validator := &CValidator{ //  CValidator with only neccessary fields
			rv: rv,
			rt: rt,
		}

		err1 := t0.ValidateValue0(validator, "Age", rv.Type().Kind(), rv, -1)
		assert.Error(t, err1)
		expectedText := "Non unsupported type 'struct' by rule 'min'"
		assert.Equal(t, expectedText, err1.Error())
	}
}

func Test_MinValidator_ValidateValue0_Int(t *testing.T) {
	TTest_MinValidator_ValidateValue0[int](t)
	TTest_MinValidator_ValidateValue0[int8](t)
	TTest_MinValidator_ValidateValue0[int16](t)
	TTest_MinValidator_ValidateValue0[int32](t)
	TTest_MinValidator_ValidateValue0[int64](t)
	TTest_MinValidator_ValidateValue0[uint](t)
	TTest_MinValidator_ValidateValue0[uint8](t)
	TTest_MinValidator_ValidateValue0[uint16](t)
	TTest_MinValidator_ValidateValue0[uint32](t)
	TTest_MinValidator_ValidateValue0[uint64](t)
	TTest_MinValidator_ValidateValue0[float32](t)
	TTest_MinValidator_ValidateValue0[float64](t)
}

/************ MaxValidator ************/
func Test_createRuleMax(t *testing.T) {
	{
		res, err := createRuleMax("44")
		require.NoError(t, err)
		require.NotNil(t, res)

		res1 := res.(*MaxValidator)
		assert.Equal(t, int64(44), res1.limit)
	}

	{
		_, err := createRuleMax("A")
		expectedText := "An invalid value in the rule 'max'"
		assert.Equal(t, true, strings.Contains(err.Error(), expectedText))
	}
}

func TTest_MaxValidator_ValidateValue0[T LimitValidatableNumericalTypes](t *testing.T) {

	t0, _ := createRuleMax("42")
	type User = TUser[T]

	// successful validating of single int
	{
		user := User{Age: 41}

		rt := reflect.TypeOf(user)
		rv := reflect.ValueOf(user.Age)

		validator := &CValidator{ //  CValidator with only neccessary fields
			rv: rv,
			rt: rt,
		}

		err1 := t0.ValidateValue0(validator, "Age", rv.Type().Kind(), rv, -1)
		assert.NoError(t, err1)
		assert.Equal(t, 0, len(validator.vErrors))
	}

	// unsuccessful validating of single int
	{
		user := User{Age: 43}

		rt := reflect.TypeOf(user)
		rv := reflect.ValueOf(user.Age)

		validator := &CValidator{ //  CValidator with only neccessary fields
			rv: rv,
			rt: rt,
		}

		err1 := t0.ValidateValue0(validator, "Age", rv.Type().Kind(), rv, -1)
		assert.NoError(t, err1)
		assert.Equal(t, 1, len(validator.vErrors))
		expectedText := fmt.Sprintf("Validating error of member 'Age' of struct '%s' by rule 'max'", rt.Name())
		assert.Equal(t, expectedText, validator.vErrors[0].Err.Error())
	}

	// unsupported type
	{
		user := User{Age: 41}

		rt := reflect.TypeOf(user)
		rv := reflect.ValueOf(user)

		validator := &CValidator{ //  CValidator with only neccessary fields
			rv: rv,
			rt: rt,
		}

		err1 := t0.ValidateValue0(validator, "Age", rv.Type().Kind(), rv, -1)
		assert.Error(t, err1)
		expectedText := "Non unsupported type 'struct' by rule 'max'"
		assert.Equal(t, expectedText, err1.Error())
	}
}

func Test_MaxValidator_ValidateValue0_Int(t *testing.T) {
	TTest_MaxValidator_ValidateValue0[int](t)
	TTest_MaxValidator_ValidateValue0[int8](t)
	TTest_MaxValidator_ValidateValue0[int16](t)
	TTest_MaxValidator_ValidateValue0[int32](t)
	TTest_MaxValidator_ValidateValue0[int64](t)
	TTest_MaxValidator_ValidateValue0[uint](t)
	TTest_MaxValidator_ValidateValue0[uint8](t)
	TTest_MaxValidator_ValidateValue0[uint16](t)
	TTest_MaxValidator_ValidateValue0[uint32](t)
	TTest_MaxValidator_ValidateValue0[uint64](t)
	TTest_MaxValidator_ValidateValue0[float32](t)
	TTest_MaxValidator_ValidateValue0[float64](t)
}

/************ RegexpValidator ************/
func Test_createRuleRegexp(t *testing.T) {
	{
		res, err := createRuleRegexp("^[AB]+$")
		require.NoError(t, err)
		require.NotNil(t, res)

		res1 := res.(*RegexpValidator)
		assert.Equal(t, "^[AB]+$", res1.s)
		assert.NotNil(t, res1.re)
	}

	{
		_, err := createRuleRegexp("[")
		expectedText := "An invalid value in the rule 'regexp'"
		assert.Equal(t, true, strings.Contains(err.Error(), expectedText))
	}
}

func Test_RegexpValidator_ValidateValue0(t *testing.T) {
	t0, _ := createRuleRegexp("^[AB]+$")
	type User = TUser[int]

	// successful validating
	{
		user := User{Name: "ABBA"}

		rt := reflect.TypeOf(user)
		rv := reflect.ValueOf(user.Name)

		validator := &CValidator{ //  CValidator with only neccessary fields
			rv: rv,
			rt: rt,
		}

		err1 := t0.ValidateValue0(validator, "Name", rv.Type().Kind(), rv, -1)
		assert.NoError(t, err1)
		assert.Equal(t, 0, len(validator.vErrors))
	}

	// unsuccessful validating
	{
		user := User{Name: "very long name"}

		rt := reflect.TypeOf(user)
		rv := reflect.ValueOf(user.Name)

		validator := &CValidator{ //  CValidator with only neccessary fields
			rv: rv,
			rt: rt,
		}

		err1 := t0.ValidateValue0(validator, "Name", rv.Type().Kind(), rv, -1)
		assert.NoError(t, err1)
		assert.Equal(t, 1, len(validator.vErrors))
		expectedText := fmt.Sprintf("Validating error of member 'Name' of struct '%s' by rule 'regexp'", rt.Name())
		assert.Equal(t, expectedText, validator.vErrors[0].Err.Error())
	}

	// unsupported type
	{
		user := User{Name: "ABCD"}

		rt := reflect.TypeOf(user)
		rv := reflect.ValueOf(user)

		validator := &CValidator{ //  CValidator with only neccessary fields
			rv: rv,
			rt: rt,
		}

		err1 := t0.ValidateValue0(validator, "Name", rv.Type().Kind(), rv, -1)
		assert.Error(t, err1)
		expectedText := "Non unsupported type 'struct' by rule 'regexp'"
		assert.Equal(t, expectedText, err1.Error())
	}
}

/************ InValidator ************/
func Test_createRuleIn(t *testing.T) {
	{
		res, err := createRuleIn("123,456")
		require.NoError(t, err)
		require.NotNil(t, res)

		res1 := res.(*InValidator)
		assert.Equal(t, []string{"123", "456"}, res1.enabled)
		assert.Equal(t, []int{123, 456}, res1.enabledInt)
		assert.Equal(t, []float64{123.0, 456.0}, res1.enabledFlt)
		require.NoError(t, res1.errEnabledInt)
		require.NoError(t, res1.errEnabledFlt)
	}

	{
		res, err := createRuleIn("123,A")
		require.NoError(t, err)
		require.NotNil(t, res)

		res1 := res.(*InValidator)
		assert.Equal(t, []string{"123", "A"}, res1.enabled)
		assert.Nil(t, res1.enabledInt)
		assert.Nil(t, res1.enabledFlt)
		require.Error(t, res1.errEnabledInt)
		require.Error(t, res1.errEnabledFlt)
	}
}

func TTest_InValidator_ValidateValue0[T LimitValidatableNumericalTypes](t *testing.T) {

	t0, _ := createRuleIn("42,43")
	type User = TUser[T]

	// successful validating of single int
	{
		user := User{Age: 43}

		rt := reflect.TypeOf(user)
		rv := reflect.ValueOf(user.Age)

		validator := &CValidator{ //  CValidator with only neccessary fields
			rv: rv,
			rt: rt,
		}

		err1 := t0.ValidateValue0(validator, "Age", rv.Type().Kind(), rv, -1)
		assert.NoError(t, err1)
		assert.Equal(t, 0, len(validator.vErrors))
	}

	// unsuccessful validating of single int
	{
		user := User{Age: 41}

		rt := reflect.TypeOf(user)
		rv := reflect.ValueOf(user.Age)

		validator := &CValidator{ //  CValidator with only neccessary fields
			rv: rv,
			rt: rt,
		}

		err1 := t0.ValidateValue0(validator, "Age", rv.Type().Kind(), rv, -1)
		assert.NoError(t, err1)
		assert.Equal(t, 1, len(validator.vErrors))
		expectedText := fmt.Sprintf("Validating error of member 'Age' of struct '%s' by rule 'in'", rt.Name())
		assert.Equal(t, expectedText, validator.vErrors[0].Err.Error())
	}

	// unsupported type
	{
		user := User{Age: 41}

		rt := reflect.TypeOf(user)
		rv := reflect.ValueOf(user)

		validator := &CValidator{ //  CValidator with only neccessary fields
			rv: rv,
			rt: rt,
		}

		err1 := t0.ValidateValue0(validator, "Age", rv.Type().Kind(), rv, -1)
		assert.Error(t, err1)
		expectedText := "Non unsupported type 'struct' by rule 'in'"
		assert.Equal(t, expectedText, err1.Error())
	}
}

func Test_InValidator_ValidateValue0_Str(t *testing.T) {

	t0, _ := createRuleIn("A,B")
	type User = TUser[int]

	// successful validating of single int
	{
		user := User{Name: "A"}

		rt := reflect.TypeOf(user)
		rv := reflect.ValueOf(user.Name)

		validator := &CValidator{ //  CValidator with only neccessary fields
			rv: rv,
			rt: rt,
		}

		err1 := t0.ValidateValue0(validator, "Name", rv.Type().Kind(), rv, -1)
		assert.NoError(t, err1)
		assert.Equal(t, 0, len(validator.vErrors))
	}

	// unsuccessful validating of single int
	{
		user := User{Name: "C"}

		rt := reflect.TypeOf(user)
		rv := reflect.ValueOf(user.Name)

		validator := &CValidator{ //  CValidator with only neccessary fields
			rv: rv,
			rt: rt,
		}

		err1 := t0.ValidateValue0(validator, "Name", rv.Type().Kind(), rv, -1)
		assert.NoError(t, err1)
		assert.Equal(t, 1, len(validator.vErrors))
		expectedText := fmt.Sprintf("Validating error of member 'Name' of struct '%s' by rule 'in'", rt.Name())
		assert.Equal(t, expectedText, validator.vErrors[0].Err.Error())
	}

	// unsupported type
	{
		user := User{Name: "A"}

		rt := reflect.TypeOf(user)
		rv := reflect.ValueOf(user)

		validator := &CValidator{ //  CValidator with only neccessary fields
			rv: rv,
			rt: rt,
		}

		err1 := t0.ValidateValue0(validator, "Name", rv.Type().Kind(), rv, -1)
		assert.Error(t, err1)
		expectedText := "Non unsupported type 'struct' by rule 'in'"
		assert.Equal(t, expectedText, err1.Error())
	}
}

func Test_InValidator_ValidateValue0_Int(t *testing.T) {
	TTest_InValidator_ValidateValue0[int](t)
	TTest_InValidator_ValidateValue0[int8](t)
	TTest_InValidator_ValidateValue0[int16](t)
	TTest_InValidator_ValidateValue0[int32](t)
	TTest_InValidator_ValidateValue0[int64](t)
	TTest_InValidator_ValidateValue0[uint](t)
	TTest_InValidator_ValidateValue0[uint8](t)
	TTest_InValidator_ValidateValue0[uint16](t)
	TTest_InValidator_ValidateValue0[uint32](t)
	TTest_InValidator_ValidateValue0[uint64](t)
	TTest_InValidator_ValidateValue0[float32](t)
	TTest_InValidator_ValidateValue0[float64](t)
}
