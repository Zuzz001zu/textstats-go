// Package stats computes simple statistics for a piece of text:
// line, word and character counts, and the most frequent words.
package stats

import (
	"sort"
	"strings"
	"unicode"
	"unicode/utf8"
)

// Summary holds the basic counts for a text.
type Summary struct {
	Lines int `json:"lines"`
	Words int `json:"words"`
	Chars int `json:"chars"`
}

// WordCount is a word and how many times it appears.
type WordCount struct {
	Word  string `json:"word"`
	Count int    `json:"count"`
}

// Analyze returns the line, word and character counts of text.
// Characters are counted as Unicode runes, excluding line breaks.
func Analyze(text string) Summary {
	return Summary{
		Lines: countLines(text),
		Words: len(Words(text)),
		Chars: utf8.RuneCountInString(strings.ReplaceAll(text, "\n", "")),
	}
}

// Words splits text into lowercase words. A word is a run of
// letters or digits, so punctuation and spaces act as separators.
func Words(text string) []string {
	fields := strings.FieldsFunc(text, func(r rune) bool {
		return !unicode.IsLetter(r) && !unicode.IsDigit(r)
	})
	for i, f := range fields {
		fields[i] = strings.ToLower(f)
	}
	return fields
}

// Top returns the n most frequent words, most frequent first.
// Words with the same count are sorted alphabetically so the
// result is always the same for the same input.
func Top(text string, n int) []WordCount {
	if n <= 0 {
		return []WordCount{}
	}
	counts := make(map[string]int)
	for _, w := range Words(text) {
		counts[w]++
	}

	result := make([]WordCount, 0, len(counts))
	for w, c := range counts {
		result = append(result, WordCount{Word: w, Count: c})
	}
	sort.Slice(result, func(i, j int) bool {
		if result[i].Count != result[j].Count {
			return result[i].Count > result[j].Count
		}
		return result[i].Word < result[j].Word
	})

	if len(result) > n {
		result = result[:n]
	}
	return result
}

// countLines counts lines the way a text editor shows them:
// a trailing newline does not start an extra empty line.
func countLines(text string) int {
	if text == "" {
		return 0
	}
	n := strings.Count(text, "\n")
	if !strings.HasSuffix(text, "\n") {
		n++
	}
	return n
}
