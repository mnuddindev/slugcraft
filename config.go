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
			Lowercase(),
			func(s string) string {
				re := regexp.MustCompile(`[^a-z0-9]+`)
				return re.ReplaceAllString(s, "-")
			},
			func(s string) string {
				return strings.Trim(s, "-")
			},
		},
		MaxLength:   220,
		UseCache:    false,
		SuffixStyle: "numeric",
		Cache:       &Cache{Store: make(map[string]struct{})},
	}
	for _, opt := range options {
		opt(cfg)
	}
	if cfg.MaxLength <= 0 {
		cfg.MaxLength = 220
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
		} else {
			cfg.MaxLength = 220
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

// WithAbbreviation adds a custom abbreviation rule.
func WithAbbreviation(from, to string) Options {
	return func(cfg *Config) {
		if cfg.Abbreviations == nil {
			cfg.Abbreviations = make(map[string]string)
		}
		cfg.Abbreviations[from] = to
	}
}

// DefaultStopWords returns a basic stopwords map.
func DefaultStopWords(lang string) map[string]struct{} {
	if lang == "en" {
		return map[string]struct{}{
			"a": {}, "about": {}, "above": {}, "after": {}, "again": {}, "against": {}, "all": {}, "am": {}, "an": {}, "and": {}, "any": {}, "are": {}, "aren't": {}, "as": {}, "at": {},
			"be": {}, "because": {}, "been": {}, "before": {}, "being": {}, "below": {}, "between": {}, "both": {}, "but": {}, "by": {},
			"can't": {}, "cannot": {}, "could": {}, "couldn't": {},
			"did": {}, "didn't": {}, "do": {}, "does": {}, "doesn't": {}, "doing": {}, "don't": {}, "down": {}, "during": {},
			"each": {},
			"few":  {}, "for": {}, "from": {}, "further": {},
			"had": {}, "hadn't": {}, "has": {}, "hasn't": {}, "have": {}, "haven't": {}, "having": {}, "he": {}, "he'd": {}, "he'll": {}, "he's": {}, "her": {}, "here": {}, "here's": {}, "hers": {}, "herself": {}, "him": {}, "himself": {}, "his": {}, "how": {}, "how's": {},
			"i": {}, "i'd": {}, "i'll": {}, "i'm": {}, "i've": {}, "if": {}, "in": {}, "into": {}, "is": {}, "isn't": {}, "it": {}, "it's": {}, "its": {}, "itself": {},
			"let's": {},
			"me":    {}, "more": {}, "most": {}, "mustn't": {}, "my": {}, "myself": {},
			"no": {}, "nor": {}, "not": {}, "of": {}, "off": {}, "on": {}, "once": {}, "only": {}, "or": {}, "other": {}, "ought": {}, "our": {}, "ours": {}, "ourselves": {}, "out": {}, "over": {}, "own": {},
			"same": {}, "shan't": {}, "she": {}, "she'd": {}, "she'll": {}, "she's": {}, "should": {}, "shouldn't": {}, "so": {}, "some": {}, "such": {},
			"than": {}, "that": {}, "that's": {}, "the": {}, "their": {}, "theirs": {}, "them": {}, "themselves": {}, "then": {}, "there": {}, "there's": {}, "these": {}, "they": {}, "they'd": {}, "they'll": {}, "they're": {}, "they've": {}, "this": {}, "those": {}, "through": {}, "to": {}, "too": {},
			"under": {}, "until": {}, "up": {},
			"very": {},
			"was":  {}, "wasn't": {}, "we": {}, "we'd": {}, "we'll": {}, "we're": {}, "we've": {}, "were": {}, "weren't": {}, "what": {}, "what's": {}, "when": {}, "when's": {}, "where": {}, "where's": {}, "which": {}, "while": {}, "who": {}, "who's": {}, "whom": {}, "why": {}, "why's": {}, "with": {}, "won't": {}, "would": {}, "wouldn't": {},
			"you": {}, "you'd": {}, "you'll": {}, "you're": {}, "you've": {}, "your": {}, "yours": {}, "yourself": {}, "yourselves": {},
		}
	}
	return nil
}
