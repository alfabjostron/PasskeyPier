package harbor

import "fmt"

// UserVerification mirrors the WebAuthn UserVerificationRequirement enum.
type UserVerification string

const (
	UVRequired    UserVerification = "required"
	UVPreferred   UserVerification = "preferred"
	UVDiscouraged UserVerification = "discouraged"
)

// Validate ensures the policy value is one of the recognized enum members.
func (uv UserVerification) Validate() error {
	switch uv {
	case UVRequired, UVPreferred, UVDiscouraged:
		return nil
	default:
		return fmt.Errorf("harbor: invalid user verification requirement %q", uv)
	}
}

// RelyingParty holds server-side configuration and the credential store used to
// validate ceremonies. Origin is the single expected origin; RPID is the
// relying-party identifier (a registrable domain suffix of the origin host).
type RelyingParty struct {
	ID     string // e.g. "harbor.example"
	Origin string // e.g. "https://harbor.example"
	// Store maps base64url credential id -> registered public-key record.
	Store map[string]*RegisteredCredential
}

// RegisteredCredential is the RP's server-side view of a registered passkey.
type RegisteredCredential struct {
	ID            []byte
	PublicKey     []byte // raw Ed25519 public key (32 bytes)
	SignCount     uint32
	RPID          string
	UserHandle    []byte
	BackupEligible bool
}

