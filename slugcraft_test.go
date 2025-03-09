package slugcraft

import (
	"context"
	"testing"
)

// TestNew tests the default Config configuration
func TestNew(t *testing.T) {
	s := New()
	if len(s.PipeLine) != 3 {
		t.Errorf("expected 2 default pipeline transformers, got %d", len(s.PipeLine))
	}
	if s.MaxLength != 220 {
		t.Errorf("expected MaxLength 220, got %d", s.MaxLength)
	}
	if s.UseCache != false {
		t.Errorf("expected UseCache false, got %v", s.UseCache)
	}
	if s.SuffixStyle != "numeric" {
		t.Errorf("expected SuffixStyle 'numeric', got %s", s.SuffixStyle)
	}
}

// TestMakeBasic tests basic slug generation
func TestMakeBasic(t *testing.T) {
	s := New()
	tests := []struct {
		Input    string
		Expected string
	}{
		{"Hello, World", "hello-world"},
		{"", ""},
		{"Simple", "simple"},
		{"UPPER CASE", "upper-case"},
	}

	for _, tt := range tests {
		t.Run(tt.Input, func(t *testing.T) {
			slug, err := s.Make(context.Background(), tt.Input)
			if err != nil {
				t.Errorf("Make(%q) returned error: %v", tt.Input, err)
			}
			if slug != tt.Expected {
				t.Errorf("Make(%q) = %q , expected %q", tt.Input, slug, err)
			}
		})
	}
}

// TestMakeWithLanguage tests slug generation with transliteration.
func TestMakeWithLanguage(t *testing.T) {
	tests := []struct {
		language string
		input    string
		expected string
	}{
		{"bn", "আমি তোমাকে", "ami-tomake"},
		{"bn", "বাংলা ভাষা", "bangla-bhasha"},
		{"bn", "গোলাপ ফুল", "golap-phul"},
		{"bn", "পাখির গান", "pakhir-gan"},
		{"bn", "রাতের তারা", "rater-tara"},
		{"bn", "ক্ষমা করো", "khoma-kro"},
		{"ru", "привет мир", "privet-mir"},
	}

	for _, tt := range tests {
		t.Run(tt.language+"/"+tt.input, func(t *testing.T) {
			s := New(WithLanguage(tt.language))
			slug, err := s.Make(context.Background(), tt.input)
			if err != nil {
				t.Errorf("Make(%q, lang=%q) returned error: %v", tt.input, tt.language, err)
			}
			if slug != tt.expected {
				t.Errorf("Make(%q, lang=%q) = %q, expected %q", tt.input, tt.language, slug, tt.expected)
			}
		})
	}
}

// TestMakeWithPipeline tests custom pipeline transformations.
func TestMakeWithPipeline(t *testing.T) {
	s := New(WithPipeline(
		Lowercase(),
		RemoveDiacritics(),
		ReplaceSpaces("_"), // Custom delimiter
	))
	tests := []struct {
		input    string
		expected string
	}{
		{"Hello, World!", "hello_world"},
		{"Café au Lait", "cafe_au_lait"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			slug, err := s.Make(context.Background(), tt.input)
			if err != nil {
				t.Errorf("Make(%q) returned error: %v", tt.input, err)
			}
			if slug != tt.expected {
				t.Errorf("Make(%q) = %q, expected %q", tt.input, slug, tt.expected)
			}
		})
	}
}

// TestMakeWithCache tests collision avoidance with in-memory cache.
func TestMakeWithCache(t *testing.T) {
	s := New(
		WithUseCache(true),
		WithSuffixStyle("numeric"),
	)
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"First", "My Post", "my-post"},
		{"Collision1", "My Post", "my-post-1"}, // Collision
		{"Collision2", "My Post", "my-post-2"}, // Another collision
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			slug, err := s.Make(context.Background(), tt.input)
			if err != nil {
				t.Errorf("Make(%q) returned error: %v", tt.input, err)
			}
			if slug != tt.expected {
				t.Errorf("Make(%q) = %q, expected %q", tt.input, slug, tt.expected)
			}
		})
	}
}

// TestMakeWithCacheAndStyles tests different suffix styles.
func TestMakeWithCacheAndStyles(t *testing.T) {
	tests := []struct {
		suffixStyle string
		inputs      []string
		expected    []string
	}{
		{
			"numeric",
			[]string{"Test", "Test", "Test"},
			[]string{"test", "test-1", "test-2"},
		},
		{
			"version",
			[]string{"Test", "Test", "Test"},
			[]string{"test", "test-v1", "test-v2"},
		},
		{
			"revision",
			[]string{"Test", "Test", "Test"},
			[]string{"test", "test-rev1", "test-rev2"},
		},
	}

	for _, tt := range tests {
		s := New(
			WithUseCache(true),
			WithSuffixStyle(tt.suffixStyle),
		)
		for i, input := range tt.inputs {
			slug, err := s.Make(context.Background(), input)
			if err != nil {
				t.Errorf("Make(%q, style=%q) returned error: %v", input, tt.suffixStyle, err)
			}
			if slug != tt.expected[i] {
				t.Errorf("Make(%q, style=%q) = %q, expected %q", input, tt.suffixStyle, slug, tt.expected[i])
			}
		}
	}
}

