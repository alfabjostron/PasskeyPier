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
)

// Authenticator data flag bits (WebAuthn L2 sec. 6.1).
const (
	FlagUserPresent    byte = 1 << 0 // UP
	FlagUserVerified   byte = 1 << 2 // UV
	FlagBackupEligible byte = 1 << 3 // BE
	FlagBackupState    byte = 1 << 4 // BS
	FlagAttestedData   byte = 1 << 6 // AT
	FlagExtensionData  byte = 1 << 7 // ED
)

// ClientData is the client-collected data hashed and signed during a ceremony.
// The JSON serialization mirrors the browser's clientDataJSON structure.
type ClientData struct {
	Type        CeremonyType `json:"type"`
	Challenge   string       `json:"challenge"` // base64url
