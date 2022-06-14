package hw02unpackstring

import (
	"errors"
	"strings"
	"unicode"
	"unicode/utf8"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(s string) (string, error) {
	if !utf8.ValidString(s) {
		return s, ErrInvalidString
	}
	if utf8.RuneCountInString(s) == 0 {
		return s, nil
	}
	b := strings.Builder{}
	var escaped bool
	var prev rune
	var prevEscaped bool
	var str string
	var prevStr string
	for i, r := range s {
		switch {
		case isSlash(r):
			if escaped {
				escaped = false
				str = string(r)
			} else {
				escaped = true
			}
		case unicode.IsDigit(r):
			if escaped {
				escaped = false
				str = string(r)
			} else {
				if (unicode.IsDigit(prev) && !prevEscaped) || i == 0 {
					return s, ErrInvalidString
				}
				m := int(r - '1')
				if m < 1 {
					prevStr = ""
				} else {
					str = strings.Repeat(string(prev), m)
				}
			}
		default:
			str = string(r)
		}
		prevEscaped = isSlash(prev) && !isSlash(r)
		prev = r
		b.WriteString(prevStr)
		prevStr = str
		str = ""
	}
	b.WriteString(prevStr)
	return b.String(), nil
}

func isSlash(r rune) bool {
	return r == '\\'
}
