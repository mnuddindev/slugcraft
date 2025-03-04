package slugcraft

import (
	"regexp"
	"strings"
	"unicode"

	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

// Common transformers for pipeline

// Lowercase converts the string to lowercase
func Lowercase() Transformer {
	return strings.ToLower
}

// RemoveDiacritics removes diacritics
func RemoveDiacritics() Transformer {
	return func(s string) string {
		t := transform.Chain(
			norm.NFD,
			runes.Remove(runes.In(unicode.Mn)),
		)
		result, _, _ := transform.String(t, s)
		return result
	}
}

// ReplaceSpaces replaces spaces with a delimeter
func ReplaceSpaces(delimeter string) Transformer {
	return func(s string) string {
		re := regexp.MustCompile(`[^a-z0-9]+`)
		s = re.ReplaceAllString(s, delimeter)
		return strings.Trim(s, delimeter)
	}
}
