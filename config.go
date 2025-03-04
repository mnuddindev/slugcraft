package slugcraft

import (
	"regexp"
	"strings"
	"sync"
)

// Config is the main struct for generating slugs.
type Config struct {
	MaxLength     int                 // Maximum allowed length of the final slug (e.g., 220 characters)
	SuffixStyle   string              // Style of suffix: "numeric" (-2), "version" (-v2), "revision" (-rev2)
	Language      string              // Language will hold the preferred Language to transliteration Default: english
	RegexReplace  string              // Will hold the things that will be replaced
	StopWords     map[string]struct{} // All words that will be removed from the input if given
	Abbreviations map[string]string   // Abbreviations that will be removed from the input if given
	UseCache      bool                // Flag to enable in-memory caching of slug lookups
	Cache         *Cache              // In-memory cache struct
	RegexFilter   *regexp.Regexp      // Regex pattern to replace certain characters from input if given
	PipeLine      []Transformer       // Pipeline for step by step process
}

// Cache is a simple in-memory store for slug uniqueness.
type Cache struct {
	Mu    sync.RWMutex
	Store map[string]struct{}
}

// Option defines a functional option for configuring Config.
type Options func(*Config)

// Transformer defines a function that transform a string in pipeline.
type Transformer func(string) string

// New creates a new Config with default settings and optional configurations.
func New(options ...Options) *Config {
	cfg := &Config{
		PipeLine: []Transformer{
			strings.ToLower,
			func(s string) string { return strings.ReplaceAll(s, " ", "-") },
		},
		MaxLength:   220,
		UseCache:    false,
		SuffixStyle: "numeric",
		Cache:       &Cache{Store: make(map[string]struct{})},
	}
	for _, opt := range options {
		opt(cfg)
	}
	return cfg
}

// WithPipeline sets custom transformations for the slug generation pipeline.
func WithPipeline(transformers ...Transformer) Options {
	return func(cfg *Config) {
		cfg.PipeLine = transformers
	}
}

// WithLanguage sets the language for transliteration.
func WithLanguage(lang string) Options {
	return func(cfg *Config) {
		cfg.Language = lang
	}
}

// WithUseCache enables or disables in-memory caching for uniqueness
func WithUseCache(use bool) Options {
	return func(cfg *Config) {
		cfg.UseCache = use
	}
}

// WithSuffixStyle sets the style for suffix generation ("numeric", "version", "revision")
func WithSuffixStyle(style string) Options {
	return func(cfg *Config) {
		switch style {
		case "numeric", "version", "revision":
			cfg.SuffixStyle = style
		default:
			cfg.SuffixStyle = "numeric"
		}
	}
}

// WithMaxLength sets the maximum length of the slug
func WithMaxLength(max int) Options {
	return func(cfg *Config) {
		if max > 0 {
			cfg.MaxLength = max
		}
	}
}

// WithRegexFilter sets a regex pattern to filter characters
func WithRegexFilter(pattern, replace string) Options {
	return func(cfg *Config) {
		cfg.RegexFilter = regexp.MustCompile(pattern)
		cfg.RegexReplace = replace
	}
}

// WithStopWords sets stopwords to remove from the slug
func WithStopWords(lang string) Options {
	return func(cfg *Config) {
		cfg.StopWords = DefaultStopWords(lang)
	}
}
