# passkeypier — Specification

Version: 0.1.0 · Report schema: `passkeypier/report/v1`

This document defines exactly what passkeypier models, how the ceremonies are
computed, and — just as importantly — what it deliberately does **not** do.

## 1. Scope and non-goals

passkeypier is a virtual conformance lab for exploring the data-model checks a
WebAuthn [relying party](https://www.w3.org/TR/webauthn-2/#relying-party)
performs during passkey registration and authentication. It is designed to run
deterministically in a test harness with no browser, no authenticator hardware,
and no network.

### In scope

- Secure random challenge generation and lifecycle.
- `base64url` (unpadded) encoding of binary values.
- Client data (`webauthn.create` / `webauthn.get`) construction and SHA-256
  hashing.
- Authenticator data layout, flag bits, and the big-endian signature counter.
- Ed25519 (COSE algorithm `-8`, OKP curve Ed25519) signing and verification.
- Relying-party verification steps: type, challenge, origin, RP ID hash,
  user-verification policy, user presence, signature, and counter monotonicity.
- User-verification policy semantics (`required` / `preferred` / `discouraged`).
- A scenario engine with positive/negative expectations and JSON/text reports.

### Explicit non-goals
