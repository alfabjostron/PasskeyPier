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
