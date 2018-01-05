// Type definitions mirroring the Go conformance report schema
// (harbor.Report, schema "passkeypier/report/v1"). Keep these in sync with
// internal/harbor/report.go and scenarios.go.

export type Outcome = "pass" | "fail";
export type Expectation = "accept" | "reject";

export interface ScenarioResult {
  name: string;
  category: string;
