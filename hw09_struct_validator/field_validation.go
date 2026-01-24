package hw09structvalidator

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

type validatingError error

type RuleValidator interface {
	ValidateValue(parent *CValidator, tp reflect.StructField, rv reflect.Value) error
	Name() string
}

type LenValidator struct {
	RuleValidator
	limit int
}

type MinValidator struct {
	RuleValidator
	limit int64
}

type MaxValidator struct {
	RuleValidator
	limit int64
}

type RegexpValidator struct {
	RuleValidator
	regex string
}

type InValidator struct {
	RuleValidator
	enabled []string
}

func (v *LenValidator) ValidateValue(parent *CValidator, tp reflect.StructField, rv reflect.Value) error {
	// get type and value:
	rt := tp.Type

	fmt.Printf("Len validator\n")
	fmt.Printf("Type: %v\n", rt)
	fmt.Printf("Value: %v\n", rv)
	return nil
}

func (v *MinValidator) ValidateValue0(parent *CValidator, name string, kind reflect.Kind, rv reflect.Value, index int) error {
	switch kind {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64: //  int
		if rv.Int() < v.limit {
			parent.appendValidatingError(v.Name(), name, index)
		}

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64: // uint
		if rv.Uint() < uint64(v.limit) {
			parent.appendValidatingError(v.Name(), name, index)
		}

	case reflect.Float32, reflect.Float64: // floatt
		if rv.Float() < float64(v.limit) {
			parent.appendValidatingError(v.Name(), name, index)
		}
	default:
		return fmt.Errorf("Non unsupported type: %s", kind.String())
	}
	return nil
}

func (v *MinValidator) ValidateValue(parent *CValidator, tp reflect.StructField, rv reflect.Value) error {
	// get type
	rt := tp.Type

	if rt.Kind() == reflect.Slice {
		elemType := rt.Elem()
		for i := 0; i < rv.Len(); i++ {
			if err := v.ValidateValue0(parent, tp.Name, elemType.Kind(), rv.Index(i), i); err != nil {
				return err
			}
		}
	} else {
		return v.ValidateValue0(parent, tp.Name, rt.Kind(), rv, -1)
	}
	return nil
}

func (v *MaxValidator) ValidateValue(parent *CValidator, tp reflect.StructField, rv reflect.Value) error {
	// get type and value:
	rt := tp.Type

	fmt.Printf("Max validator\n")
	fmt.Printf("Type: %v\n", rt)
	fmt.Printf("Value: %v\n", rv)
	return nil
}

func (v *RegexpValidator) ValidateValue(parent *CValidator, tp reflect.StructField, rv reflect.Value) error {
	// get type and value:
	rt := tp.Type

	fmt.Printf("Regexp validator\n")
	fmt.Printf("Type: %v\n", rt)
	fmt.Printf("Value: %v\n", rv)
	return nil
}

func (v *InValidator) ValidateValue(parent *CValidator, tp reflect.StructField, rv reflect.Value) error {
	// get type and value:
	rt := tp.Type

	fmt.Printf("In validator\n")
	fmt.Printf("Type: %v\n", rt)
	fmt.Printf("Value: %v\n", rv)
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
		return nil, fmt.Errorf("An invalid value in the rule 'min': %v", err)
	}

	return &MinValidator{
		limit: int64(limit)}, nil
}

func createRuleMax(value string) (RuleValidator, error) {
	limit, err := strconv.Atoi(value)
	if err != nil {
		return nil, fmt.Errorf("invalid value in rule %v", err)
	}

	return &MaxValidator{
		limit: int64(limit)}, nil
}
func createRuleRegexp(value string) (RuleValidator, error) {
	// regex.Compile()

	return &RegexpValidator{
		regex: value}, nil
}

func createRuleIn(value string) (RuleValidator, error) {
	return &InValidator{
		enabled: strings.Split(value, ",")}, nil
}

func (v *LenValidator) Name() string    { return "len" }
func (v *MinValidator) Name() string    { return "min" }
func (v *MaxValidator) Name() string    { return "max" }
func (v *RegexpValidator) Name() string { return "regexp" }
func (v *InValidator) Name() string     { return "in" }
