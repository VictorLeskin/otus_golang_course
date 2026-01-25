package hw09structvalidator

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

type RuleValidator interface {
	ValidateValue0(parent *CValidator, name string, kind reflect.Kind, rv reflect.Value, index int) error
	Name() string
}

func ValidateValue(v RuleValidator, parent *CValidator, tp reflect.StructField, rv reflect.Value) error {
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
	s  string        // initial string
	re regexp.Regexp // trasformed to regular expression
}

type InValidator struct {
	RuleValidator
	enabled []string

	enabledInt    []int
	errEnabledInt error

	enabledFlt    []float64
	errEnabledFlt error
}

func (v *LenValidator) ValidateValue0(parent *CValidator, name string, kind reflect.Kind, rv reflect.Value, index int) error {
	switch kind {
	case reflect.String:
		if len(rv.String()) != v.limit {
			parent.appendValidatingError(v.Name(), name, index)
		}

	default:
		return fmt.Errorf("Non unsupported type '%s' by rule '%s'", kind.String(), v.Name())
	}
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
		return fmt.Errorf("Non unsupported type '%s' by rule '%s'", kind.String(), v.Name())
	}
	return nil
}

func (v *MaxValidator) ValidateValue0(parent *CValidator, name string, kind reflect.Kind, rv reflect.Value, index int) error {
	switch kind {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64: //  int
		if rv.Int() > v.limit {
			parent.appendValidatingError(v.Name(), name, index)
		}

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64: // uint
		if rv.Uint() > uint64(v.limit) {
			parent.appendValidatingError(v.Name(), name, index)
		}

	case reflect.Float32, reflect.Float64: // floatt
		if rv.Float() > float64(v.limit) {
			parent.appendValidatingError(v.Name(), name, index)
		}
	default:
		return fmt.Errorf("Non unsupported type '%s' by rule '%s'", kind.String(), v.Name())
	}
	return nil
}

func (v *RegexpValidator) ValidateValue0(parent *CValidator, name string, kind reflect.Kind, rv reflect.Value, index int) error {
	switch kind {
	case reflect.String:
		if !v.re.MatchString(rv.String()) {
			parent.appendValidatingError(v.Name(), name, index)
		}

	default:
		return fmt.Errorf("Non unsupported type '%s' by rule '%s'", kind.String(), v.Name())
	}
	return nil
}

func (v *InValidator) toInts() (ret []int, err error) {
	for _, value := range v.enabled {
		i, err := strconv.Atoi(value)
		if err != nil {
			return nil, fmt.Errorf("invalid integer value %s for in value in rule 'in' %v", value, err)
		}
		ret = append(ret, i)
	}

	return ret, err
}

func (v *InValidator) toFloats() (ret []float64, err error) {
	for _, value := range v.enabled {
		f, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid integer value %s for in value in rule 'in' %v", value, err)
		}
		ret = append(ret, f)
	}

	return ret, err
}

func (v *InValidator) ValidateValue0(parent *CValidator, name string, kind reflect.Kind, rv reflect.Value, index int) error {
	switch kind {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64: //  int
		if v.errEnabledInt != nil {
			return v.errEnabledInt
		}
		for _, e := range v.enabledInt {
			if rv.Int() == int64(e) {
				return nil
			}
		}

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64: // uint
		if v.errEnabledInt != nil {
			return v.errEnabledInt
		}
		for _, e := range v.enabledInt {
			if rv.Uint() == uint64(e) {
				return nil
			}
		}

	case reflect.Float32, reflect.Float64: // floatt
		if v.errEnabledFlt != nil {
			return v.errEnabledFlt
		}
		for _, e := range v.enabledFlt {
			if rv.Float() == e {
				return nil
			}
		}
	case reflect.String:
		for _, e := range v.enabled {
			if rv.String() == e {
				return nil
			}
		}
	default:
		return fmt.Errorf("Non unsupported type '%s' by rule '%s'", kind.String(), v.Name())
	}
	parent.appendValidatingError(v.Name(), name, index)
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

	return nil, fmt.Errorf("Wrong rule '%s'", name)
}

func createRuleLen(value string) (RuleValidator, error) {
	limit, err := strconv.Atoi(value)
	if err != nil {
		return nil, fmt.Errorf("An invalid value in the rule 'len': %v", err)
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
	re, err := regexp.Compile(value) // возвращает ошибку
	if err != nil {
		return nil, fmt.Errorf("An invalid value in the rule 'regexp': %v", err)
	}

	return &RegexpValidator{
		s: value, re: *re}, nil
}

func createRuleIn(value string) (RuleValidator, error) {
	ret := &InValidator{
		enabled: strings.Split(value, ",")}
	ret.enabledInt, ret.errEnabledInt = ret.toInts()
	ret.enabledFlt, ret.errEnabledFlt = ret.toFloats()

	return ret, nil
}

func (v *LenValidator) Name() string    { return "len" }
func (v *MinValidator) Name() string    { return "min" }
func (v *MaxValidator) Name() string    { return "max" }
func (v *RegexpValidator) Name() string { return "regexp" }
func (v *InValidator) Name() string     { return "in" }
