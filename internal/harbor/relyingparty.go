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

// NewRelyingParty constructs an RP with an empty credential store.
func NewRelyingParty(id, origin string) *RelyingParty {
	return &RelyingParty{
		ID:     id,
		Origin: origin,
		Store:  make(map[string]*RegisteredCredential),
	}
}

// RegistrationOptions are the parameters the RP issues to begin a create
// ceremony (analogous to PublicKeyCredentialCreationOptions).
type RegistrationOptions struct {
	Challenge        Challenge
	UserHandle       []byte
	UserVerification UserVerification
}

// AuthenticationOptions are the parameters the RP issues to begin a get
// ceremony (analogous to PublicKeyCredentialRequestOptions).
type AuthenticationOptions struct {
	Challenge        Challenge
	UserVerification UserVerification
	// AllowCredentialID optionally restricts which credential may respond.
	AllowCredentialID []byte
// review note
