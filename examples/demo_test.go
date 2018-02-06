// Package examples contains runnable, documentation-style examples for the
// passkeypier lab. These use Go's Example test convention so they compile and
// their output is verified by `go test`.
package examples

import (
	"fmt"

	"github.com/alfabjostron/passkeypier/internal/harbor"
)

// ExampleRegisterAndAuthenticate walks through an honest passkey ceremony:
// the relying party issues a challenge, the virtual authenticator mints an
// Ed25519 credential and signs an assertion, and the RP verifies everything.
func Example_registerAndAuthenticate() {
	const (
		rpID   = "harbor.example"
		origin = "https://harbor.example"
	)

	rp := harbor.NewRelyingParty(rpID, origin)
	va := harbor.NewVirtualAuthenticator(true) // UV-capable (biometric/PIN)

	// --- Registration ceremony (webauthn.create) ---
	regChallenge, err := harbor.NewChallenge(harbor.DefaultChallengeLen)
	if err != nil {
		panic(err)
	}
	reg, err := harbor.Register(rp, va, harbor.RegistrationOptions{
