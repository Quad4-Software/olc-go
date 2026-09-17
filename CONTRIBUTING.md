# Contributing

## Development

Requires Go 1.27.1 or later. No external dependencies.

Common tasks:

| Task | Command |
|------|---------|
| Test | `go test ./...` |
| Race | `go test ./... -race` |
| Lint | `make lint` |
| Format | `make fmt` |
| Local CI | `make ci` |

Install the lint tools once with `make tools`.

## Guidelines

- Run `make fmt` before committing. CI enforces gofumpt and goimports.
- Keep the library dependency-free.
- Tests must pass against the official OLC CSV corpus in `olc/testdata/`.
- Preserve the zero-allocation guarantees covered by `alloc_test.go`.

## Reporting issues

Open an issue on GitHub with a failing code, coordinate pair, or corpus row
when reporting encode/decode mismatches.
