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
}

// clientOrigin is the origin the virtual client reports. In a browser this is
// enforced by the user agent; here the client faithfully reports the RP origin
// unless a scenario deliberately overrides it.
type clientOrigin struct {
	origin      string
	crossOrigin bool
}

// Register performs a full create ceremony against the RP:
//
//	1. the RP issues options with a fresh challenge;
//	2. the client builds client data (webauthn.create);
//	3. the authenticator mints an Ed25519 credential and authenticator data;
