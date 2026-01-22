package hw09structvalidator

import (
	"fmt"
	"reflect"
	"strconv"
)

type RuleValidator interface {
	ValidateValue(tf reflect.StructField, vf reflect.Value) error
}

type LenValidator struct {
	RuleValidator
	limit int
}

func (v *LenValidator) ValidateValue(tp reflect.StructField, vf reflect.Value) error {
	// get type and value:
	typ := tp.Type // тип: int

	fmt.Printf("Type: %v\n", typ) // Type: int
	fmt.Printf("Value: %v\n", vf) // Value: 42
	return nil
}

type MinValidator struct {
	RuleValidator
	limit int
}

func (v *MinValidator) ValidateValue(tf reflect.StructField, vf reflect.Value) error {
	return nil
}

func CreateRule(name string, value string) (RuleValidator, error) {
	switch name {
	case "len":
		return createRuleLen(value)
	case "min":
		return createRuleMin(value)
	case "max":
		return createRuleMax(value)
	case "regexp":
		return createRuleRegexp(value)
	case "in":
		return createRuleIn(value)
	default:
		break
	}

	return nil, nil
}

func createRuleLen(value string) (RuleValidator, error) {
	limit, err := strconv.Atoi(value)
	if err != nil {
		return nil, fmt.Errorf("invalid value in rule  %v", err)
	}

	return &LenValidator{
		limit: limit}, nil
}

func createRuleMin(value string) (RuleValidator, error) {
	limit, err := strconv.Atoi(value)
	if err != nil {
		return nil, fmt.Errorf("invalid value in rule %v", err)
	}

	return &MinValidator{
		limit: limit}, nil
}

func createRuleMax(value string) (RuleValidator, error)    { _ = value; return nil, nil }
func createRuleRegexp(value string) (RuleValidator, error) { _ = value; return nil, nil }
func createRuleIn(value string) (RuleValidator, error)     { _ = value; return nil, nil }
