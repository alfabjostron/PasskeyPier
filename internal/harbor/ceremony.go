package harbor

import (
	"bytes"
	"crypto/ed25519"
	"errors"
	"fmt"
)

// RegistrationResult is the authenticator/client response to a create ceremony,
// mirroring the fields an RP consumes from an AuthenticatorAttestationResponse.
type RegistrationResult struct {
	CredentialID      []byte
	ClientDataJSON    []byte
	AuthenticatorData []byte
	PublicKey         []byte // raw Ed25519 public key
	UserVerified      bool
}

// AuthenticationResult is the response to a get ceremony, mirroring the fields
// of an AuthenticatorAssertionResponse.
type AuthenticationResult struct {
	CredentialID      []byte
	ClientDataJSON    []byte
	AuthenticatorData []byte
	Signature         []byte
	UserVerified      bool
	UserHandle        []byte
}

// clientOrigin is the origin the virtual client reports. In a browser this is
// enforced by the user agent; here the client faithfully reports the RP origin
// unless a scenario deliberately overrides it.
type clientOrigin struct {
	origin      string
	crossOrigin bool
}

// Register performs a full create ceremony against the RP:
//
//	1. the RP issues options with a fresh challenge;
//	2. the client builds client data (webauthn.create);
//	3. the authenticator mints an Ed25519 credential and authenticator data;
//	4. the RP verifies origin, type, challenge, RP ID hash and UV policy, then
//	   stores the credential.
//
// The origin argument lets scenarios inject a mismatched origin. Pass the RP's
// own origin for the honest path.
func Register(rp *RelyingParty, va *VirtualAuthenticator, opts RegistrationOptions, origin clientOrigin) (*RegistrationResult, error) {
	if err := opts.UserVerification.Validate(); err != nil {
		return nil, err
	}
	uv := decideUV(va, opts.UserVerification)
	if opts.UserVerification == UVRequired && !uv {
		return nil, errUVUnsupported
	}

	cred, err := va.makeCredential(rp.ID)
	if err != nil {
		return nil, err
	}

	cd := ClientData{
		Type:        TypeCreate,
		Challenge:   opts.Challenge.String(),
		Origin:      origin.origin,
		CrossOrigin: origin.crossOrigin,
	}
	cdJSON, err := cd.Marshal()
	if err != nil {
		return nil, err
	}

	ad := cred.buildAuthData(uv, true, va.AAGUID)
	adBytes := ad.Marshal()

	// RP-side verification of the create ceremony (WebAuthn L2 sec. 7.1).
	if err := verifyClientData(cdJSON, TypeCreate, opts.Challenge, rp.Origin); err != nil {
		return nil, err
	}
	if err := verifyRPIDHash(ad, rp.ID); err != nil {
		return nil, err
	}
	if err := verifyUVFlag(ad, opts.UserVerification); err != nil {
		return nil, err
	}
	if !ad.Has(FlagUserPresent) {
		return nil, errors.New("harbor: user presence flag not set")
	}

	rp.Store[EncodeBase64URL(cred.id)] = &RegisteredCredential{
		ID:             cred.id,
		PublicKey:      append([]byte(nil), cred.pub...),
		SignCount:      ad.SignCount,
		RPID:           rp.ID,
		UserHandle:     append([]byte(nil), opts.UserHandle...),
		BackupEligible: cred.backedUp,
	}

	return &RegistrationResult{
		CredentialID:      append([]byte(nil), cred.id...),
		ClientDataJSON:    cdJSON,
		AuthenticatorData: adBytes,
		PublicKey:         append([]byte(nil), cred.pub...),
		UserVerified:      uv,
	}, nil
}

// Authenticate performs a full get ceremony: the authenticator signs
// authData || clientDataHash with the credential's Ed25519 key and the RP
// verifies the signature, origin, type, challenge, RP ID hash, UV policy and
// signature counter monotonicity (WebAuthn L2 sec. 7.2).
func Authenticate(rp *RelyingParty, va *VirtualAuthenticator, opts AuthenticationOptions, origin clientOrigin) (*AuthenticationResult, error) {
	if err := opts.UserVerification.Validate(); err != nil {
		return nil, err
	}

	cred, err := selectCredential(va, rp.ID, opts.AllowCredentialID)
	if err != nil {
		return nil, err
	}

	uv := decideUV(va, opts.UserVerification)
	if opts.UserVerification == UVRequired && !uv {
		return nil, errUVUnsupported
	}

	// The counter increments on each assertion (typical hardware behavior).
	cred.signCount++

	cd := ClientData{
		Type:        TypeGet,
		Challenge:   opts.Challenge.String(),
		Origin:      origin.origin,
		CrossOrigin: origin.crossOrigin,
	}
	cdJSON, err := cd.Marshal()
	if err != nil {
		return nil, err
	}
	cdHash, err := cd.Hash()
	if err != nil {
		return nil, err
	}

	ad := cred.buildAuthData(uv, false, va.AAGUID)
	adBytes := ad.Marshal()

	// The authenticator signs the concatenation authData || SHA-256(clientData).
	signed := append(append([]byte(nil), adBytes...), cdHash...)
	sig := ed25519.Sign(cred.priv, signed)

	res := &AuthenticationResult{
		CredentialID:      append([]byte(nil), cred.id...),
		ClientDataJSON:    cdJSON,
		AuthenticatorData: adBytes,
		Signature:         sig,
		UserVerified:      uv,
	}

	// RP-side verification.
	stored, ok := rp.Store[EncodeBase64URL(cred.id)]
	if !ok {
		return nil, errors.New("harbor: asserted credential is not registered")
	}
	res.UserHandle = append([]byte(nil), stored.UserHandle...)

	if err := verifyClientData(cdJSON, TypeGet, opts.Challenge, rp.Origin); err != nil {
		return nil, err
	}
	if err := verifyRPIDHash(ad, rp.ID); err != nil {
		return nil, err
	}
	if err := verifyUVFlag(ad, opts.UserVerification); err != nil {
		return nil, err
	}
	if !ad.Has(FlagUserPresent) {
