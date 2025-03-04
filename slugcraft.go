package slugcraft

import (
	"context"
	"fmt"
	"strings"

	"golang.org/x/text/transform"
)

// Make generates a slug from the input string with the configured options.
func (cfg *Config) Make(ctx context.Context, input string) (string, error) {
	if input == "" {
		return "", nil
	}

	if err := ctx.Err(); err != nil {
		return "", err
	}

	slug := input

	// Apply abbreviaton
	if cfg.Abbreviations != nil {
		for from, to := range cfg.Abbreviations {
			slug = strings.ReplaceAll(slug, from, to)
		}
	}

	// Remove stopwords
	if cfg.StopWords != nil {
		words := strings.Fields(slug)
		var filtered []string
		for _, w := range words {
			if _, ok := cfg.StopWords[strings.ToLower(w)]; !ok {
				filtered = append(filtered, w)
			}
		}
		slug = strings.Join(filtered, " ")
	}

	// Apply language-specific transliteration
	if cfg.Language != "" {
		var err error
		slug, err = cfg.Transliterate(slug, cfg.Language)
		if err != nil {
			return "", nil
		}
	}

	// Apply pipeline transformations
	for _, t := range cfg.PipeLine {
		slug = t(slug)
	}

	// Apply regex filter
	if cfg.RegexFilter != nil {
		slug = cfg.RegexFilter.ReplaceAllString(slug, cfg.RegexReplace)
	}

	// Truncate to max length
	if cfg.MaxLength > 0 && len(slug) > cfg.MaxLength {
		slug = slug[:cfg.MaxLength]
	}

	// Handle uniqueness with in-memory cache
	if cfg.UseCache {
		slug = cfg.EnsureUnique(ctx, slug)
		cfg.Cache.Set(slug)
	}

	return slug, nil
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
	trans := LanguageTransformer(lang)
	result, _, err := transform.String(trans, input)
	if err != nil {
		return input, err
	}
	return result, nil
}
