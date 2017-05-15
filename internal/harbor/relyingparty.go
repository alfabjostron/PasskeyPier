package harbor

import "fmt"

// UserVerification mirrors the WebAuthn UserVerificationRequirement enum.
type UserVerification string

const (
	UVRequired    UserVerification = "required"
	UVPreferred   UserVerification = "preferred"
	UVDiscouraged UserVerification = "discouraged"
)
