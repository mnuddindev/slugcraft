package slugcraft

import (
	"context"
	"strings"
)

// Make generates a slug from the input string with the configured options.
func (cfg *Config) Make(ctx context.Context, input string) (string, error) {
	if input == "" {
		return "", nil
	}
	if err := ctx.Err(); err != nil {
		return "", err
	}

	// Initialize builder for both modes
	cfg.Builder.Reset()
	cfg.Builder.Grow(len(input))
	cfg.Builder.WriteString(input)

	// Apply abbreviations
	if cfg.Abbreviations != nil {
		var temp strings.Builder
		temp.WriteString(cfg.Builder.String())
		for from, to := range cfg.Abbreviations {
			tempStr := strings.ReplaceAll(temp.String(), from, to)
			temp.Reset()
			temp.WriteString(tempStr)
		}
		cfg.Builder.Reset()
		cfg.Builder.WriteString(temp.String())
	}

	// Remove stopwords
	if cfg.StopWords != nil {
		words := strings.Fields(cfg.Builder.String())
		cfg.Builder.Reset()
		for i, w := range words {
			if _, ok := cfg.StopWords[strings.ToLower(w)]; !ok {
				if i > 0 {
					cfg.Builder.WriteByte(' ')
				}
				cfg.Builder.WriteString(w)
			}
		}
	}

	// Apply language-specific transliteration
	if cfg.Language != "" {
		result, err := cfg.Transliterate(cfg.Builder.String()) // Use builder output only
		if err != nil {
			return "", err // Fix: return error
		}
		cfg.Builder.Reset()
		cfg.Builder.WriteString(result)
	}

	// Apply pipeline transformations
	for _, t := range cfg.PipeLine {
		t(&cfg.Builder)
	}

	// Apply regex filter
	if cfg.RegexFilter != nil {
		temp := cfg.RegexFilter.ReplaceAllString(cfg.Builder.String(), cfg.RegexReplace)
		cfg.Builder.Reset()
		cfg.Builder.WriteString(temp)
	}

	// Truncate to max length
	if cfg.MaxLength > 0 && cfg.Builder.Len() > cfg.MaxLength {
		runes := []rune(cfg.Builder.String())
		if len(runes) > cfg.MaxLength {
			cfg.Builder.Reset()
			cfg.Builder.WriteString(string(runes[:cfg.MaxLength]))
		}
	}

	// Handle uniqueness with in-memory cache
	if cfg.UseCache {
		cfg.EnsureUnique(ctx)
	}

	// Return final result
	result := cfg.Builder.String()
	if !cfg.ZeroAlloc {
		// Non-zero-alloc mode creates a copy (for compatibility)
		slug := make([]byte, len(result))
		copy(slug, result)
		return string(slug), nil
	}
	return result, nil
}

// MakeBulk generates slugs for multilple inputs.
func (cfg *Config) MakeBulk(ctx context.Context, inputs []string) ([]string, error) {
	if len(inputs) == 0 {
		return nil, nil
	}
	slugs := make([]string, len(inputs))
	for i, input := range inputs {
		slug, err := cfg.Make(ctx, input)
		if err != nil {
			return nil, err
		}
		slugs[i] = slug
	}
	return slugs, nil
}

// EnsureUnique ensures the slug is unique using the in-memory cache.
func (cfg *Config) EnsureUnique(ctx context.Context) {
	if err := ctx.Err(); err != nil {
		return
	}

	baseSlug := cfg.Builder.String()
	cfg.Cache.Mu.Lock()
	defer cfg.Cache.Mu.Unlock()

	count, exists := cfg.Cache.Store[baseSlug]
	if !exists {
		cfg.Cache.Store[baseSlug] = 0
		return
	}

	count++
	cfg.Cache.Store[baseSlug] = count
	cfg.Builder.WriteByte('-')
	switch cfg.SuffixStyle {
	case "numeric":
		cfg.Builder.WriteString(itoa(count))
	case "version":
		cfg.Builder.WriteString("v")
		cfg.Builder.WriteString(itoa(count))
	case "revision":
		cfg.Builder.WriteString("rev")
		cfg.Builder.WriteString(itoa(count))
	}
	finalSlug := cfg.Builder.String()
	cfg.Cache.Store[finalSlug] = 0
}

// itoa converts an int to string without allocation
func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	var digits [20]byte
	i := len(digits) - 1
	for n > 0 {
		digits[i] = byte('0' + n%10)
		n /= 10
		i--
	}
	return string(digits[i+1:])
}

// Transliterate converts text to a Latin-based slug using language-specific rules
func (cfg *Config) Transliterate(input string) (string, error) {
	cfg.Builder.Reset()
	cfg.Builder.Grow(len(input))

	switch cfg.Language {
	case "bn":
		TransliterateBangla(input, &cfg.Builder)
	case "ru":
		TransliterateRussian(input, &cfg.Builder)
	default:
		if cfg.UseUnidecode {
			TransliterateUnidecode(input, &cfg.Builder)
		} else {
			TransliterateGeneric(input, &cfg.Builder)
		}
	}

	// Fail-safe ASCII fallback
	if cfg.Builder.Len() == 0 {
		transliterateASCIISafe(input, &cfg.Builder)
	}

	return cfg.Builder.String(), nil
}
