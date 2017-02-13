package harbor

import (
	"bytes"
	"encoding/json"
)

// decodeStrict decodes JSON, rejecting unknown fields to catch malformed or
