package slugcraft

import (
	"strings"
	"unicode"
)

// TransliterateBangla converts Bengali text to Banglish.
func TransliterateBangla(input string, b *strings.Builder) string {
	runes := []rune(input)
	bengaliToBanglish := map[string]string{
		// Vowels
		"অ": "o", "আ": "a", "ই": "i", "ঈ": "ee", "উ": "u", "ঊ": "oo",
		"এ": "e", "ঐ": "oi", "ও": "o", "ঔ": "ou",

		// Consonants
		"ক": "k", "খ": "kh", "গ": "g", "ঘ": "gh", "ঙ": "ng",
		"চ": "ch", "ছ": "chh", "জ": "j", "ঝ": "jh", "ঞ": "ny",
		"ট": "t", "ঠ": "th", "ড": "d", "ঢ": "dh", "ণ": "n",
		"ত": "t", "থ": "th", "দ": "d", "ধ": "dh", "ন": "n",
		"প": "p", "ফ": "ph", "ব": "b", "ভ": "bh", "ম": "m",
		"য": "j", "র": "r", "ল": "l", "শ": "sh", "ষ": "sh", "স": "s", "হ": "h",

		// Special Cases
		"ৎ": "t", "ড়": "r", "ঢ়": "rh", "য়": "yo", "ং": "ng",

		// Dependent Vowel Signs
		"া": "a", "ি": "i", "ী": "ee", "ু": "u", "ূ": "oo",
		"ে": "e", "ৈ": "oi", "ো": "o", "ৌ": "ou", "্র": "r",

		// Jukto Borno (Conjunct Consonants)
		"ক্ত": "kt", "গ্ন": "gn", "স্ট": "st", "স্প": "sp", "শ্চ": "sch", "স্ফ": "sph",
		"স্ত": "st", "স্ত্র": "str", "ন্ত্র": "ntr", "ম্প": "mp", "ন্ড": "nd",
		"ঙ্ক": "nk", "ঙ্গ": "ngg", "ষ্ক": "shk", "ষ্ঠ": "shth", "ক্ষ": "kho",
	}

	for i := 0; i < len(runes); {
		// Check for three-rune conjunct (e.g., স্ত্র)
		if i+2 < len(runes) {
			triplet := string(runes[i : i+3])
			if mapped, ok := bengaliToBanglish[triplet]; ok {
				b.WriteString(mapped)
				i += 3
				continue
			}
		}

		// Check for two-rune conjunct (e.g., ক্ত)
		if i+1 < len(runes) {
			pair := string(runes[i : i+2])
			if mapped, ok := bengaliToBanglish[pair]; ok {
				b.WriteString(mapped)
				i += 2
				continue
			}
		}

		// Single rune
		current := string(runes[i])
		if mapped, ok := bengaliToBanglish[current]; ok {
			b.WriteString(mapped)
		} else {
			b.WriteRune(runes[i])
		}
		i++
	}

	return b.String()
}

// transliterateCyrillic handles Cyrillic script (Russian).
func TransliterateRussian(input string, b *strings.Builder) {
	for _, r := range input {
		switch r {
		case 'а':
			b.WriteString("a")
		case 'б':
			b.WriteString("b")
		case 'в':
			b.WriteString("v")
		case 'г':
			b.WriteString("g")
		case 'д':
			b.WriteString("d")
		case 'е':
			b.WriteString("e")
		case 'ё':
			b.WriteString("yo")
		case 'ж':
			b.WriteString("zh")
		case 'з':
			b.WriteString("z")
		case 'и':
			b.WriteString("i")
		case 'й':
			b.WriteString("y")
		case 'к':
			b.WriteString("k")
		case 'л':
			b.WriteString("l")
		case 'м':
			b.WriteString("m")
		case 'н':
			b.WriteString("n")
		case 'о':
			b.WriteString("o")
		case 'п':
			b.WriteString("p")
		case 'р':
			b.WriteString("r")
		case 'с':
			b.WriteString("s")
		case 'т':
			b.WriteString("t")
		case 'у':
			b.WriteString("u")
		case 'ф':
			b.WriteString("f")
		case 'х':
			b.WriteString("kh")
		case 'ц':
			b.WriteString("ts")
		case 'ч':
			b.WriteString("ch")
		case 'ш':
			b.WriteString("sh")
		case 'щ':
			b.WriteString("shch")
		case 'ъ':
			b.WriteString("")
		case 'ы':
			b.WriteString("y")
		case 'ь':
			b.WriteString("")
		case 'э':
			b.WriteString("e")
		case 'ю':
			b.WriteString("yu")
		case 'я':
			b.WriteString("ya")
		case ' ':
			b.WriteByte(' ')
		default:
			if unicode.Is(unicode.Cyrillic, r) {
				b.WriteRune(unicode.ToLower(r))
			}
		}
	}
}

// transliterateGeneric handles basic Latin normalization.
func TransliterateGeneric(input string, b *strings.Builder) {
	for _, r := range input {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			b.WriteRune(unicode.ToLower(r))
		} else if unicode.IsSpace(r) {
			b.WriteByte(' ')
		}
	}
}

// transliterateUnidecode default (no dependency) version.
func TransliterateUnidecode(input string, b *strings.Builder) {
	TransliterateGeneric(input, b)
}

// transliterateASCIISafe ensures a non-empty ASCII output as a last resort.
func transliterateASCIISafe(input string, b *strings.Builder) {
	for _, r := range input {
		if r <= 127 { // ASCII range
			b.WriteRune(unicode.ToLower(r))
		}
	}
	if b.Len() == 0 {
		b.WriteString("slug") // Ultimate fallback
	}
}
