// Command passkeypier is the CLI front-end for the Passkey/WebAuthn conformance
// lab. It runs virtual registration and authentication ceremonies and emits
// conformance reports in text or JSON form.
//
// passkeypier is an educational conformance-exploration tool. It is not a
// FIDO-certified product and performs no attestation trust verification.
package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/alfabjostron/passkeypier/internal/harbor"
)

const usage = `passkeypier - virtual Passkey/WebAuthn conformance lab

usage:
  passkeypier run     [-format text|json] [-out FILE]   run the conformance suite
  passkeypier demo    [-uv required|preferred|discouraged]  run one register+auth ceremony
  passkeypier list                                       list scenarios
  passkeypier version                                    print version

passkeypier is an educational tool and is not FIDO-certified.
`

const version = "0.1.0"

