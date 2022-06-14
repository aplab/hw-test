package hw03frequencyanalysis

import (
	"regexp"
	"sort"
	"strings"
)

const limit int = 10

type part struct {
	Word  string
	Count uint
}

var re = regexp.MustCompile("[!-,.-/:-@\\[-`{-~]|\\s-+|-+\\s")

func Top10(s string) []string {
	freq := make(map[string]uint)
	for _, v := range strings.Fields(s) {
		freq[v]++
	}
	parts := make([]part, len(freq))
	i := 0
	for k, v := range freq {
		parts[i] = part{
			Word:  k,
			Count: v,
		}
		i++
	}
	sort.Slice(parts, func(i, j int) bool {
		if parts[i].Count == parts[j].Count {
			return parts[i].Word < parts[j].Word
		}
		return parts[i].Count > parts[j].Count
	})
	l := len(parts)
	if l > limit {
		l = limit
	}
	ret := make([]string, l)
	for i := range parts[:l] {
		ret[i] = parts[i].Word
	}
	return ret
}

func Top10s(s string) []string {
	s = re.ReplaceAllString(" "+s+" ", " ")
	freq := make(map[string]uint)
	for _, v := range strings.Fields(s) {
		freq[strings.ToLower(v)]++
	}
	parts := make([]part, len(freq))
	i := 0
	for k, v := range freq {
		parts[i] = part{
			Word:  k,
			Count: v,
		}
		i++
	}
	sort.Slice(parts, func(i, j int) bool {
		if parts[i].Count == parts[j].Count {
			return parts[i].Word < parts[j].Word
		}
		return parts[i].Count > parts[j].Count
	})
	l := len(parts)
	if l > limit {
		l = limit
	}
	ret := make([]string, l)
	for i := range parts[:l] {
		ret[i] = parts[i].Word
	}
	return ret
}
