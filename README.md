# passkeypier

> A harbor for passkeys. Every credential that sails in is checked at the
> lighthouse before it is allowed to dock.

<p align="center">
  <img src="docs/assets/credential-lighthouse.svg" alt="A harbor lighthouse acting as a verification checkpoint, with a credential passing five gates before docking" width="640" />
</p>

**passkeypier** is a mixed-language Passkey and WebAuthn conformance lab. The Go
core and CLI run virtual registration and authentication ceremonies: secure
random challenges, `base64url`, origin and relying-party checks, signature
counters, user-verification policies, and Ed25519 signing and verification. A
small scenario engine drives those ceremonies and emits machine-readable
conformance reports. A dependency-light TypeScript browser lab reads those
reports and renders them, entirely offline.

The purpose is to make the logic of a passkey ceremony executable and testable.
It is not FIDO-certified and makes no certification claims. It deliberately
leaves out attestation-statement trust, CBOR and COSE_Key parsing, and
transport bindings. The exact scope and the honest list of non-goals live in
[`docs/SPEC.md`](docs/SPEC.md).

---

## Chart of the harbor

- [Why a pier](#why-a-pier)
- [What sails in and out](#what-sails-in-and-out)
- [Quick departure](#quick-departure)
- [The two ceremonies at the dock](#the-two-ceremonies-at-the-dock)
- [Command transcripts](#command-transcripts)
- [Conformance scenarios](#conformance-scenarios)
- [The report format](#the-report-format)
- [The browser lab](#the-browser-lab)
- [Layout of the harbor](#layout-of-the-harbor)
- [Building, testing, and CI](#building-testing-and-ci)
- [Design notes and honest limits](#design-notes-and-honest-limits)
- [License](#license)

---
