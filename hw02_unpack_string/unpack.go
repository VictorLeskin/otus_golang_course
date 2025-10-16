package hw02unpackstring

import (
	"errors"
	"strings"
	"unicode"
	"unicode/utf8"
)

var ErrInvalidString = errors.New("invalid string")

const (
	INIT = iota
	SYMBOL
	BACKSLASH
)

func Unpack(s string) (string, error) {
	// Place your code here.
	state := INIT
	prev := utf8.RuneError
	var sb strings.Builder

	for _, k := range s {
		switch state {
		case INIT:
			switch {
			case unicode.IsDigit(k):
				return "", ErrInvalidString
			case k == '\\':
				state = BACKSLASH
			default:
				prev = k
				state = SYMBOL
			}
		case SYMBOL:
			switch {
			case unicode.IsDigit(k):
				n := int(k - '0')
				for i := 0; i < n; i++ {
					sb.WriteRune(prev)
				}
				prev = utf8.RuneError
				state = INIT
			case k == '\\':
				sb.WriteRune(prev)
				state = BACKSLASH
			default:
				sb.WriteRune(prev)
				prev = k
			}
		case BACKSLASH:
			switch {
			case unicode.IsDigit(k) || k == '\\':
				prev = k
				state = SYMBOL
			default:
				return "", nil // only \\ or \0-9 are allowed
			}
		}
		// fmt.Printf("%d %d = %c = %s\n", state, prev, k, sb.String())
	}
	if prev != utf8.RuneError {
		sb.WriteRune(prev)
	}
	return sb.String(), nil
}
