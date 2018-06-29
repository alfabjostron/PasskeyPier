# passkeypier — build, test and report targets.
#
# Go core/CLI uses only the standard library. The TypeScript browser lab needs
# a local TypeScript compiler (npm i inside web/, or a global `tsc`).

GO      ?= go
PKG     := ./...
BIN     := passkeypier
REPORT  := report.json

.PHONY: all
all: build test ## build and test everything (default)

.PHONY: build
build: ## compile the Go CLI
	$(GO) build -o $(BIN) ./cmd/passkeypier

.PHONY: test
test: ## run Go unit tests and examples
