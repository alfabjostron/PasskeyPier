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

passkeypier is **not** a FIDO-certified implementation and makes no
certification claims. The following are intentionally out of scope:

- **Attestation statement verification.** No `packed`, `tpm`, `android-key`,
  `android-safetynet`, `fido-u2f`, or `apple` attestation formats are parsed or
  trusted. The lab treats registration attestation as `none`-equivalent.
- **CBOR / COSE_Key parsing.** The attested credential public key is modeled as
  the raw 32-byte Ed25519 key rather than a full COSE_Key CBOR map. Signature
  verification uses the stored key object directly.
- **Additional COSE algorithms.** Only Ed25519 is implemented. ES256 (`-7`) and
  RS256 (`-257`) are not.
- **Transport / CTAP wire framing, extensions, `largeBlob`, `hmac-secret`, and
  enterprise attestation.**
- **RP ID derivation rules** beyond an exact-match hash check (no registrable
  domain suffix walking).

If you need certified behavior, use a certified stack. passkeypier exists to
make the *logic* of the ceremony legible and testable.

## 2. Data model

### 2.1 Challenge

A challenge is `n` random bytes from `crypto/rand`. The WebAuthn-recommended
minimum of 16 bytes is enforced; the default is 32 bytes. On the wire the
challenge appears inside client data as unpadded base64url.

### 2.2 Client data

