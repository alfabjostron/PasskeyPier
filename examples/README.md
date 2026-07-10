# Examples

This directory holds runnable examples for the passkeypier conformance lab.

## `sample-report.json`

A conformance report produced by the Go CLI. Regenerate it with:

```sh
go run ./cmd/passkeypier run -format json -out examples/sample-report.json
```

The TypeScript browser lab (`web/`) loads this file to render the report.

## `demo_test.go`

An example-style Go test that documents the honest register + authenticate
ceremony as a self-contained, runnable transcript. Run it with:

```sh
go test ./examples/ -run Example -v
```
