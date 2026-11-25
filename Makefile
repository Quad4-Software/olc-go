# Dev targets for olc-go. The tools target installs pinned versions into $(GOBIN).
GO ?= go
GOBIN ?= $(shell $(GO) env GOPATH)/bin
GOLANGCI_LINT ?= $(GOBIN)/golangci-lint
GOFUMPT ?= $(GOBIN)/gofumpt
GOIMPORTS ?= $(GOBIN)/goimports
STATICCHECK ?= $(GOBIN)/staticcheck

GOLANGCI_LINT_VERSION ?= v2.15.0
STATICCHECK_VERSION ?= v0.8.2
GOFUMPT_VERSION ?= v0.12.0
GOIMPORTS_VERSION ?= v0.53.0
MODULE ?= git.quad4.io/libraries/olc-go

.PHONY: help tools fmt fmt-check vet lint staticcheck test test-race bench fuzz ci clean

help:
	@printf '%s\n' \
		'Targets:' \
		'  tools         Install or update golangci-lint, gofumpt, goimports, staticcheck' \
		'  fmt           Format with gofumpt and goimports' \
		'  fmt-check     Check gofumpt and goimports formatting' \
		'  vet           Run go vet' \
		'  lint          Run golangci-lint' \
		'  staticcheck   Run staticcheck' \
		'  test          Run go test ./...' \
		'  test-race     Run tests with the race detector' \
		'  bench         Bench smoke for EncodeTo' \
		'  fuzz          Short bounded fuzz' \
		'  ci            Run all local checks' \
		'  clean         Remove local build artifacts'

tools:
	$(GO) install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@$(GOLANGCI_LINT_VERSION)
	$(GO) install mvdan.cc/gofumpt@$(GOFUMPT_VERSION)
	$(GO) install golang.org/x/tools/cmd/goimports@$(GOIMPORTS_VERSION)
	$(GO) install honnef.co/go/tools/cmd/staticcheck@$(STATICCHECK_VERSION)

$(GOLANGCI_LINT):
	$(GO) install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@$(GOLANGCI_LINT_VERSION)

$(GOFUMPT):
	$(GO) install mvdan.cc/gofumpt@$(GOFUMPT_VERSION)

$(GOIMPORTS):
	$(GO) install golang.org/x/tools/cmd/goimports@$(GOIMPORTS_VERSION)

$(STATICCHECK):
	$(GO) install honnef.co/go/tools/cmd/staticcheck@$(STATICCHECK_VERSION)

fmt: $(GOFUMPT) $(GOIMPORTS)
	$(GOFUMPT) -l -w .
	$(GOIMPORTS) -local $(MODULE) -w .

fmt-check: $(GOFUMPT) $(GOIMPORTS)
	@files="$$($(GOFUMPT) -l .)"; test -z "$$files" || { echo "$$files"; exit 1; }
	@files="$$($(GOIMPORTS) -local $(MODULE) -l .)"; test -z "$$files" || { echo "$$files"; exit 1; }

vet:
	$(GO) vet ./...

lint: $(GOLANGCI_LINT)
	$(GOLANGCI_LINT) run ./...

staticcheck: $(STATICCHECK)
	$(STATICCHECK) ./...

test:
	$(GO) test ./... -count=1

test-race:
	$(GO) test ./... -race -count=1

bench:
	$(GO) test ./olc -bench=BenchmarkEncodeTo_preallocated -benchmem -count=1 -benchtime=100ms

fuzz:
	$(GO) test ./olc -fuzz=FuzzEncodeDecodeRoundTrip -fuzztime=5s
	$(GO) test ./olc -fuzz=FuzzCodeOperations -fuzztime=5s

ci: fmt-check vet lint staticcheck test test-race bench fuzz

clean:
	rm -f coverage.out cover.out cpu.out mem.out olc.test example/example
