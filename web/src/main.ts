// Browser lab entrypoint. Wires the file picker, a "load bundled sample"
// button and drag-and-drop to the report parser and renderer. Fully offline:
// no network requests, no remote assets.

import { parseReportText, ReportError } from "./report.js";
import { renderReport, renderError } from "./render.js";

