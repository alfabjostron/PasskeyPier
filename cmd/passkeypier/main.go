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

func main() {
	if len(os.Args) < 2 {
		fmt.Fprint(os.Stderr, usage)
		os.Exit(2)
	}
	switch os.Args[1] {
	case "run":
		os.Exit(cmdRun(os.Args[2:]))
	case "demo":
		os.Exit(cmdDemo(os.Args[2:]))
	case "list":
		os.Exit(cmdList())
	case "version":
		fmt.Printf("passkeypier %s\n", version)
	case "-h", "--help", "help":
		fmt.Print(usage)
	default:
		fmt.Fprintf(os.Stderr, "unknown command %q\n\n%s", os.Args[1], usage)
		os.Exit(2)
	}
}

func cmdRun(args []string) int {
	fs := flag.NewFlagSet("run", flag.ExitOnError)
	format := fs.String("format", "text", "report format: text or json")
	out := fs.String("out", "", "output file (default stdout)")
	_ = fs.Parse(args)

	results := harbor.RunScenarios(harbor.DefaultScenarios())
	report := harbor.BuildReport(results)

	w := os.Stdout
	if *out != "" {
		f, err := os.Create(*out)
		if err != nil {
			fmt.Fprintf(os.Stderr, "passkeypier: %v\n", err)
			return 1
		}
		defer f.Close()
		w = f
	}

	var err error
