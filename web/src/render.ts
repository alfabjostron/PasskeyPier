// DOM rendering for the passkeypier browser lab. Dependency-free; builds nodes
// with textContent (never innerHTML from report data) to avoid injection.

import { Report, ScenarioResult } from "./types.js";

function el<K extends keyof HTMLElementTagNameMap>(
  tag: K,
  className?: string,
  text?: string,
): HTMLElementTagNameMap[K] {
  const node = document.createElement(tag);
  if (className) node.className = className;
  if (text !== undefined) node.textContent = text;
  return node;
}

function formatDuration(ns: number): string {
  if (ns < 1000) return `${ns} ns`;
  if (ns < 1_000_000) return `${(ns / 1000).toFixed(1)} µs`;
  return `${(ns / 1_000_000).toFixed(2)} ms`;
}

function renderSummary(report: Report): HTMLElement {
  const wrap = el("div", "summary");
  const banner = el(
    "div",
    report.summary.all_passed ? "banner banner-pass" : "banner banner-fail",
    report.summary.all_passed
      ? "ALL SCENARIOS PASSED"
      : "CONFORMANCE FAILURES PRESENT",
  );
  wrap.appendChild(banner);

  const meta = el("div", "meta");
  meta.appendChild(el("span", "meta-item", `tool: ${report.tool}`));
  meta.appendChild(el("span", "meta-item", `schema: ${report.schema}`));
  meta.appendChild(el("span", "meta-item", `generated: ${report.generated_at}`));
  meta.appendChild(
    el(
      "span",
      "meta-item",
