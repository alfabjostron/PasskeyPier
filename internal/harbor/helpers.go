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
