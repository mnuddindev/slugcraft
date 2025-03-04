package slugcraft

import (
	"context"
	"testing"
)

// TestNew tests the default Config configuration
func TestNew(t *testing.T) {
	s := New()
	if len(s.PipeLine) != 2 {
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
		{"ru", "Привет, мир!", "privet-mir"},
		{"zh", "你好世界", "nhsj"},
		{"hi", "नमस्ते दुनिया", "nmste-dunya"},
		{"ja", "こんにちは世界", "knnicha-sekai"},
		{"de", "Hallo Welt!", "hallo-welt"},
		{"es", "Café Olé!", "cafe-ole"},
	}

	for _, tt := range tests {
		s := New(WithLanguage(tt.language))
		slug, err := s.Make(context.Background(), tt.input)
		if err != nil {
			t.Errorf("Make(%q, lang=%q) returned error: %v", tt.input, tt.language, err)
		}
		if slug != tt.expected {
			t.Errorf("Make(%q, lang=%q) = %q, expected %q", tt.input, tt.language, slug, tt.expected)
		}
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
		slug, err := s.Make(context.Background(), tt.input)
		if err != nil {
			t.Errorf("Make(%q) returned error: %v", tt.input, err)
		}
		if slug != tt.expected {
			t.Errorf("Make(%q) = %q, expected %q", tt.input, slug, tt.expected)
		}
	}
}

// TestMakeWithCache tests collision avoidance with in-memory cache.
func TestMakeWithCache(t *testing.T) {
	s := New(
		WithUseCache(true),
		WithSuffixStyle("numeric"),
	)
	tests := []struct {
		input    string
		expected string
	}{
		{"My Post", "my-post"},
		{"My Post", "my-post-1"}, // Collision
		{"My Post", "my-post-2"}, // Another collision
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
	c := &cache{store: make(map[string]struct{})}

	// Test set and get
	c.set("slug1")
	if !c.get("slug1") {
		t.Errorf("cache.get('slug1') = false, expected true")
	}
	if c.get("slug2") {
		t.Errorf("cache.get('slug2') = true, expected false")
	}

	// Test delete
	c.set("slug2")
	c.delete("slug2")
	if c.get("slug2") {
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
