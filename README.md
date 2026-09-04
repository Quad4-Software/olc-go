# OLC-Go

[![CI](https://github.com/Quad4-Software/olc-go/actions/workflows/ci.yml/badge.svg)](https://github.com/Quad4-Software/olc-go/actions/workflows/ci.yml)
[![CodeQL](https://github.com/Quad4-Software/olc-go/actions/workflows/codeql.yml/badge.svg)](https://github.com/Quad4-Software/olc-go/actions/workflows/codeql.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/Quad4-Software/olc-go.svg)](https://pkg.go.dev/github.com/Quad4-Software/olc-go)
[![License](https://img.shields.io/github/license/Quad4-Software/olc-go)](LICENSE)

Dependency-free Go implementation of [Open Location Code](https://github.com/google/open-location-code) (OLC).

**Module:** `github.com/Quad4-Software/olc-go`

## Requirements

- Go 1.26.5 or later

## Install

```bash
go get github.com/Quad4-Software/olc-go@latest
```

## Usage

```go
package main

import (
	"fmt"

	"github.com/Quad4-Software/olc-go/olc"
)

func main() {
	code := olc.Encode(-1.286386, 36.817223, 10)
	fmt.Println(code) // 6GCRPR78+CV

	area, err := olc.Decode(code)
	if err != nil {
		panic(err)
	}
	lat, lng := area.Center()
	fmt.Println(lat, lng)

	short, err := olc.Shorten(code, -1.286386, 36.817223)
	if err != nil {
		panic(err)
	}
	full, err := olc.RecoverNearest(short, -1.286386, 36.817223)
	if err != nil {
		panic(err)
	}
	fmt.Println(short, full)

	// Zero-allocation encode into a reused buffer.
	var buf [olc.MaxEncodedLen]byte
	n, err := olc.EncodeTo(buf[:], -1.286386, 36.817223, 10)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(buf[:n]))
}
```

### Example CLI

```bash
go run ./example -lat -1.286386 -lon 36.817223
go run ./example -decode 6GCRPR78+CV
go run ./example -lat -1.286386 -lon 36.817223 -shorten
go run ./example -recover '+CV' -ref-lat -1.286386 -ref-lon 36.817223
```

## Validation

Tests load the official Open Location Code CSV corpus from
[google/open-location-code](https://github.com/google/open-location-code) (`olc/testdata/`):

| Corpus | Checks |
|--------|--------|
| `validityTests.csv` | `IsValid` / `IsShort` / `IsFull` |
| `encoding.csv` | integer mapping + `Encode` exact strings |
| `decoding.csv` | `Decode` bounding boxes |
| `shortCodeTests.csv` | `Shorten` / `RecoverNearest` |

`EncodeTo` and `Decode` are covered by zero-allocation tests. Fuzz coverage lives in `FuzzEncodeDecodeRoundTrip`.

## Performance

Prefer `EncodeTo` or `AppendEncode` when encoding in a hot path. `Encode` returns a `string` and therefore allocates once. `Decode` of a full code is allocation-free on the success path.

## Testing

| Goal | Command |
|------|---------|
| Default | `go test ./...` |
| Race | `go test ./... -race` |
| Lint | `make lint` |
| Bench | `go test ./olc -bench=. -benchmem` |
| Local CI | `make ci` |

## License

Copyright 2026 Quad4. See [LICENSE](LICENSE) (0BSD).

Official OLC test CSV files are Apache-2.0; see [NOTICE](NOTICE).
