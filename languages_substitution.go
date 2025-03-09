package slugcraft

import (
	"strings"
	"unicode"
	"unicode/utf8"
)

// TransliterateBangla converts Bengali text to Banglish.
func TransliterateBangla(input string, b *strings.Builder) {
	var lastRune rune
	for i := 0; i < len(input); {
		r, size := utf8.DecodeRuneInString(input[i:])
		switch r {
		// Vowels
		case 'অ':
			b.WriteString("o")
		case 'আ':
			b.WriteString("a")
		case 'ই':
			b.WriteString("i")
		case 'ঈ':
			b.WriteString("i")
		case 'উ':
			b.WriteString("u")
		case 'ঊ':
			b.WriteString("u")
		case 'ঋ':
			b.WriteString("ri")
		case 'এ':
			b.WriteString("e")
		case 'ঐ':
			b.WriteString("oi")
		case 'ও':
			b.WriteString("o")
		case 'ঔ':
			b.WriteString("ou")
		// Consonants
		case 'ক':
			b.WriteString("k")
		case 'খ':
			b.WriteString("kh")
		case 'গ':
			b.WriteString("g")
		case 'ঘ':
			b.WriteString("gh")
		case 'ঙ':
			b.WriteString("ng")
		case 'চ':
			b.WriteString("ch")
		case 'ছ':
			b.WriteString("chh")
		case 'জ':
			b.WriteString("j")
		case 'ঝ':
			b.WriteString("jh")
		case 'ঞ':
			b.WriteString("n")
		case 'ট':
			b.WriteString("t")
		case 'ঠ':
			b.WriteString("th")
		case 'ড':
			b.WriteString("d")
		case 'ঢ':
			b.WriteString("dh")
		case 'ণ':
			b.WriteString("n")
		case 'ত':
			b.WriteString("t")
		case 'থ':
			b.WriteString("th")
		case 'দ':
			b.WriteString("d")
		case 'ধ':
			b.WriteString("dh")
		case 'ন':
			b.WriteString("n")
		case 'প':
			b.WriteString("p")
		case 'ফ':
			b.WriteString("ph")
		case 'ব':
			b.WriteString("b")
		case 'ভ':
			b.WriteString("bh")
		case 'ম':
			b.WriteString("m")
		case 'য':
			b.WriteString("y")
		case 'র':
			b.WriteString("r")
		case 'ল':
			b.WriteString("l")
		case 'শ':
			b.WriteString("sh")
		case 'ষ':
			b.WriteString("sh")
		case 'স':
			b.WriteString("s")
		case 'হ':
			b.WriteString("h")
		// Diacritics (vowel signs)
		case 'ি':
			if lastRune != 0 {
				b.WriteString("i")
			}
		case 'ী':
			if lastRune != 0 {
				b.WriteString("i")
			}
		case 'ু':
			if lastRune != 0 {
				b.WriteString("u")
			}
		case 'ূ':
			if lastRune != 0 {
				b.WriteString("u")
			}
		case 'ৃ':
			if lastRune != 0 {
				b.WriteString("ri")
			}
		case 'ে':
			if lastRune != 0 {
				b.WriteString("e")
			}
		case 'ৈ':
			if lastRune != 0 {
				b.WriteString("oi")
			}
		case 'ো':
			if lastRune != 0 {
				b.WriteString("o")
			}
		case 'ৌ':
			if lastRune != 0 {
				b.WriteString("ou")
			}
		case 'ৎ':
			if lastRune != 0 {
				b.WriteString("t")
			}
		case 'ড়':
			if lastRune != 0 {
				b.WriteString("r")
			}
		case 'ঢ়':
			if lastRune != 0 {
				b.WriteString("rh")
			}
		case 'য়':
			if lastRune != 0 {
				b.WriteString("yo")
			}
		case 'ং':
			if lastRune != 0 {
				b.WriteString("ng")
			}
		// Special cases (conjuncts)
		case '্': // Halant: check for conjuncts
			if lastRune != 0 && i+size < len(input) {
				nextR, nextSize := utf8.DecodeRuneInString(input[i+size:])
				switch string([]rune{lastRune, r, nextR}) {
				case "ক্ষ":
					b.WriteString("kho")
					size += nextSize
					lastRune = 0
					continue
				case "জ্ঞ":
					b.WriteString("gy")
					size += nextSize
					lastRune = 0
					continue
				case "ক্ত":
					b.WriteString("kt")
					size += nextSize
					lastRune = 0
					continue
				case "গ্ন":
					b.WriteString("gn")
					size += nextSize
					lastRune = 0
					continue
				case "স্ট":
					b.WriteString("st")
					size += nextSize
					lastRune = 0
					continue
				case "স্প":
					b.WriteString("sp")
					size += nextSize
					lastRune = 0
					continue
				case "শ্চ":
					b.WriteString("sch")
					size += nextSize
					lastRune = 0
					continue
				case "স্ফ":
					b.WriteString("sph")
					size += nextSize
					lastRune = 0
					continue
				case "স্ত":
					b.WriteString("st")
					size += nextSize
					lastRune = 0
					continue
				case "স্ত্র":
					b.WriteString("str")
					size += nextSize
					lastRune = 0
					continue
				case "ন্ত্র":
					b.WriteString("ntr")
					size += nextSize
					lastRune = 0
					continue
				case "ম্প":
					b.WriteString("mp")
					size += nextSize
					lastRune = 0
					continue
				case "ন্ড":
					b.WriteString("nd")
					size += nextSize
					lastRune = 0
					continue
				case "ঙ্ক":
					b.WriteString("nk")
					size += nextSize
					lastRune = 0
					continue
				case "ঙ্গ":
					b.WriteString("ng")
					size += nextSize
					lastRune = 0
					continue
				case "ষ্ক":
					b.WriteString("shk")
					size += nextSize
					lastRune = 0
					continue
				case "ষ্ঠ":
					b.WriteString("shth")
					size += nextSize
					lastRune = 0
					continue
				}
			}
		default:
			if unicode.IsSpace(r) {
				b.WriteByte(' ')
			} else if unicode.Is(unicode.Bengali, r) {
				b.WriteRune(unicode.ToLower(r)) // Fallback for unmapped
			}
		}
		lastRune = r
		i += size
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
