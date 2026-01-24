package hw09structvalidator

import (
	"encoding/json"
	"fmt"
	"testing"

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
