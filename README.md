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

## Why a pier

A passkey ceremony is a small piece of maritime choreography. The relying party
(the harbor) issues a one-time challenge, a bottle thrown into the water. The
authenticator (a ship carrying a private Ed25519 key) writes its reply, seals it
with a signature, and sends it back. The harbor's lighthouse checks every
returning vessel. Is this the right origin? The challenge I actually issued?
Bound to my relying-party identifier? Did the signature counter advance, or is
this a cloned ship? Only when every check passes may the credential dock.

passkeypier makes that choreography runnable and testable, with no browser or
hardware key in sight.

---

## What sails in and out

Implemented against the Go standard library only (`crypto/ed25519`,
`crypto/rand`, `crypto/sha256`, `encoding/base64`, `encoding/json`,
`encoding/binary`):

- **Secure challenges** from `crypto/rand`, with the 16-byte spec minimum
  enforced and a 32-byte default.
- **`base64url`** encode and decode without padding, with lenient decoding for
  hand-authored fixtures.
- **Client data** for `webauthn.create` and `webauthn.get`: canonical JSON,
  SHA-256 `clientDataHash`, strict decoding that rejects unknown fields.
- **Authenticator data** in the exact wire layout
  `rpIdHash(32) then flags(1) then signCount(4, big-endian) then trailing`,
  with UP, UV, BE, BS, AT, and ED flags.
- **Virtual authenticator** holding resident (discoverable) Ed25519 credentials,
  each with a signature counter and a configurable user-verification capability.
- **Relying-party verification**: origin, ceremony type, challenge equality, RP
  ID hash, user-verification policy, user presence, the Ed25519 assertion
  signature, and signature-counter monotonicity.
- **Scenario engine** covering positive and negative expectations across the
  registration, authentication, policy, and security categories.
- **Reports** in JSON (schema `passkeypier/report/v1`) and human-readable text.

The TypeScript lab in `web/` validates the JSON report against the schema before
rendering, then shows a summary, per-category totals, and an expandable list of
scenarios. It carries no runtime dependencies and makes no network access.

---

## Quick departure

You need **Go 1.24+**. The web lab additionally needs a TypeScript compiler
(`tsc`) if you want to typecheck or build it.

```sh
# Run the full conformance suite (human-readable)
go run ./cmd/passkeypier run

# Emit a JSON report for the browser lab
go run ./cmd/passkeypier run -format json -out report.json

# Perform a single honest register plus authenticate ceremony
go run ./cmd/passkeypier demo

# List every scenario and what it asserts
go run ./cmd/passkeypier list
```

Or use the Makefile:

```sh
make            # build + test
make run        # text conformance report
make report     # writes report.json
make demo       # one ceremony
make help       # list all targets
```

---

## The two ceremonies at the dock

<p align="center">
  <img src="docs/assets/ceremony-dock.svg" alt="Data flow between relying party and authenticator: a challenge goes out, authenticator data and an Ed25519 signature come back" width="680" />
</p>

**Registration (`webauthn.create`).** The harbor issues creation options with a
fresh challenge and a user-verification policy. The ship mints a new Ed25519
credential bound to the relying-party ID, builds client data, and returns
authenticator data with the `AT` (attested-credential-data) flag set. The harbor
checks type, challenge, and origin, the RP ID hash, the UV policy, and user
presence, then stores the public key, sign count, user handle, and backup
eligibility.

**Authentication (`webauthn.get`).** The harbor issues request options with a
fresh challenge. The ship selects a resident credential, increments its counter,
and signs the concatenation `authenticatorData then SHA-256(clientData)` with
Ed25519. The harbor verifies everything registration checks, plus the signature
against the stored public key and the counter's monotonic advance.

The signed message is identical in structure to a real WebAuthn assertion:

```
message   = authenticatorData || SHA-256(clientDataJSON)
signature = Ed25519_Sign(credentialPrivateKey, message)
```

---

## Command transcripts

These are real outputs from the tool. Random values (credential ids, keys,
signatures, timestamps) will differ per run.

### `passkeypier demo -uv required`

```text
$ go run ./cmd/passkeypier demo -uv required
registration OK
  credential id: wxr0Wjqp7DvSVVJ9CSf4bg
  public key:    ZtSsH6A-tfSIGUR8esyG4cjJgyVt2u5CbfoLmT6nTlg
  user verified: true
authentication OK
  credential id: wxr0Wjqp7DvSVVJ9CSf4bg
  signature:     3bH_oSthdboI8ijkQafBh3uE_2BpgF0UZ48VykA_PnDsmePMzSmukiSm9_CHcMJ_Z7qARP9XUo5jMn6JSskiDA
  user verified: true
  user handle:   mariner-1
```

Note that the credential id in the assertion matches the one from registration,
and the base64url signature is 86 characters (64 raw Ed25519 bytes).

### `passkeypier run`

