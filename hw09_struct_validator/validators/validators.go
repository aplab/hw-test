package validators

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"unicode/utf8"
)

var (
	ErrUnknownRule         = errors.New("unknown rule")
	ErrInvalidSyntax       = errors.New("invalid syntax")
	ErrInvalidRegexpSyntax = errors.New("invalid regexp syntax")
)

var (
	RuleIntMin       = regexp.MustCompile(`min:(\d+)`)
	RuleIntMax       = regexp.MustCompile(`max:(\d+)`)
	RuleIntIn        = regexp.MustCompile(`in:(.*)`)
	RuleStringLen    = regexp.MustCompile(`len:(\d+)`)
	RuleStringIn     = regexp.MustCompile(`in:(.*)`)
	RuleStringRegexp = regexp.MustCompile(`regexp:(.*)`)
	RulesSplit       = regexp.MustCompile(`(?:\\|)\|`)
)

func IntMin(value, limit int) error {
	if value < limit {
		return fmt.Errorf("value %v less than %v", value, limit)
	}
	return nil
}

func IntMax(value, limit int) error {
	if value > limit {
		return fmt.Errorf("value %v greater than %v", value, limit)
	}
	return nil
}

func IntIn(value int, in []int) error {
	for _, i := range in {
		if value == i {
			return nil
		}
	}
	return fmt.Errorf("value %v not in set %v", value, in)
}

func StringLen(value string, length int) error {
	if utf8.RuneCountInString(value) > length {
		return fmt.Errorf("value %v is longer than %v", value, length)
	}
	return nil
}

func StringIn(value string, in []string) error {
	for _, i := range in {
		if value == i {
			return nil
		}
	}
	return fmt.Errorf("value %v not in set %v", value, in)
}

func StringRegexp(value string, regexp *regexp.Regexp) error {
	if regexp.MatchString(value) {
		return nil
	}
	return fmt.Errorf("value %v is not match %v", value, regexp.String())
}

func SliceIntMin(value []int, limit int) error {
	for _, i := range value {
		if err := IntMin(i, limit); err != nil {
			return err
		}
	}
	return nil
}

func SliceIntMax(value []int, limit int) error {
	for _, i := range value {
		if err := IntMax(i, limit); err != nil {
			return err
		}
	}
	return nil
}

func SliceIntIn(value, in []int) error {
	for _, i := range value {
		if err := IntIn(i, in); err != nil {
			return err
		}
	}
	return nil
}

func SliceStringLen(value []string, length int) error {
	for _, i := range value {
		if err := StringLen(i, length); err != nil {
			return err
		}
	}
	return nil
}

func SliceStringIn(value, in []string) error {
	for _, i := range value {
		if err := StringIn(i, in); err != nil {
			return err
		}
	}
	return nil
}

func SliceStringRegexp(value []string, regexp *regexp.Regexp) error {
	for _, i := range value {
		if err := StringRegexp(i, regexp); err != nil {
			return err
		}
	}
	return nil
}

func IntGetRules(rule string) ([]func(int) error, error) {
	rules := RulesSplit.Split(rule, -1)
	functions := make([]func(int) error, len(rules))
	for i, rule := range rules {
		rules[i] = strings.ReplaceAll(rule, "\\|", "|")
		f, err := IntGetRule(rules[i])
		if err != nil {
			return functions, err
		}
		functions[i] = f
	}
	return functions, nil
}

func StringGetRules(rule string) ([]func(string) error, error) {
	rules := RulesSplit.Split(rule, -1)
	functions := make([]func(string) error, len(rules))
	for i, rule := range rules {
		rules[i] = strings.ReplaceAll(rule, "\\|", "|")
		f, err := StringGetRule(rules[i])
		if err != nil {
			return functions, err
		}
		functions[i] = f
	}
	return functions, nil
}

func StringGetRule(rule string) (f func(string) error, err error) {
	data := RuleStringLen.FindStringSubmatch(rule)
	if data != nil {
		l, err := strconv.Atoi(data[1])
		if err != nil {
			return nil, ErrInvalidSyntax
		}
		return func(s string) error {
			return StringLen(s, l)
		}, nil
	}
	data = RuleStringIn.FindStringSubmatch(rule)
	if data != nil {
		return func(s string) error {
			return StringIn(s, strings.Split(data[1], ","))
		}, nil
	}
	data = RuleStringRegexp.FindStringSubmatch(rule)
	if data != nil {
		reg, err := regexp.Compile(data[1])
		if err != nil {
			return nil, ErrInvalidRegexpSyntax
		}
		return func(s string) error {
			return StringRegexp(s, reg)
		}, nil
	}
	return nil, ErrUnknownRule
}

func IntGetRule(rule string) (f func(int) error, err error) {
	data := RuleIntMax.FindStringSubmatch(rule)
	if data != nil {
		limit, err := strconv.Atoi(data[1])
		if err != nil {
			return nil, ErrInvalidSyntax
		}
		return func(i int) error {
			return IntMax(i, limit)
		}, nil
	}
	data = RuleIntMin.FindStringSubmatch(rule)
	if data != nil {
		limit, err := strconv.Atoi(data[1])
		if err != nil {
			return nil, ErrInvalidSyntax
		}
		return func(i int) error {
			return IntMin(i, limit)
		}, nil
	}
	data = RuleIntIn.FindStringSubmatch(rule)
	if data != nil {
		variants := strings.Split(data[1], ",")
		in := make([]int, len(variants))
		for i, v := range variants {
			iv, err := strconv.Atoi(v)
			if err != nil {
				return nil, ErrInvalidSyntax
			}
			in[i] = iv
		}
		return func(i int) error {
			return IntIn(i, in)
		}, nil
	}
	return nil, ErrUnknownRule
}
