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
