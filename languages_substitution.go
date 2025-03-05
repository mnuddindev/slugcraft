package slugcraft

import (
	"strings"
)

// TransliterateBangla converts Bengali text to Banglish.
func TransliterateBangla(input string) string {
	var b strings.Builder
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
		"ৎ": "t", "ড়": "r", "ঢ়": "rh", "য়": "y", "ং": "ng",

		// Dependent Vowel Signs
		"া": "a", "ি": "i", "ী": "ee", "ু": "u", "ূ": "oo",
		"ে": "e", "ৈ": "oi", "ো": "o", "ৌ": "ou",

		// Jukto Borno (Conjunct Consonants)
		"ক্ত": "kt", "গ্ন": "gn", "স্ট": "st", "স্প": "sp", "শ্চ": "sch", "স্ফ": "sph",
		"স্ত": "st", "স্ত্র": "str", "ন্ত্র": "ntr", "ম্প": "mp", "ন্ড": "nd",
		"ঙ্ক": "nk", "ঙ্গ": "ngg", "ষ্ক": "shk", "ষ্ঠ": "shth", "ক্ষ": "kh",
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
