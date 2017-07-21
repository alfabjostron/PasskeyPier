package harbor

import (
	"crypto/ed25519"
	"errors"
	"time"
)

// nowNS returns the current monotonic clock reading in nanoseconds.
func nowNS() int64 { return time.Now().UnixNano() }

// authenticateWithChallengeMismatch signs an assertion over signChallenge but
// asks the RP to verify against verifyChallenge, modeling a replayed or
// mismatched challenge. It returns the RP verification error.
func authenticateWithChallengeMismatch(rp *RelyingParty, va *VirtualAuthenticator, signChallenge, verifyChallenge Challenge, origin string) error {
	cred, err := selectCredential(va, rp.ID, nil)
	if err != nil {
		return err
	}
	cred.signCount++

	cd := ClientData{
		Type:      TypeGet,
		Challenge: signChallenge.String(),
		Origin:    origin,
	}
	cdJSON, err := cd.Marshal()
	if err != nil {
		return err
	}
	// The RP verifies the client data against the challenge it actually issued.
	return verifyClientData(cdJSON, TypeGet, verifyChallenge, rp.Origin)
}

// tamperedSignatureCheck produces a valid assertion, flips a signature byte,
// and runs the RP signature verification, which must fail.
func tamperedSignatureCheck(rp *RelyingParty, reg *RegistrationResult) error {
	stored, ok := rp.Store[EncodeBase64URL(reg.CredentialID)]
	if !ok {
		return errors.New("harbor: credential missing from store")
	}
	// Reconstruct a signable payload from the registration authenticator data
	// plus a fresh client-data hash.
	cd := ClientData{Type: TypeGet, Challenge: "AAAAAAAAAAAAAAAAAAAAAA", Origin: rp.Origin}
	cdHash, err := cd.Hash()
	if err != nil {
		return err
	}
	ad, err := ParseAuthenticatorData(reg.AuthenticatorData)
	if err != nil {
		return err
	}
	// Rebuild minimal auth data (no attested credential data) for the assertion.
	assertAD := AuthenticatorData{RPIDHash: ad.RPIDHash, Flags: FlagUserPresent | FlagUserVerified, SignCount: 1}.Marshal()

	// We do not have the private key here, so simulate a forged signature by
