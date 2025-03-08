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

// Lowercase converts text to lowercase in-place.
func Lowercase() Transformer {
	return func(b *strings.Builder) {
		temp := strings.ToLower(b.String())
		b.Reset()
		b.WriteString(temp)
	}
}

// ReplaceSpaces replaces spaces with dashes.
func ReplaceSpaces(delimeter string) Transformer {
	return func(b *strings.Builder) {
		re := regexp.MustCompile(`[^a-z0-9]+`)
		temp := re.ReplaceAllString(b.String(), delimeter)
		b.Reset()
		b.WriteString(strings.Trim(temp, delimeter))
	}
}

// RemoveDiacritics removes diacritics using Unicode normalization.
func RemoveDiacritics() Transformer {
	return func(b *strings.Builder) {
		t := transform.Chain(
			norm.NFD,
			runes.Remove(runes.In(unicode.Mn)),
		)
		result, _, _ := transform.String(t, b.String())
		b.Reset()
		b.WriteString(result)
	}
}

// TrimDashes trims leading and trailing dashes.
func TrimDashes() Transformer {
	return func(b *strings.Builder) {
		temp := strings.Trim(b.String(), "-")
		b.Reset()
		b.WriteString(temp)
	}
}
