package slugcraft

import "testing"

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
