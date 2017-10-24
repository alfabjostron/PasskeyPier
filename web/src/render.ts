// DOM rendering for the passkeypier browser lab. Dependency-free; builds nodes
// with textContent (never innerHTML from report data) to avoid injection.

import { Report, ScenarioResult } from "./types.js";

function el<K extends keyof HTMLElementTagNameMap>(
  tag: K,
  className?: string,
  text?: string,
): HTMLElementTagNameMap[K] {
  const node = document.createElement(tag);
