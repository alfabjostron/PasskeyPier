# Changelog

All notable changes to passkeypier are documented here. The format follows
[Keep a Changelog](https://keepachangelog.com/en/1.1.0/) and this project
adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [0.1.0] - 2026-01-01

### Added

- Go standard-library core (`internal/harbor`) modeling passkey ceremonies:
  - Cryptographically secure challenge generation (`crypto/rand`), 16-byte
    spec minimum with a 32-byte default.
  - Unpadded base64url encode/decode helpers.
  - Client data (`webauthn.create` / `webauthn.get`) marshaling and SHA-256
    hashing.
  - Authenticator data marshaling/parsing with UP/UV/BE/BS/AT/ED flag handling
    and a big-endian signature counter.
  - Virtual authenticator holding Ed25519 discoverable credentials with a
    per-credential counter and configurable user-verification capability.
  - Relying-party verification: origin, ceremony type, challenge equality,
    RP ID hash, user-verification policy, user-presence, Ed25519 assertion
    signature verification, and signature-counter monotonicity.
- Built-in conformance scenario suite (registration, authentication, policy and
  security categories) with positive and negative expectations.
- Conformance report engine emitting JSON (schema `passkeypier/report/v1`) and
  human-readable text.
- `passkeypier` CLI with `run`, `demo`, `list` and `version` subcommands.
- Focused unit tests and runnable Go examples.
- Dependency-light TypeScript browser lab (`web/`) that validates and renders
  reports, fully offline (no remote assets, embedded sample report).
- Documentation: themed README, `docs/SPEC.md`, two original animated SVGs.
- Makefile, MIT license, `.gitignore`, and GitHub Actions CI (Go + TypeScript).

### Notes

- This is an educational conformance-exploration tool. It is **not**
  FIDO-certified and deliberately omits attestation-statement verification,
  CBOR/COSE_Key parsing, and transport bindings. See `docs/SPEC.md` for the
  precise scope and non-goals.

[0.1.0]: https://example.com/passkeypier/releases/tag/v0.1.0

// draft note 1
