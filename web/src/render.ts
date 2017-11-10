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
      `${report.summary.passed}/${report.summary.total} passed`,
    ),
  );
  wrap.appendChild(meta);
  return wrap;
}

function renderCategories(report: Report): HTMLElement {
  const wrap = el("div", "categories");
  wrap.appendChild(el("h2", undefined, "By category"));
  const list = el("div", "cat-grid");
  for (const c of report.categories) {
    const card = el("div", c.failed === 0 ? "cat-card ok" : "cat-card bad");
    card.appendChild(el("div", "cat-name", c.category));
    card.appendChild(
      el("div", "cat-counts", `pass ${c.passed} · fail ${c.failed}`),
    );
    list.appendChild(card);
  }
  wrap.appendChild(list);
  return wrap;
}

function renderScenario(r: ScenarioResult): HTMLElement {
  const row = el("details", r.outcome === "pass" ? "scenario pass" : "scenario fail");
  const summary = el("summary");
  const badge = el(
    "span",
    r.outcome === "pass" ? "badge badge-pass" : "badge badge-fail",
    r.outcome.toUpperCase(),
  );
  summary.appendChild(badge);
  summary.appendChild(el("span", "sc-name", r.name));
  summary.appendChild(el("span", "sc-cat", r.category));
  summary.appendChild(
    el("span", "sc-exp", `expect ${r.expectation}`),
  );
  summary.appendChild(el("span", "sc-dur", formatDuration(r.duration_ns)));
  row.appendChild(summary);

  const body = el("div", "sc-body");
  body.appendChild(el("p", "sc-desc", r.description));
  body.appendChild(el("p", "sc-detail", r.detail));
  row.appendChild(body);
  return row;
}

// renderReport clears the target and renders the full report into it.
export function renderReport(target: HTMLElement, report: Report): void {
  target.textContent = "";
