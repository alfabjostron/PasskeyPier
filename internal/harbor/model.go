package harbor

import (
	"crypto/sha256"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
)

// CeremonyType distinguishes the two WebAuthn ceremonies.
type CeremonyType string

const (
	// TypeCreate is the value of clientDataJSON.type during registration.
	TypeCreate CeremonyType = "webauthn.create"
	// TypeGet is the value of clientDataJSON.type during authentication.
	TypeGet CeremonyType = "webauthn.get"
