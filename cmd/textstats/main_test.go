package main

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
)

// runCLI runs the program with args and the given stdin text.
func runCLI(t *testing.T, stdin string, args ...string) (stdout, stderr string, code int) {
	t.Helper()
	var out, errOut bytes.Buffer
	code = run(args, strings.NewReader(stdin), &out, &errOut)
	return out.String(), errOut.String(), code
}

func TestCountText(t *testing.T) {
	out, _, code := runCLI(t, "hello world\nhello go\n", "count", "-")
	if code != 0 {
		t.Fatalf("exit code = %d, want 0", code)
	}
	want := "Lines: 2\nWords: 4\nChars: 19\n"
	if out != want {
		t.Errorf("output = %q, want %q", out, want)
	}
}

func TestCountJSON(t *testing.T) {
	out, _, code := runCLI(t, "a b c", "count", "-", "--json")
	if code != 0 {
		t.Fatalf("exit code = %d, want 0", code)
	}
	var got struct{ Lines, Words, Chars int }
	if err := json.Unmarshal([]byte(out), &got); err != nil {
		t.Fatalf("output is not valid JSON: %v\n%s", err, out)
	}
	if got.Lines != 1 || got.Words != 3 || got.Chars != 5 {
		t.Errorf("got %+v", got)
	}
}

func TestTopJSON(t *testing.T) {
	out, _, code := runCLI(t, "go go cli", "top", "--json", "-n", "1", "-")
	if code != 0 {
		t.Fatalf("exit code = %d, want 0", code)
	}
	var got []struct {
		Word  string
		Count int
	}
	if err := json.Unmarshal([]byte(out), &got); err != nil {
		t.Fatalf("output is not valid JSON: %v\n%s", err, out)
	}
	if len(got) != 1 || got[0].Word != "go" || got[0].Count != 2 {
		t.Errorf("got %+v", got)
	}
}

func TestErrors(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		wantCode int
		wantErr  string
	}{
		{"no command", nil, 2, "Usage"},
		{"unknown command", []string{"jump"}, 2, "unknown command"},
		{"missing file", []string{"count"}, 2, "expected exactly one file"},
		{"file not found", []string{"count", "nope.txt"}, 1, "file not found"},
		{"bad -n", []string{"top", "-n", "0", "-"}, 2, "-n must be at least 1"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, errOut, code := runCLI(t, "", tt.args...)
			if code != tt.wantCode {
				t.Errorf("exit code = %d, want %d", code, tt.wantCode)
			}
			if !strings.Contains(errOut, tt.wantErr) {
				t.Errorf("stderr = %q, want it to contain %q", errOut, tt.wantErr)
			}
		})
	}
}

func TestVersion(t *testing.T) {
	out, _, code := runCLI(t, "", "version")
	if code != 0 || !strings.HasPrefix(out, "textstats ") {
		t.Errorf("version output = %q, code %d", out, code)
	}
}
