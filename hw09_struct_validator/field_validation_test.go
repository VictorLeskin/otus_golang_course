package hw09structvalidator

import (
	"testing"
)

func Test_LenValidator_ctor(t *testing.T) {
	t0 := &LenValidator{limit: 8}

	var t1 RuleValidator
	t1 = t0

	_ = t1

}
