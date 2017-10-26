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
