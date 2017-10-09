// Browser lab entrypoint. Wires the file picker, a "load bundled sample"
// button and drag-and-drop to the report parser and renderer. Fully offline:
// no network requests, no remote assets.

import { parseReportText, ReportError } from "./report.js";
import { renderReport, renderError } from "./render.js";

function getTarget(): HTMLElement {
  const t = document.getElementById("report");
  if (!t) throw new Error("missing #report element");
  return t;
}

function handleText(text: string): void {
  const target = getTarget();
  try {
    const report = parseReportText(text);
    renderReport(target, report);
  } catch (e) {
    if (e instanceof ReportError) {
      renderError(target, e.message);
    } else {
      renderError(target, `unexpected error: ${(e as Error).message}`);
    }
  }