// TestMakeWithoutCache tests slug generation without caching.
func TestMakeWithoutCache(t *testing.T) {
	s := New(WithUseCache(false))
	tests := []struct {
		input    string
		expected string
	}{
		{"My Post", "my-post"},
		{"My Post", "my-post"}, // No collision avoidance
	}

	for _, tt := range tests {
		slug, err := s.Make(context.Background(), tt.input)
		if err != nil {
			t.Errorf("Make(%q) returned error: %v", tt.input, err)
		}
		if slug != tt.expected {
			t.Errorf("Make(%q) = %q, expected %q", tt.input, slug, tt.expected)
		}
	}
}

// TestMakeBulk tests bulk slug generation.
func TestMakeBulk(t *testing.T) {
	s := New()
	inputs := []string{"Post One", "Post Two", "Post Three"}
	expected := []string{"post-one", "post-two", "post-three"}

	slugs, err := s.MakeBulk(context.Background(), inputs)
	if err != nil {
		t.Errorf("MakeBulk returned error: %v", err)
	}
	if len(slugs) != len(expected) {
		t.Errorf("MakeBulk returned %d slugs, expected %d", len(slugs), len(expected))
	}
	for i, slug := range slugs {
		if slug != expected[i] {
			t.Errorf("MakeBulk[%d] = %q, expected %q", i, slug, expected[i])
		}
	}
}

// TestCacheOperations tests the in-memory cache directly.
func TestCacheOperations(t *testing.T) {
	c := &Cache{Store: make(map[string]int)}

	// Test set and get
	c.Set("slug1")
	if !c.Get("slug1") {
		t.Errorf("cache.get('slug1') = false, expected true")
	}
	if c.Get("slug2") {
		t.Errorf("cache.get('slug2') = true, expected false")
	}

	// Test delete
	c.Set("slug2")
	c.Del("slug2")
	if c.Get("slug2") {
		t.Errorf("cache.get('slug2') after delete = true, expected false")
	}
}

// TestMakeWithMaxLength tests slug truncation.
func TestMakeWithMaxLength(t *testing.T) {
	s := New(WithMaxLength(5))
	slug, err := s.Make(context.Background(), "Hello, World!")
	if err != nil {
		t.Errorf("Make returned error: %v", err)
	}
	if slug != "hello" {
		t.Errorf("Make = %q, expected %q", slug, "hello")
	}
	if len(slug) > 5 {
		t.Errorf("slug length = %d, expected <= 5", len(slug))
	}
}

// TestMakeWithContextCancel tests context cancellation.
func TestMakeWithContextCancel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	s := New()
	_, err := s.Make(ctx, "Hello, World!")
	if err == nil {
		t.Errorf("Make with canceled context did not return error")
	}
}

// TestMakeBangla tests basic Bangla transliteration.
func TestMakeBangla(t *testing.T) {
	s := New(WithLanguage("bn"))
	tests := []struct {
		input    string
		expected string
	}{
		{"বাংলা", "bangla"},
		{"আমি ভালো", "ami-bhalo"},
		{"ক্ত", "kt"},
		{"Hello বাংলা", "hello-bangla"},
		{"ষ্ঠান", "shthan"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			slug, err := s.Make(context.Background(), tt.input)
			if err != nil {
				t.Errorf("Make(%q) returned error: %v", tt.input, err)
			}
			if slug != tt.expected {
				t.Errorf("Make(%q) = %q, expected %q", tt.input, slug, tt.expected)
			}
		})
	}
}

// TestRegexFilter tests regex filter functionality.
func TestRegexFilter(t *testing.T) {
	s := New(
		WithRegexFilter(`[^a-z0-9-]`, ""), // Remove everything except a-z, 0-9, -
	)
	tests := []struct {
		input    string
		expected string
	}{
		{"Hello, World", "hello-world"},
		{"Asho Khela hobe", "asho-khela-hobe"},
		{"Hello, Bangla", "hello-bangla"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			slug, err := s.Make(context.Background(), tt.input)
			if err != nil {
				t.Errorf("Make(%q) returned error: %v", tt.input, err)
			}
			if slug != tt.expected {
				t.Errorf("Make(%q) = %q, expected %q", tt.input, slug, tt.expected)
			}
		})
	}
}

