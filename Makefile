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
	$(GO) test $(PKG)

.PHONY: cover
cover: ## run tests with coverage summary
	$(GO) test -cover $(PKG)

.PHONY: vet
vet: ## run go vet
	$(GO) vet $(PKG)

.PHONY: run
run: ## run the conformance suite (text)
	$(GO) run ./cmd/passkeypier run

.PHONY: report
report: ## write a JSON conformance report to report.json
	$(GO) run ./cmd/passkeypier run -format json -out $(REPORT)

.PHONY: demo
demo: ## run one register+authenticate ceremony
	$(GO) run ./cmd/passkeypier demo

.PHONY: web-sample
web-sample: ## generate the sample report consumed by the web lab
	$(GO) run ./cmd/passkeypier run -format json -out web/sample-report.json
	$(GO) run ./cmd/passkeypier run -format json -out examples/sample-report.json

.PHONY: web-typecheck
web-typecheck: ## typecheck the TypeScript lab (needs tsc)
	cd web && tsc --noEmit

.PHONY: web-build
web-build: ## compile the TypeScript lab to web/dist
	cd web && tsc

.PHONY: clean
clean: ## remove build artifacts
	-rm -f $(BIN) $(BIN).exe $(REPORT)
	-rm -rf web/dist

.PHONY: help
help: ## list targets
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "  %-16s %s\n", $$1, $$2}'
