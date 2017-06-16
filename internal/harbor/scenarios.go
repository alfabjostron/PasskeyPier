package harbor

import (
	"fmt"
	"sort"
)

// Outcome is the result of running a single scenario check.
type Outcome string

const (
	OutcomePass Outcome = "pass"
	OutcomeFail Outcome = "fail"
)

// Expectation declares whether a scenario is expected to succeed or to be
// rejected by the RP. Negative scenarios ("expect rejection") are how the lab
// exercises the security checks: a scenario that *should* be rejected but is
// accepted is a conformance failure.
type Expectation string

const (
	ExpectAccept Expectation = "accept"
	ExpectReject Expectation = "reject"
)

// ScenarioResult records the outcome of one scenario.
type ScenarioResult struct {
	Name        string      `json:"name"`
	Category    string      `json:"category"`
	Description string      `json:"description"`
	Expectation Expectation `json:"expectation"`
	Outcome     Outcome     `json:"outcome"`
	// Detail carries the error message on rejection, or a success note.
	Detail string `json:"detail"`
	// DurationNS is the wall-clock duration of the scenario in nanoseconds.
	DurationNS int64 `json:"duration_ns"`
}

// Scenario is a named, self-contained conformance check.
type Scenario struct {
	Name        string
	Category    string
	Description string
	Expectation Expectation
	// Run executes the scenario and returns nil on RP acceptance or a non-nil
	// error when the RP rejects the ceremony.
	Run func() error
}

// evaluate compares the actual error against the expectation and yields an
// Outcome plus a human-readable detail string.
func (s Scenario) evaluate(err error) (Outcome, string) {
	switch s.Expectation {
	case ExpectAccept:
		if err == nil {
			return OutcomePass, "ceremony accepted as expected"
		}
		return OutcomeFail, fmt.Sprintf("expected acceptance but was rejected: %v", err)
	case ExpectReject:
		if err != nil {
			return OutcomePass, fmt.Sprintf("rejected as expected: %v", err)
		}
		return OutcomeFail, "expected rejection but ceremony was accepted"
	default:
		return OutcomeFail, fmt.Sprintf("unknown expectation %q", s.Expectation)
	}
}

// DefaultScenarios returns the built-in conformance scenario suite. Each
// scenario builds its own fresh RP and authenticator so runs are independent.
func DefaultScenarios() []Scenario {
	const (
		rpID     = "harbor.example"
		origin   = "https://harbor.example"
		badOrig  = "https://evil.example"
		otherRP  = "phish.example"
	)

	newPair := func(uv bool) (*RelyingParty, *VirtualAuthenticator) {
		return NewRelyingParty(rpID, origin), NewVirtualAuthenticator(uv)
	}

	// registerHelper performs an honest registration and returns state.
	registerHelper := func(rp *RelyingParty, va *VirtualAuthenticator, uv UserVerification) (*RegistrationResult, error) {
		ch, err := NewChallenge(DefaultChallengeLen)
		if err != nil {
			return nil, err
		}
		return Register(rp, va, RegistrationOptions{
			Challenge:        ch,
			UserHandle:       []byte("user-42"),
			UserVerification: uv,
		}, Origin(origin))
	}

	return []Scenario{
