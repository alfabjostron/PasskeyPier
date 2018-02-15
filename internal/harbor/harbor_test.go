package harbor

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
)

const (
	testRPID   = "harbor.example"
	testOrigin = "https://harbor.example"
)

func mustChallenge(t *testing.T) Challenge {
	t.Helper()
	ch, err := NewChallenge(DefaultChallengeLen)
	if err != nil {
		t.Fatalf("NewChallenge: %v", err)
	}
	return ch
}

func sha256Of(t *testing.T, clientDataJSON []byte) []byte {
	t.Helper()
	var cd ClientData
	if err := json.Unmarshal(clientDataJSON, &cd); err != nil {
		t.Fatalf("unmarshal client data: %v", err)
	}
	h, err := cd.Hash()
	if err != nil {
		t.Fatalf("hash: %v", err)
	}
	return h
}

func TestNewChallengeLength(t *testing.T) {
	ch, err := NewChallenge(DefaultChallengeLen)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(ch) != DefaultChallengeLen {
		t.Fatalf("challenge length = %d, want %d", len(ch), DefaultChallengeLen)
	}
	if _, err := NewChallenge(8); err == nil {
		t.Fatal("expected error for short challenge, got nil")
	}
}

func TestChallengeUniqueness(t *testing.T) {
	seen := map[string]bool{}
	for i := 0; i < 256; i++ {
		ch := mustChallenge(t)
		s := ch.String()
		if seen[s] {
			t.Fatalf("duplicate challenge at iteration %d", i)
		}
		seen[s] = true
	}
}

func TestBase64URLRoundTrip(t *testing.T) {
	raw := []byte{0x00, 0xff, 0x10, 0x7f, 0x80, 0xab}
	enc := EncodeBase64URL(raw)
	if strings.ContainsAny(enc, "+/=") {
		t.Fatalf("encoding %q contains non-url-safe or padding characters", enc)
	}
	back, err := DecodeBase64URL(enc)
	if err != nil {
		t.Fatalf("decode: %v", err)
	}
	if !bytes.Equal(raw, back) {
		t.Fatalf("round trip mismatch: got %x want %x", back, raw)
	}
}

func TestAuthenticatorDataMarshalParse(t *testing.T) {
	ad := AuthenticatorData{
		RPIDHash:  RPIDHash(testRPID),
		Flags:     FlagUserPresent | FlagUserVerified,
		SignCount: 7,
		Trailing:  []byte{1, 2, 3},
	}
	raw := ad.Marshal()
	if len(raw) != 37+3 {
		t.Fatalf("marshaled length = %d, want %d", len(raw), 40)
	}
	parsed, err := ParseAuthenticatorData(raw)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	if parsed.SignCount != 7 {
		t.Fatalf("signCount = %d, want 7", parsed.SignCount)
	}
	if !parsed.Has(FlagUserVerified) {
		t.Fatal("expected UV flag set")
	}
	if !bytes.Equal(parsed.Trailing, []byte{1, 2, 3}) {
		t.Fatalf("trailing = %x, want 010203", parsed.Trailing)
	}
	if _, err := ParseAuthenticatorData(raw[:10]); err == nil {
		t.Fatal("expected error parsing truncated data")
	}
}

func TestClientDataHashDeterministic(t *testing.T) {
	cd := ClientData{Type: TypeGet, Challenge: "abc", Origin: testOrigin}
	h1, err := cd.Hash()
	if err != nil {
		t.Fatalf("hash: %v", err)
	}
	h2, _ := cd.Hash()
	if !bytes.Equal(h1, h2) {
		t.Fatal("client data hash not deterministic")
	}
