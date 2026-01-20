package hw09structvalidator

import (
	"fmt"
	"strconv"
	"strings"
)

type RuleValidator interface {
	ValidateValue(value interface{}) error
}

type LengthValidator struct {
	RuleValidator
	ruleType string
	limit    int
}

func NewLengthValidator(rule string) (RuleValidator, error) {
	parts := strings.Split(rule, ":")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid rule format: %s", rule)
	}

	ruleType := parts[0]
	limitStr := parts[1]

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		return nil, fmt.Errorf("invalid limit value in rule %s: %v", rule, err)
	}

	// Проверяем, что ruleType один из известных
	if ruleType != "len" {
		return nil, fmt.Errorf("wrong len rule type: %s", ruleType)
	}

	return &LengthValidator{
		ruleType: ruleType,
		limit:    limit}, nil
}

func CreateRule(name string, value string) RuleValidator {
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

	return nil
}

func createRuleLen(value string) RuleValidator    { return nil }
func createRuleMin(value string) RuleValidator    { return nil }
func createRuleMax(value string) RuleValidator    { return nil }
func createRuleRegexp(value string) RuleValidator { return nil }
func createRuleIn(value string) RuleValidator     { return nil }
