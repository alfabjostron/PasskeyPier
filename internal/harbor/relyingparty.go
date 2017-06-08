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
