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
