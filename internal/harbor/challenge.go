package harbor

import (
	"crypto/rand"
	"fmt"
)

// minChallengeLen is the minimum challenge length recommended by the WebAuthn
// specification (16 bytes). The default used by the lab is 32 bytes.
const (
	minChallengeLen     = 16
	DefaultChallengeLen = 32
)

// Challenge is a server-generated cryptographic random nonce that binds a
// ceremony to a single attempt and prevents replay.
type Challenge []byte

// NewChallenge returns a cryptographically secure random challenge of the
// requested length. Lengths below the 16-byte spec minimum are rejected.
func NewChallenge(n int) (Challenge, error) {
	if n < minChallengeLen {
		return nil, fmt.Errorf("harbor: challenge length %d below minimum %d", n, minChallengeLen)
	}
	buf := make([]byte, n)
	if _, err := rand.Read(buf); err != nil {
		return nil, fmt.Errorf("harbor: reading random challenge: %w", err)
	}
	return Challenge(buf), nil
}

// String renders the challenge as base64url, matching how it appears on the
// wire inside client data JSON.
func (c Challenge) String() string { return EncodeBase64URL(c) }
