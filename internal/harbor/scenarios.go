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
