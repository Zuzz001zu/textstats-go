# textstats

[![CI](https://github.com/Zuzz001zu/textstats-go/actions/workflows/ci.yml/badge.svg)](https://github.com/Zuzz001zu/textstats-go/actions/workflows/ci.yml)

A small command-line tool written in Go that prints statistics about a text file:
line, word and character counts, and the most frequent words.

Built with the Go standard library only — no external dependencies.

## Features

- `count` - lines, words and characters (Unicode-aware)
- `top` - the N most frequent words, with stable ordering
- `--json` output for scripts and automation
- Reads from a file or from standard input (`-`)
- Clear error messages and exit codes (`0` OK, `1` runtime error, `2` usage error)
- Unit and CLI tests, CI on Linux, Windows and macOS

## Installation

Requires Go 1.22 or newer.

```bash
go install github.com/Zuzz001zu/textstats-go/cmd/textstats@latest
```

Or build from source:

```bash
git clone https://github.com/Zuzz001zu/textstats-go.git
cd textstats-go
go build -o textstats ./cmd/textstats
```

## Usage

```bash
textstats count testdata/sample.txt
```

```
Lines: 4
Words: 36
Chars: 200
```

```bash
textstats top -n 3 testdata/sample.txt
```

```
a               4
is              3
for             2
```

JSON output:

```bash
textstats count testdata/sample.txt --json
```

```json
{
  "lines": 4,
  "words": 36,
  "chars": 200
}
```

Read from standard input:

```bash
echo "go go cli" | textstats top -n 1 -
```

## Project structure

```
textstats-go/
├── cmd/textstats/        # CLI: commands, flags, output
│   ├── main.go
│   └── main_test.go
├── internal/stats/       # text statistics logic (no I/O)
│   ├── stats.go
│   └── stats_test.go
├── testdata/sample.txt   # example input
├── .github/workflows/    # GitHub Actions CI
└── go.mod
```

The logic lives in `internal/stats` and has no input/output, so it is easy to test.
The CLI in `cmd/textstats` only parses arguments, reads input and formats output.
`run()` returns an exit code instead of calling `os.Exit`, which lets the tests
check the full CLI behavior.

## Development

```bash
go test -cover ./...   # run tests with coverage
go vet ./...           # static checks
gofmt -l .             # list unformatted files
```

## License

MIT — see [LICENSE](LICENSE).
