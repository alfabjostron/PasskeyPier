package harbor

import (
	"crypto/ed25519"
	"crypto/rand"
	"errors"
	"fmt"
)

// COSE algorithm identifier for EdDSA over Ed25519 (RFC 8152 / RFC 9053).
const COSEAlgEdDSA = -8
