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
