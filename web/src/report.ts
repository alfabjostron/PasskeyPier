// Report parsing and validation. Dependency-free: uses only the DOM and
// standard JS. Validates the untrusted JSON against the expected schema shape
// before the UI renders it.

import { Report, ScenarioResult, EXPECTED_SCHEMA } from "./types.js";

export class ReportError extends Error {}

function isObject(v: unknown): v is Record<string, unknown> {
  return typeof v === "object" && v !== null && !Array.isArray(v);
}

function requireString(o: Record<string, unknown>, key: string): string {
  const v = o[key];
  if (typeof v !== "string") {
    throw new ReportError(`field "${key}" must be a string`);
  }
  return v;
}

function requireNumber(o: Record<string, unknown>, key: string): number {
  const v = o[key];
  if (typeof v !== "number" || Number.isNaN(v)) {
    throw new ReportError(`field "${key}" must be a number`);
  }
  return v;
}

function requireBool(o: Record<string, unknown>, key: string): boolean {
  const v = o[key];
  if (typeof v !== "boolean") {
    throw new ReportError(`field "${key}" must be a boolean`);
  }
  return v;
}

function parseScenario(v: unknown): ScenarioResult {
  if (!isObject(v)) {
    throw new ReportError("scenario result must be an object");
  }
  const outcome = requireString(v, "outcome");
  if (outcome !== "pass" && outcome !== "fail") {
    throw new ReportError(`invalid outcome "${outcome}"`);
  }
  const expectation = requireString(v, "expectation");
  if (expectation !== "accept" && expectation !== "reject") {
    throw new ReportError(`invalid expectation "${expectation}"`);
  }
  return {
    name: requireString(v, "name"),
    category: requireString(v, "category"),
    description: requireString(v, "description"),
    expectation,
