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
