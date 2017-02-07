package harbor

import (
	"crypto/ed25519"
	"crypto/rand"
	"errors"
	"fmt"
)

// COSE algorithm identifier for EdDSA over Ed25519 (RFC 8152 / RFC 9053).
const COSEAlgEdDSA = -8

// VirtualAuthenticator models a roaming or platform authenticator holding one
// or more discoverable credentials. It signs assertions with Ed25519 and
// maintains a per-credential signature counter.
type VirtualAuthenticator struct {
	// AAGUID identifies the authenticator model (16 bytes). Zero for privacy.
	AAGUID [16]byte
	// SupportsUV indicates the authenticator can perform user verification
	// (biometric or PIN). If false, UV-required ceremonies must fail.
	SupportsUV bool
	// BackupEligible marks credentials as multi-device (synced passkeys).
	BackupEligible bool

	creds map[string]*credential // keyed by base64url credential id
}

type credential struct {
