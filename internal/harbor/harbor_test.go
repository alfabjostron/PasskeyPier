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
	if len(h1) != 32 {
		t.Fatalf("hash length = %d, want 32", len(h1))
	}
}

func TestRegisterAndAuthenticateHappyPath(t *testing.T) {
	rp := NewRelyingParty(testRPID, testOrigin)
	va := NewVirtualAuthenticator(true)

	reg, err := Register(rp, va, RegistrationOptions{
		Challenge:        mustChallenge(t),
		UserHandle:       []byte("user-1"),
		UserVerification: UVPreferred,
	}, Origin(testOrigin))
	if err != nil {
		t.Fatalf("Register: %v", err)
	}
	if !reg.UserVerified {
		t.Fatal("expected user verified on UV-capable authenticator")
	}
	if rp.Store[EncodeBase64URL(reg.CredentialID)] == nil {
		t.Fatal("credential not stored by RP")
	}

	auth, err := Authenticate(rp, va, AuthenticationOptions{
		Challenge:        mustChallenge(t),
		UserVerification: UVPreferred,
	}, Origin(testOrigin))
	if err != nil {
		t.Fatalf("Authenticate: %v", err)
	}
	if string(auth.UserHandle) != "user-1" {
		t.Fatalf("user handle = %q, want user-1", auth.UserHandle)
	}
	// The stored counter must have advanced.
	if rp.Store[EncodeBase64URL(reg.CredentialID)].SignCount == 0 {
		t.Fatal("expected signature counter to advance")
	}
}

func TestUVRequiredWithoutSupportRejected(t *testing.T) {
	rp := NewRelyingParty(testRPID, testOrigin)
	va := NewVirtualAuthenticator(false) // no UV capability
	_, err := Register(rp, va, RegistrationOptions{
		Challenge:        mustChallenge(t),
		UserVerification: UVRequired,
	}, Origin(testOrigin))
	if err == nil {
		t.Fatal("expected rejection when UV required but unsupported")
	}
}

func TestWrongOriginRejected(t *testing.T) {
	rp := NewRelyingParty(testRPID, testOrigin)
	va := NewVirtualAuthenticator(true)
	if _, err := Register(rp, va, RegistrationOptions{
		Challenge:        mustChallenge(t),
		UserVerification: UVPreferred,
	}, Origin(testOrigin)); err != nil {
		t.Fatalf("Register: %v", err)
	}
	_, err := Authenticate(rp, va, AuthenticationOptions{
		Challenge:        mustChallenge(t),
		UserVerification: UVPreferred,
	}, CrossOrigin("https://evil.example"))
	if err == nil {
		t.Fatal("expected origin mismatch rejection")
	}
	if !strings.Contains(err.Error(), "origin") {
		t.Fatalf("expected origin error, got: %v", err)
	}
}

func TestCounterRegressionRejected(t *testing.T) {
	if err := verifyCounter(10, 5); err == nil {
		t.Fatal("expected counter regression rejection")
	}
	if err := verifyCounter(0, 0); err != nil {
		t.Fatalf("counter 0->0 should be tolerated: %v", err)
	}
	if err := verifyCounter(1, 2); err != nil {
		t.Fatalf("counter 1->2 should pass: %v", err)
	}
}

func TestSignatureVerificationFailsOnTamper(t *testing.T) {
	rp := NewRelyingParty(testRPID, testOrigin)
	va := NewVirtualAuthenticator(true)
	if _, err := Register(rp, va, RegistrationOptions{
		Challenge:        mustChallenge(t),
		UserVerification: UVPreferred,
	}, Origin(testOrigin)); err != nil {
		t.Fatalf("Register: %v", err)
	}
	auth, err := Authenticate(rp, va, AuthenticationOptions{
		Challenge:        mustChallenge(t),
		UserVerification: UVPreferred,
	}, Origin(testOrigin))
	if err != nil {
		t.Fatalf("Authenticate: %v", err)
	}
	stored := rp.Store[EncodeBase64URL(auth.CredentialID)]
	cdHash := sha256Of(t, auth.ClientDataJSON)
	// Flip a signature byte.
	tampered := append([]byte(nil), auth.Signature...)
	tampered[0] ^= 0xff
	if err := VerifyAssertionSignature(stored.PublicKey, auth.AuthenticatorData, cdHash, tampered); err == nil {
		t.Fatal("expected tampered signature to fail verification")
	}
	// The untampered signature must verify.
	if err := VerifyAssertionSignature(stored.PublicKey, auth.AuthenticatorData, cdHash, auth.Signature); err != nil {
		t.Fatalf("valid signature failed verification: %v", err)
	}
}

func TestUserVerificationValidate(t *testing.T) {
	for _, uv := range []UserVerification{UVRequired, UVPreferred, UVDiscouraged} {
		if err := uv.Validate(); err != nil {
			t.Fatalf("%q should be valid: %v", uv, err)
		}
	}
	if err := UserVerification("bogus").Validate(); err == nil {
		t.Fatal("expected invalid UV to error")
	}
