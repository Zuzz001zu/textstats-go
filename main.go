// Command textstats prints statistics about a text file.
//
// Usage:
//
//	textstats count <file> [--json]
//	textstats top [-n 5] <file> [--json]
//	textstats version
//
// Use "-" as the file name to read from standard input.
package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/Zuzz001zu/textstats-go/internal/stats"
)

const version = "1.1.0"

const usage = `textstats - text statistics from the command line

Usage:
  textstats count <file> [--json]       count lines, words and characters
  textstats top [-n N] <file> [--json]  show the N most frequent words (default 5)
  textstats version                     print the version

Use "-" as <file> to read from standard input.
`

func main() {
	os.Exit(run(os.Args[1:], os.Stdin, os.Stdout, os.Stderr))
}

// run contains the whole program logic. It returns the exit code
// instead of calling os.Exit, so it can be tested.
func run(args []string, stdin io.Reader, stdout, stderr io.Writer) int {
	if len(args) == 0 {
		fmt.Fprint(stderr, usage)
		return 2
	}

	switch args[0] {
	case "count":
		return runCount(args[1:], stdin, stdout, stderr)
	case "top":
		return runTop(args[1:], stdin, stdout, stderr)
	case "version":
		fmt.Fprintln(stdout, "textstats", version)
		return 0
	case "help", "-h", "--help":
		fmt.Fprint(stdout, usage)
		return 0
	default:
		fmt.Fprintf(stderr, "error: unknown command %q\n\n%s", args[0], usage)
		return 2
	}
}

func runCount(args []string, stdin io.Reader, stdout, stderr io.Writer) int {
	fs := flag.NewFlagSet("count", flag.ContinueOnError)
	fs.SetOutput(stderr)
	asJSON := fs.Bool("json", false, "print the result as JSON")
	text, code := parseAndRead(fs, args, stdin, stderr)
	if code != 0 {
		return code
	}

	s := stats.Analyze(text)
	if *asJSON {
		return printJSON(stdout, stderr, s)
	}
	fmt.Fprintf(stdout, "Lines: %d\nWords: %d\nChars: %d\n", s.Lines, s.Words, s.Chars)
	return 0
}

func runTop(args []string, stdin io.Reader, stdout, stderr io.Writer) int {
	fs := flag.NewFlagSet("top", flag.ContinueOnError)
	fs.SetOutput(stderr)
	n := fs.Int("n", 5, "number of words to show")
	asJSON := fs.Bool("json", false, "print the result as JSON")
	text, code := parseAndRead(fs, args, stdin, stderr)
	if code != 0 {
		return code
	}
	if *n < 1 {
		fmt.Fprintln(stderr, "error: -n must be at least 1")
		return 2
	}

	top := stats.Top(text, *n)
	if *asJSON {
		return printJSON(stdout, stderr, top)
	}
	for _, wc := range top {
		fmt.Fprintf(stdout, "%-15s %d\n", wc.Word, wc.Count)
	}
	return 0
}

// parseAndRead parses the flags and reads the file named by the
// single positional argument. Flags may come before or after the file.
func parseAndRead(fs *flag.FlagSet, args []string, stdin io.Reader, stderr io.Writer) (string, int) {
	var files []string
	for {
		if err := fs.Parse(args); err != nil {
			return "", 2
		}
		if fs.NArg() == 0 {
			break
		}
		files = append(files, fs.Arg(0))
		args = fs.Args()[1:]
	}
	if len(files) != 1 {
		fmt.Fprintf(stderr, "error: expected exactly one file, got %d\n", len(files))
		return "", 2
	}

	text, err := readInput(files[0], stdin)
	if err != nil {
		fmt.Fprintln(stderr, "error:", err)
		return "", 1
	}
	return text, 0
}

func readInput(name string, stdin io.Reader) (string, error) {
	var data []byte
	var err error
	if name == "-" {
		data, err = io.ReadAll(stdin)
	} else {
		data, err = os.ReadFile(name)
	}
	if errors.Is(err, os.ErrNotExist) {
		return "", fmt.Errorf("file not found: %s", name)
	}
	return string(data), err
}

func printJSON(stdout, stderr io.Writer, v any) int {
	enc := json.NewEncoder(stdout)
	enc.SetIndent("", "  ")
	if err := enc.Encode(v); err != nil {
		fmt.Fprintln(stderr, "error:", err)
		return 1
	}
	return 0
}
