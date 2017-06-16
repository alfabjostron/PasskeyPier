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
