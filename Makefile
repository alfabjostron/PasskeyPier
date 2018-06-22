# passkeypier — build, test and report targets.
#
# Go core/CLI uses only the standard library. The TypeScript browser lab needs
# a local TypeScript compiler (npm i inside web/, or a global `tsc`).

GO      ?= go
