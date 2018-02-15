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
