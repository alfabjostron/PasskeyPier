# Contributing to PasskeyPier

PasskeyPier is deliberately scoped: run passkey ceremonies locally against
a matrix of relying-party scenarios, inspect the bytes, render
deterministic JSON reports. No network during ceremonies, no dependencies
in the core. Contributions that respect that scope are welcome.

## Development setup

```bash
git clone https://github.com/alfabjostron/PasskeyPier.git
cd PasskeyPier
go build ./...
go test -cover ./...
```

The browser lab lives in `web/` (`npm install && npx tsc`). CI verifies
`go mod tidy` produces no diff, and the smoke run must emit a non-empty
report.

## Ground rules

- **Deterministic reports.** The same scenario matrix must always produce
  the same report bytes. No wall-clock, no randomness, no map-iteration
  order in output.
- **Std-only core.** The harbor package (challenge bytes, base64url,
  ceremony, relying-party matrix) stays on the Go standard library.
- **Byte fidelity.** Challenge/credential byte handling is the whole
  point; changes there need boundary coverage in
  `internal/harbor/harbor_test.go`.
- **Tests on behaviour changes.** Scenario and attestation-format changes
  need coverage across the full matrix.

## Commit style

Short imperative subjects (`feat: ...`, `fix: ...`, `docs: ...`). Body only
when the "why" is not obvious from the diff.

## Reporting issues

Attach the scenario matrix (sanitized), the exact CLI flags, and the
report section that looks wrong.