// TestStopwords tests stopword removal.
func TestStopwords(t *testing.T) {
	s := New(
		WithLanguage("bn"),
		WithStopWords("en"), // Using English stopwords for mixed text
	)
	tests := []struct {
		input    string
		expected string
	}{
		{"বাংলা the আমি", "bangla-ami"},
		{"প্রিয় and ক্ষমা", "priyo-khoma"},
		{"Hello a বাংলা", "hello-bangla"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			slug, err := s.Make(context.Background(), tt.input)
			if err != nil {
				t.Errorf("Make(%q) returned error: %v", tt.input, err)
			}
			if slug != tt.expected {
				t.Errorf("Make(%q) = %q, expected %q", tt.input, slug, tt.expected)
			}
		})
	}
}

// TestAbbreviations tests abbreviation replacement.
func TestAbbreviations(t *testing.T) {
	s := New(
		WithLanguage("bn"),
		WithAbbreviation("বাংলা", "BN"),
		WithAbbreviation("প্রিয়", "PR"),
	)
	tests := []struct {
		input    string
		expected string
	}{
		{"বাংলা আমি", "bn-ami"},
		{"প্রিয় ক্ষমা", "pr-khoma"},
		{"Hello বাংলা", "hello-bn"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			slug, err := s.Make(context.Background(), tt.input)
			if err != nil {
				t.Errorf("Make(%q) returned error: %v", tt.input, err)
			}
			if slug != tt.expected {
				t.Errorf("Make(%q) = %q, expected %q", tt.input, slug, tt.expected)
			}
		})
	}
}

// BenchmarkTransliterateBangla measures Bangla transliteration performance.
func BenchmarkTransliterateBangla(b *testing.B) {
	s := New(WithLanguage("bn"))
	input := "বাংলা প্রিয় ক্ষমা"
	for i := 0; i < b.N; i++ {
		s.Make(context.Background(), input)
	}
}

// BenchmarkPipeline measures pipeline transformation performance.
func BenchmarkPipeline(b *testing.B) {
	s := New(WithLanguage("bn"))
	input := "হ্যালো, বিশ্ব! Hello World"
	for i := 0; i < b.N; i++ {
		s.Make(context.Background(), input)
	}
}

// BenchmarkRegexFilter measures regex filter performance.
func BenchmarkRegexFilter(b *testing.B) {
	s := New(
		WithLanguage("bn"),
		WithRegexFilter(`[^a-z0-9-]`, ""),
	)
	input := "বাংলা!@# প্রিয় 123 ক্ষমা"
	for i := 0; i < b.N; i++ {
		s.Make(context.Background(), input)
	}
}

// BenchmarkStopwords measures stopword removal performance.
func BenchmarkStopwords(b *testing.B) {
	s := New(
		WithLanguage("bn"),
		WithStopWords("en"),
	)
	input := "বাংলা the প্রিয় and ক্ষমা"
	for i := 0; i < b.N; i++ {
		s.Make(context.Background(), input)
	}
}

// BenchmarkAbbreviations measures abbreviation replacement performance.
func BenchmarkAbbreviations(b *testing.B) {
	s := New(
		WithLanguage("bn"),
		WithAbbreviation("বাংলা", "BN"),
		WithAbbreviation("প্রিয়", "PR"),
	)
	input := "বাংলা প্রিয় ক্ষমা"
	for i := 0; i < b.N; i++ {
		s.Make(context.Background(), input)
	}
}

// BenchmarkCache measures cache-based collision avoidance performance.
func BenchmarkCache(b *testing.B) {
	cfg := New(WithLanguage("bn"), func(cfg *Config) {
		cfg.UseCache = true
	})
	ctx := context.Background()

	inputs := []string{
		"আমি তোমাকে", "বাংলা ভাষা", "ক্ষমা করো",
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, input := range inputs {
			_, _ = cfg.Make(ctx, input)
		}
	}
}

// BenchmarkMakeBulk measures bulk slug generation performance.
func BenchmarkMakeBulk(b *testing.B) {
	s := New(WithLanguage("bn"))
	inputs := []string{"বাংলা", "প্রিয়", "ক্ষমা"}
	for i := 0; i < b.N; i++ {
		s.MakeBulk(context.Background(), inputs)
	}
}

func BenchmarkMakeZeroAlloc(b *testing.B) {
	s := New(WithZeroAlloc(true))
	input := "বাংলা প্রিয়"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s.Make(context.Background(), input)
	}
}

func BenchmarkMakeLegacy(b *testing.B) {
	s := New(WithZeroAlloc(false))
	input := "বাংলা প্রিয়"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s.Make(context.Background(), input)
	}
}
