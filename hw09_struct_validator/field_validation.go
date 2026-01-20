package hw09structvalidator

import (
	"fmt"
	"strconv"
)

type RuleValidator interface {
	ValidateValue(value interface{}) error
}

type LenValidator struct {
	RuleValidator
	limit int
}

func (v *LenValidator) ValidateValue(value interface{}) error {
	return nil
}

type MinValidator struct {
	RuleValidator
	limit int
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

func createRuleMax(value string) (RuleValidator, error)    { return nil, nil }
func createRuleRegexp(value string) (RuleValidator, error) { return nil, nil }
func createRuleIn(value string) (RuleValidator, error)     { return nil, nil }
