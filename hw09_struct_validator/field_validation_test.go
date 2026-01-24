package hw09structvalidator

import (
	"reflect"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_LenValidator_ctor(t *testing.T) {
	t0 := &LenValidator{limit: 8}

	var t1 RuleValidator
	t1 = t0

	_ = t1

}

func Test_MinValidator_createRuleMin(t *testing.T) {
	{
		res, err := createRuleMin("44")
		require.NoError(t, err)
		require.NotNil(t, res)

		res1 := res.(*MinValidator)
		assert.Equal(t, 44, res1.limit)
	}

	{
		_, err := createRuleMin("A")
		expectedText := "An invalid value in the rule 'min'"
		assert.Equal(t, true, strings.Contains(err.Error(), expectedText))
	}
}

type LimitValidatableNumericalTypes interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 |
		~float32 | ~float64
}

type TUser[T LimitValidatableNumericalTypes] struct {
	Age    T
	Scores []T
}

func Test_MinValidator_ValidateValue_Int(t *testing.T) {
	t0, _ := createRuleMin("42")

	type User = TUser[int]

	// successful validated of single int
	{
		user := User{Age: 43}
		rt := reflect.TypeOf(user)
		rv := reflect.ValueOf(user)
		field, _ := rt.FieldByName("Age")

		validator := &CValidator{ //  CValidator with only neccessary fields
			rv: rv,
			rt: rt,
		}

		err1 := t0.ValidateValue(validator, field, rv.FieldByName("Age"))
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

		err1 := t0.ValidateValue(validator, field, rv.FieldByName("Age"))
		assert.NoError(t, err1)
		assert.Equal(t, 1, len(validator.vErrors))
		expectedText := "Validating error of member 'Age' of struct 'TUser[int]' by rule 'min'"
		assert.Equal(t, expectedText, validator.vErrors[0].Error())
	}

	// successful validated of slice int
	{
		user := User{Scores: []int{43, 45}}
		rt := reflect.TypeOf(user)
		rv := reflect.ValueOf(user)
		field, _ := rt.FieldByName("Scores")

		validator := &CValidator{ //  CValidator with only neccessary fields
			rv: rv,
			rt: rt,
		}

		err1 := t0.ValidateValue(validator, field, rv.FieldByName("Scores"))
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

		err1 := t0.ValidateValue(validator, field, rv.FieldByName("Scores"))
		assert.NoError(t, err1)
		assert.Equal(t, 2, len(validator.vErrors))
		expectedText0 := "Validating error of member 'Scores[0]' of struct 'TUser[int]' by rule 'min'"
		expectedText1 := "Validating error of member 'Scores[1]' of struct 'TUser[int]' by rule 'min'"
		assert.Equal(t, expectedText0, validator.vErrors[0].Error())
		assert.Equal(t, expectedText1, validator.vErrors[1].Error())
	}

}

func TTest_MinValidator_ValidateValue[T LimitValidatableNumericalTypes](t *testing.T) {
	t0, _ := createRuleMin("42")

	type User = TUser[T]

	// successful validated of single int
	{
		user := User{Age: 43}
		rt := reflect.TypeOf(user)
		rv := reflect.ValueOf(user)
		field, _ := rt.FieldByName("Age")

		validator := &CValidator{ //  CValidator with only neccessary fields
			rv: rv,
			rt: rt,
		}

		err1 := t0.ValidateValue(validator, field, rv.FieldByName("Age"))
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

		err1 := t0.ValidateValue(validator, field, rv.FieldByName("Age"))
		assert.NoError(t, err1)
		assert.Equal(t, 1, len(validator.vErrors))
		expectedText := "Validating error of member 'Age' of struct 'TUser[int]' by rule 'min'"
		assert.Equal(t, expectedText, validator.vErrors[0].Error())
	}

	// successful validated of slice int
	{
		user := User{Scores: []T{43, 45}}
		rt := reflect.TypeOf(user)
		rv := reflect.ValueOf(user)
		field, _ := rt.FieldByName("Scores")

		validator := &CValidator{ //  CValidator with only neccessary fields
			rv: rv,
			rt: rt,
		}

		err1 := t0.ValidateValue(validator, field, rv.FieldByName("Scores"))
		assert.NoError(t, err1)
		assert.Equal(t, 0, len(validator.vErrors))
	}

	// failed validated of slice int
	{
		user := User{Scores: []T{41, 40}}
		rt := reflect.TypeOf(user)
		rv := reflect.ValueOf(user)
		field, _ := rt.FieldByName("Scores")

		validator := &CValidator{ //  CValidator with only neccessary fields
			rv: rv,
			rt: rt,
		}

		err1 := t0.ValidateValue(validator, field, rv.FieldByName("Scores"))
		assert.NoError(t, err1)
		assert.Equal(t, 2, len(validator.vErrors))
		expectedText0 := "Validating error of member 'Scores[0]' of struct 'TUser[int]' by rule 'min'"
		expectedText1 := "Validating error of member 'Scores[1]' of struct 'TUser[int]' by rule 'min'"
		assert.Equal(t, expectedText0, validator.vErrors[0].Error())
		assert.Equal(t, expectedText1, validator.vErrors[1].Error())
	}
}

func Test_MinValidator_ValidateValue_AllNumericTypes(t *testing.T) {
	TTest_MinValidator_ValidateValue[int](t)
	TTest_MinValidator_ValidateValue[int8](t)
	TTest_MinValidator_ValidateValue[int16](t)
	TTest_MinValidator_ValidateValue[int32](t)
	TTest_MinValidator_ValidateValue[int64](t)
	TTest_MinValidator_ValidateValue[uint](t)
	TTest_MinValidator_ValidateValue[uint8](t)
	TTest_MinValidator_ValidateValue[uint16](t)
	TTest_MinValidator_ValidateValue[uint32](t)
	TTest_MinValidator_ValidateValue[uint64](t)
	TTest_MinValidator_ValidateValue[float32](t)
	TTest_MinValidator_ValidateValue[float64](t)
}
