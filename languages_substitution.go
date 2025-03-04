package slugcraft

import (
	"strings"
	"unicode"

	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
)

func LanguageTransformer(lang string) transform.Transformer {
	switch strings.ToLower(lang) {
	case "en":
		return transform.Nop
	case "bn":
		return transform.Chain(BengaliToTransformer(), StripNonLatin())
	default:
		return transform.Chain(StripNonLatin())
	}
}

// stripNonLatin removes any remaining non-Latin characters.
func StripNonLatin() transform.Transformer {
	return runes.Map(func(r rune) rune {
		if unicode.IsLetter(r) && !unicode.Is(unicode.Latin, r) {
			return -1 // Remove non-Latin letters
		}
		return r
	})
}

func BengaliToTransformer() transform.Transformer {
	return &RuneToStringTransformer{
		mapping: map[rune]string{
			// Vowels (স্বরবর্ণ)
			'অ': "o", 'আ': "a", 'ই': "i", 'ঈ': "i", 'উ': "u",
			'ঊ': "u", 'ঋ': "ri", 'এ': "e", 'ঐ': "oi", 'ও': "o",
			'ঔ': "ou",

			// Consonants (ব্যঞ্জনবর্ণ) with inherent "a"
			'ক': "k", 'খ': "kh", 'গ': "g", 'ঘ': "gh", 'ঙ': "ng",
			'চ': "ch", 'ছ': "chh", 'জ': "j", 'ঝ': "jh", 'ঞ': "ng",
			'ট': "t", 'ঠ': "th", 'ড': "d", 'ঢ': "dh", 'ণ': "n",
			'ত': "t", 'থ': "th", 'দ': "d", 'ধ': "dh", 'ন': "n",
			'প': "p", 'ফ': "ph", 'ব': "b", 'ভ': "bh", 'ম': "m",
			'য': "z", 'র': "r", 'ল': "l", 'শ': "sh", 'ষ': "sh",
			'স': "s", 'হ': "h", 'ড়': "r", 'ঢ়': "rh", 'য়': "y",
			'ৎ': "t",

			// Vowel Signs (কার)
			'া': "a", 'ি': "i", 'ী': "i", 'ু': "u", 'ূ': "u",
			'ৃ': "ri", 'ে': "e", 'ৈ': "oi", 'ো': "o", 'ৌ': "ou",

			// Special Characters
			'ং': "ng", // অনুস্বার
			'ঃ': "h",  // বিসর্গ
			'ঁ': "",   // চন্দ্রবিন্দু
			'্': "",   // হসন্ত
		},
	}
}

// runeToStringTransformer implements transform.Transformer for multi-character mappings.
type RuneToStringTransformer struct {
	mapping map[rune]string
}

func (t *RuneToStringTransformer) Transform(dst, src []byte, atEOF bool) (nDst, nSrc int, err error) {
	var b strings.Builder
	input := string(src)

	for _, r := range input {
		if mapped, ok := t.mapping[r]; ok {
			b.WriteString(mapped)
		} else {
			b.WriteRune(r)
		}
	}

	result := b.String()
	if len(result) > len(dst) {
		return 0, 0, transform.ErrShortDst
	}
	n := copy(dst, result)
	return n, len(src), nil
}

func (t *RuneToStringTransformer) Reset() {}
