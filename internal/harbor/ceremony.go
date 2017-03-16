package harbor

import (
	"bytes"
	"crypto/ed25519"
	"errors"
	"fmt"
)

// RegistrationResult is the authenticator/client response to a create ceremony,
// mirroring the fields an RP consumes from an AuthenticatorAttestationResponse.
type RegistrationResult struct {
	CredentialID      []byte
	ClientDataJSON    []byte
	AuthenticatorData []byte
	PublicKey         []byte // raw Ed25519 public key
	UserVerified      bool
}

// AuthenticationResult is the response to a get ceremony, mirroring the fields
// of an AuthenticatorAssertionResponse.
type AuthenticationResult struct {
	CredentialID      []byte
	ClientDataJSON    []byte
	AuthenticatorData []byte
	Signature         []byte
	UserVerified      bool
	UserHandle        []byte
