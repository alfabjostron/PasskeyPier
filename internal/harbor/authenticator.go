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
	id        []byte
	priv      ed25519.PrivateKey
	pub       ed25519.PublicKey
	rpID      string
	signCount uint32
	backedUp  bool
}

// NewVirtualAuthenticator constructs an authenticator with the given UV
// capability. The AAGUID is left zero (privacy-preserving default).
func NewVirtualAuthenticator(supportsUV bool) *VirtualAuthenticator {
	return &VirtualAuthenticator{
		SupportsUV: supportsUV,
		creds:      make(map[string]*credential),
	}
}

// CredentialCount reports the number of resident credentials.
func (va *VirtualAuthenticator) CredentialCount() int { return len(va.creds) }

// makeCredential generates a fresh Ed25519 credential bound to rpID and stores
// it as a discoverable (resident) credential.
func (va *VirtualAuthenticator) makeCredential(rpID string) (*credential, error) {
	pub, priv, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return nil, fmt.Errorf("harbor: generating ed25519 key: %w", err)
	}
	id := make([]byte, 16)
	if _, err := rand.Read(id); err != nil {
		return nil, fmt.Errorf("harbor: generating credential id: %w", err)
	}
	c := &credential{
		id:        id,
		priv:      priv,
		pub:       pub,
		rpID:      rpID,
		signCount: 0,
		backedUp:  va.BackupEligible,
	}
	va.creds[EncodeBase64URL(id)] = c
	return c, nil
}

// lookup finds a resident credential by base64url id.
func (va *VirtualAuthenticator) lookup(credIDB64 string) (*credential, bool) {
	c, ok := va.creds[credIDB64]
	return c, ok
}

// buildAuthData assembles authenticator data with the appropriate flags.
// includeAttested controls whether the AT flag and attested credential data
