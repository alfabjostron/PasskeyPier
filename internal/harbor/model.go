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
	Origin      string       `json:"origin"`
	CrossOrigin bool         `json:"crossOrigin"`
}

// Marshal serializes the client data to its canonical JSON byte form. This is
// the exact byte sequence whose SHA-256 hash is signed by the authenticator.
func (cd ClientData) Marshal() ([]byte, error) {
	return json.Marshal(cd)
}

// Hash returns the SHA-256 of the marshaled client data (the "clientDataHash").
func (cd ClientData) Hash() ([]byte, error) {
	raw, err := cd.Marshal()
	if err != nil {
		return nil, err
	}
	sum := sha256.Sum256(raw)
	return sum[:], nil
}

// AuthenticatorData models the authenticator data structure. Attested
// credential data and extensions are represented as opaque trailing bytes so
// that signature computation over the concatenation stays byte-accurate.
type AuthenticatorData struct {
	RPIDHash  [32]byte
	Flags     byte
	SignCount uint32
	Trailing  []byte // attested credential data + extensions, if present
}

// Marshal encodes authenticator data in the wire layout: rpIdHash(32) ||
// flags(1) || signCount(4, big-endian) || trailing.
func (ad AuthenticatorData) Marshal() []byte {
	out := make([]byte, 37+len(ad.Trailing))
	copy(out[0:32], ad.RPIDHash[:])
	out[32] = ad.Flags
	binary.BigEndian.PutUint32(out[33:37], ad.SignCount)
	copy(out[37:], ad.Trailing)
	return out
}

// ParseAuthenticatorData decodes the fixed-length prefix of authenticator data.
func ParseAuthenticatorData(raw []byte) (AuthenticatorData, error) {
	if len(raw) < 37 {
		return AuthenticatorData{}, fmt.Errorf("harbor: authenticator data too short: %d bytes", len(raw))
	}
	var ad AuthenticatorData
	copy(ad.RPIDHash[:], raw[0:32])
	ad.Flags = raw[32]
	ad.SignCount = binary.BigEndian.Uint32(raw[33:37])
	if len(raw) > 37 {
