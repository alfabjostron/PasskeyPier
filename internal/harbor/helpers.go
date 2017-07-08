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
