package harbor

import (
	"encoding/json"
	"fmt"
	"io"
	"sort"
	"strings"
	"time"
)

// ReportSchemaVersion identifies the JSON report shape consumed by the
// TypeScript browser lab. Bump on breaking changes.
const ReportSchemaVersion = "passkeypier/report/v1"

// Report is the top-level conformance report emitted after a scenario run.
type Report struct {
	Schema      string           `json:"schema"`
	Tool        string           `json:"tool"`
	GeneratedAt string           `json:"generated_at"`
	Summary     ReportSummary    `json:"summary"`
	Categories  []CategoryTotals `json:"categories"`
	Results     []ScenarioResult `json:"results"`
}

// ReportSummary aggregates high-level pass/fail counts.
type ReportSummary struct {
	Total     int  `json:"total"`
	Passed    int  `json:"passed"`
	Failed    int  `json:"failed"`
	AllPassed bool `json:"all_passed"`
}

// CategoryTotals aggregates results per scenario category.
type CategoryTotals struct {
	Category string `json:"category"`
	Passed   int    `json:"passed"`
	Failed   int    `json:"failed"`
}

// BuildReport turns raw scenario results into a structured report.
func BuildReport(results []ScenarioResult) Report {
	summary := ReportSummary{Total: len(results)}
	catMap := map[string]*CategoryTotals{}
	for _, r := range results {
		if r.Outcome == OutcomePass {
			summary.Passed++
		} else {
			summary.Failed++
		}
		c, ok := catMap[r.Category]
		if !ok {
			c = &CategoryTotals{Category: r.Category}
			catMap[r.Category] = c
		}
		if r.Outcome == OutcomePass {
			c.Passed++
		} else {
			c.Failed++
		}
	}
	summary.AllPassed = summary.Failed == 0 && summary.Total > 0

	cats := make([]CategoryTotals, 0, len(catMap))
	for _, c := range catMap {
		cats = append(cats, *c)
	}
	sort.Slice(cats, func(i, j int) bool { return cats[i].Category < cats[j].Category })

	return Report{
		Schema:      ReportSchemaVersion,
		Tool:        "passkeypier",
		GeneratedAt: time.Now().UTC().Format(time.RFC3339),
		Summary:     summary,
		Categories:  cats,
		Results:     results,
	}
}

// WriteJSON writes the report as indented JSON.
func (r Report) WriteJSON(w io.Writer) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(r)
}

// WriteText writes a human-readable text report suitable for a terminal.
func (r Report) WriteText(w io.Writer) error {
	var b strings.Builder
	b.WriteString("passkeypier conformance report\n")
	b.WriteString("==============================\n")
	fmt.Fprintf(&b, "generated: %s\n", r.GeneratedAt)
	fmt.Fprintf(&b, "schema:    %s\n\n", r.Schema)

	for _, res := range r.Results {
		mark := "PASS"
		if res.Outcome != OutcomePass {
			mark = "FAIL"
		}
		fmt.Fprintf(&b, "[%s] %-42s (%s, expect %s)\n", mark, res.Name, res.Category, res.Expectation)
		fmt.Fprintf(&b, "       %s\n", res.Detail)
	}

	b.WriteString("\nby category:\n")
	for _, c := range r.Categories {
		fmt.Fprintf(&b, "  %-14s pass=%d fail=%d\n", c.Category, c.Passed, c.Failed)
	}

	fmt.Fprintf(&b, "\nsummary: %d passed, %d failed of %d total\n",
		r.Summary.Passed, r.Summary.Failed, r.Summary.Total)
	if r.Summary.AllPassed {
		b.WriteString("result: ALL SCENARIOS PASSED\n")
	} else {
		b.WriteString("result: CONFORMANCE FAILURES PRESENT\n")
	}

	_, err := io.WriteString(w, b.String())
	return err
}