```text
$ go run ./cmd/passkeypier run
passkeypier conformance report
==============================
generated: 2026-08-31T17:54:56Z
schema:    passkeypier/report/v1

[PASS] authenticate/cloned-counter-regression     (security, expect reject)
       rejected as expected: harbor: signature counter did not increase (stored=99, got=1): possible cloned authenticator
[PASS] authenticate/happy-path                    (authentication, expect accept)
       ceremony accepted as expected
[PASS] authenticate/replayed-challenge            (authentication, expect reject)
       rejected as expected: harbor: challenge mismatch (possible replay or wrong ceremony)
[PASS] authenticate/tampered-signature            (security, expect reject)
       rejected as expected: harbor: assertion signature verification failed
[PASS] authenticate/uv-required-satisfied         (policy, expect accept)
       ceremony accepted as expected
[PASS] authenticate/wrong-origin                  (authentication, expect reject)
       rejected as expected: harbor: origin "https://evil.example", want "https://harbor.example"
[PASS] authenticate/wrong-rp-binding              (security, expect reject)
       rejected as expected: harbor: RP ID hash mismatch
[PASS] register/happy-path-uv-preferred           (registration, expect accept)
       ceremony accepted as expected
[PASS] register/uv-required-without-uv-support    (registration, expect reject)
       rejected as expected: harbor: authenticator does not support user verification

by category:
  authentication pass=3 fail=0
  policy         pass=1 fail=0
  registration   pass=2 fail=0
  security       pass=3 fail=0

summary: 9 passed, 0 failed of 9 total
result: ALL SCENARIOS PASSED
```

The process exits `0` when every scenario meets its expectation and `1` when any
conformance failure is present, which is convenient for CI gating.

### `passkeypier list`

```text
$ go run ./cmd/passkeypier list
register/happy-path-uv-preferred             [registration, expect accept]
    Honest registration with a UV-capable authenticator and preferred policy.
register/uv-required-without-uv-support      [registration, expect reject]
    UV=required against an authenticator that cannot verify the user must be rejected.
authenticate/happy-path                      [authentication, expect accept]
    Register then authenticate honestly; signature, origin and counter all valid.
authenticate/wrong-origin                    [authentication, expect reject]
    An assertion whose client data reports a foreign origin must be rejected.
authenticate/replayed-challenge              [authentication, expect reject]
    Reusing a stale challenge from a prior ceremony must fail the challenge check.
authenticate/uv-required-satisfied           [policy, expect accept]
    UV=required with a UV-capable authenticator sets the UV flag and passes.
authenticate/cloned-counter-regression       [security, expect reject]
    A signature counter that fails to advance signals a cloned authenticator and must be rejected.
authenticate/tampered-signature              [security, expect reject]
    Flipping a byte of the assertion signature must fail Ed25519 verification.
authenticate/wrong-rp-binding                [security, expect reject]
    Authenticator data carrying a foreign RP ID hash must fail the RP ID hash check.
```

---

## Conformance scenarios

Each scenario is a self-contained check with an expectation: either the harbor
should accept the ceremony, or it should reject it. A negative scenario that is
accidentally accepted counts as a conformance failure, which is how the lab
exercises the security checks.

| Scenario | Category | Expect | What it proves |
| --- | --- | --- | --- |
| `register/happy-path-uv-preferred` | registration | accept | Honest create ceremony with a UV-capable authenticator. |
| `register/uv-required-without-uv-support` | registration | reject | `UV=required` fails when the authenticator cannot verify the user. |
| `authenticate/happy-path` | authentication | accept | Signature, origin, and counter all valid end to end. |
| `authenticate/wrong-origin` | authentication | reject | A foreign origin in client data is caught. |
| `authenticate/replayed-challenge` | authentication | reject | A stale or mismatched challenge fails the equality check. |
| `authenticate/uv-required-satisfied` | policy | accept | `UV=required` sets and enforces the UV flag. |
| `authenticate/cloned-counter-regression` | security | reject | A non-advancing counter signals a cloned authenticator. |
| `authenticate/tampered-signature` | security | reject | A forged or flipped signature fails Ed25519 verification. |
| `authenticate/wrong-rp-binding` | security | reject | Authenticator data with a foreign RP ID hash is rejected. |

---

## The report format

Reports carry the schema tag `passkeypier/report/v1`. The JSON shape:

```json
{
  "schema": "passkeypier/report/v1",
  "tool": "passkeypier",
  "generated_at": "2026-08-31T17:59:21Z",
  "summary": { "total": 9, "passed": 9, "failed": 0, "all_passed": true },
  "categories": [
    { "category": "authentication", "passed": 3, "failed": 0 },
    { "category": "policy", "passed": 1, "failed": 0 },
    { "category": "registration", "passed": 2, "failed": 0 },
    { "category": "security", "passed": 3, "failed": 0 }
  ],
  "results": [
    {
      "name": "authenticate/wrong-origin",
      "category": "authentication",
      "description": "An assertion whose client data reports a foreign origin must be rejected.",
      "expectation": "reject",
      "outcome": "pass",
      "detail": "rejected as expected: harbor: origin \"https://evil.example\", want \"https://harbor.example\"",
      "duration_ns": 41000
    }
  ]
}
```

A full sample lives at [`examples/sample-report.json`](examples/sample-report.json).
Regenerate it any time with `make report` or the `-out` flag.

Programmatic use from Go is equally direct (see
[`examples/demo_test.go`](examples/demo_test.go)):

```go
results := harbor.RunScenarios(harbor.DefaultScenarios())
report := harbor.BuildReport(results)
_ = report.WriteJSON(os.Stdout) // or report.WriteText(os.Stdout)
```

---

## The browser lab

The lab in [`web/`](web/) is a small, dependency-light TypeScript app. It:

- validates an untrusted report against the `passkeypier/report/v1` schema
  before rendering, rejecting bad types, unknown outcomes, and a wrong schema;
- renders a pass or fail banner, per-category cards, and expandable scenarios;
- builds DOM nodes with `textContent` and never `innerHTML` from report data;
- runs fully offline, with no CDN, fonts, images, or telemetry. A sample report
