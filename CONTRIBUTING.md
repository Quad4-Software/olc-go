# Contributing

## Development

Requires Go 1.27.4 or later. No external dependencies.

Common tasks:

| Task | Command |
|------|---------|
| Test | `go test ./...` |
| Race | `go test ./... -race` |
| Lint | `make lint` |
| Format | `make fmt` |
| Local CI | `make ci` |

Install or update the lint tools with `make tools`.

## Guidelines

- Run `make fmt` before committing. `make ci` checks gofumpt and goimports.
- Keep the library dependency-free.
- Tests must pass against the official OLC CSV corpus in `olc/testdata/`.
- Preserve the zero-allocation guarantees covered by `alloc_test.go`.

## Reporting issues

Open an issue at <https://git.quad4.io/libraries/olc-go/issues> with a failing
code, coordinate pair, or corpus row when reporting encode/decode mismatches.
