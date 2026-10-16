package stats

import (
	"reflect"
	"testing"
)

func TestAnalyze(t *testing.T) {
	tests := []struct {
		name string
		text string
		want Summary
	}{
		{"empty", "", Summary{0, 0, 0}},
		{"one line", "Hello world", Summary{1, 2, 11}},
		{"trailing newline", "Hello\nworld\n", Summary{2, 2, 10}},
		{"unicode", "Привет мир\nпривет Git", Summary{2, 4, 20}},
		{"punctuation", "Go, Go... GO!", Summary{1, 3, 13}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Analyze(tt.text); got != tt.want {
				t.Errorf("Analyze(%q) = %+v, want %+v", tt.text, got, tt.want)
			}
		})
	}
}

func TestWords(t *testing.T) {
	got := Words("Git-based CLI, v2!")
	want := []string{"git", "based", "cli", "v2"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Words() = %v, want %v", got, want)
	}
}

func TestTop(t *testing.T) {
	text := "b a c a b a"
	got := Top(text, 2)
	want := []WordCount{{"a", 3}, {"b", 2}}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Top() = %v, want %v", got, want)
	}
}

func TestTopTiesAreSortedAlphabetically(t *testing.T) {
	got := Top("zeta alpha mid", 3)
	want := []WordCount{{"alpha", 1}, {"mid", 1}, {"zeta", 1}}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Top() = %v, want %v", got, want)
	}
}

func TestTopEdgeCases(t *testing.T) {
	if got := Top("a b", 0); len(got) != 0 {
		t.Errorf("Top(n=0) = %v, want empty", got)
	}
	if got := Top("", 5); len(got) != 0 {
		t.Errorf("Top(empty text) = %v, want empty", got)
	}
	if got := Top("a b", 10); len(got) != 2 {
		t.Errorf("Top(n > words) returned %d items, want 2", len(got))
	}
}
