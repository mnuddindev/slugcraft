package slugcraft

import (
	"strings"
	"unicode"
)

// Common transformers for pipeline

// Lowercase converts the string to lowercase
func Lowercase() Transformer {
	return strings.ToLower
}

// RemoveDiacritics removes diacritics
func RemoveDiacritics() Transformer {
	return func(s string) string {
		return strings.Map(func(r rune) rune {
			if unicode.Is(unicode.Mn, r) {
				return -1
			}
			return r
		}, s)
	}
}

// ReplaceSpaces replaces spaces with a delimeter
func ReplaceSpaces(delimeter string) Transformer {
	return func(s string) string {
		return strings.ReplaceAll(s, " ", delimeter)
	}
}
