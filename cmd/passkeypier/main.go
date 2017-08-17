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
