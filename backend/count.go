package main

import (
	"strings"
	"unicode"
	"unicode/utf8"
)

func CitiesStartingWith(names []string, letter string) []string {
	want, _ := utf8.DecodeRuneInString(letter)
	want = unicode.ToLower(want)

	out := []string{}
	for _, n := range names {
		first, _ := utf8.DecodeRuneInString(strings.TrimSpace(n))
		if unicode.ToLower(first) == want {
			out = append(out, n)
		}
	}
	return out
}
