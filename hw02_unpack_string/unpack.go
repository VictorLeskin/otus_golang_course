package hw02unpackstring

import (
	"errors"
	"fmt"
	"strings"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

const (
	INIT = iota
	SYMBOL
)

func Unpack(s string) (string, error) {
	// Place your code here.
	var state int
	state = INIT
	var sb strings.Builder
	var prev rune

	for _, k := range s {
		switch state {

		case INIT:
			switch {
			case unicode.IsDigit(k):
				return "", ErrInvalidString
			default:
				prev = rune(k)
				state = SYMBOL
			}
		case SYMBOL:
			switch {
			case unicode.IsDigit(k):
				n := int(k)
				for i := 0; i < n; i++ {
					sb.WriteRune(prev)
				}

				return "", ErrInvalidString
			default:
				prev = rune(k)
				state = SYMBOL
			}

		}

		fmt.Printf("%d %d = %c = %s\n", state, prev, k, sb.String())

	}
	return sb.String(), nil
}
