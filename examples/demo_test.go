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
		Challenge:        regChallenge,
		UserHandle:       []byte("mariner-1"),
		UserVerification: harbor.UVPreferred,
	}, harbor.Origin(origin))
	if err != nil {
		panic(err)
	}
	fmt.Printf("registered=%v userVerified=%v\n", rp.Store[harbor.EncodeBase64URL(reg.CredentialID)] != nil, reg.UserVerified)

	// --- Authentication ceremony (webauthn.get) ---
	authChallenge, err := harbor.NewChallenge(harbor.DefaultChallengeLen)
	if err != nil {
		panic(err)
	}
	auth, err := harbor.Authenticate(rp, va, harbor.AuthenticationOptions{
		Challenge:        authChallenge,
		UserVerification: harbor.UVPreferred,
	}, harbor.Origin(origin))
	if err != nil {
		panic(err)
	}
	fmt.Printf("authenticated userHandle=%s sigLen=%d\n", string(auth.UserHandle), len(auth.Signature))

	// Output:
	// registered=true userVerified=true
	// authenticated userHandle=mariner-1 sigLen=64
}

// ExampleConformanceSuite runs the built-in scenario suite and prints the
// summary counts, showing how a report is produced programmatically.
func Example_conformanceSuite() {
	results := harbor.RunScenarios(harbor.DefaultScenarios())
	report := harbor.BuildReport(results)
	fmt.Printf("all passed: %v (%d scenarios)\n", report.Summary.AllPassed, report.Summary.Total)
	// Output:
	// all passed: true (9 scenarios)
}

// draft note 11
