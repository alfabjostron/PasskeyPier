# Security Policy

## Supported versions

| Version | Supported |
|---------|-----------|
| latest release on `main` | yes |
| older tags | best effort |

## Reporting a vulnerability

PasskeyPier runs ceremonies locally and writes local reports; it performs
no network I/O and executes no downloaded code. The realistic surface is
malformed scenario inputs (parser panics, extreme allocations) and path
traversal in `-out` paths.

**How to report:** email security concerns to the maintainer rather than
opening a public issue. Include a minimal crashing scenario and the
version tag you tested against.

**What to expect:**

- acknowledgement within 7 days;
- a fix on `main` and a patch release within 30 days for confirmed
  vulnerabilities;
- credit in CHANGELOG.md and the release notes unless you prefer to stay
  anonymous.