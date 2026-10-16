package main

import (
	"bytes"
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
