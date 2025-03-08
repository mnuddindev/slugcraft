package slugcraft

import (
	"context"
	"fmt"
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
		result, err := cfg.Transliterate(cfg.Builder.String(), cfg.Language) // Use builder output only
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
		slug := cfg.EnsureUnique(ctx, cfg.Builder.String())
		cfg.Builder.Reset()
		cfg.Builder.WriteString(slug)
		cfg.Cache.Set(slug)
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
func (cfg *Config) EnsureUnique(ctx context.Context, slug string) string {
	if !cfg.Cache.Get(slug) {
		return slug // Slug is free
	}

	base := slug
	for i := 1; ; i++ {
		if err := ctx.Err(); err != nil {
			return ""
		}
		var candidate string
		switch cfg.SuffixStyle {
		case "numeric":
			candidate = fmt.Sprintf("%s-%d", base, i)
		case "version":
			candidate = fmt.Sprintf("%s-v%d", base, i)
		case "revision":
			candidate = fmt.Sprintf("%s-rev%d", base, i)
		}
		if !cfg.Cache.Get(candidate) {
			return candidate // When find a new one
		}
	}
}

// Transliterate converts text to a Latin-based slug using language-specific rules
func (cfg *Config) Transliterate(input, lang string) (string, error) {
	var result string
	switch lang {
	case "bn":
		result = TransliterateBangla(input)
	}
	return result, nil
}
