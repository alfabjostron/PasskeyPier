// Package harbor implements a virtual Passkey/WebAuthn conformance lab.
//
// It models the client/authenticator/relying-party interactions of a passkey
// ceremony (registration and authentication) using only the Go standard
// library. The implementation is intentionally focused on the parts of the
// WebAuthn Level 2 / CTAP data model that can be exercised deterministically
// in a test harness: challenge generation, base64url handling, origin and
// relying-party identifier checks, signature counters, user-verification
// policy evaluation and Ed25519 (COSE alg -8 / OKP) signing and verification.
//
// This package is a teaching and conformance-exploration tool. It is NOT a
// FIDO-certified implementation and makes no certification claims. It does not
// implement attestation statement verification, CBOR, or transport bindings.
package harbor

import "encoding/base64"

// b64url is the unpadded base64url encoding mandated by WebAuthn for the
// serialization of binary values in the JSON client-data and credential
// exchange (see WebAuthn L2 sec. 5.2, base64url without trailing '=').
var b64url = base64.RawURLEncoding

