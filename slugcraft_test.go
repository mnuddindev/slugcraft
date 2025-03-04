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
