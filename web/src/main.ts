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
}

function wireFileInput(): void {
  const input = document.getElementById("file") as HTMLInputElement | null;
  if (!input) return;
  input.addEventListener("change", () => {
    const file = input.files?.[0];
    if (!file) return;
    file.text().then(handleText);
  });
}

function wireDrop(): void {
  const zone = document.getElementById("drop");
  if (!zone) return;
  zone.addEventListener("dragover", (e) => {
    e.preventDefault();
    zone.classList.add("drag");
  });
  zone.addEventListener("dragleave", () => zone.classList.remove("drag"));
  zone.addEventListener("drop", (e) => {
    e.preventDefault();
    zone.classList.remove("drag");
    const file = e.dataTransfer?.files?.[0];
    if (file) file.text().then(handleText);
  });
}

function wireSample(): void {
  const btn = document.getElementById("load-sample");
  if (!btn) return;
  const script = document.getElementById("sample-report");
  btn.addEventListener("click", () => {
    if (script && script.textContent) {
      handleText(script.textContent);
    } else {
