package harbor

import (
	"bytes"
	"encoding/json"
)

// decodeStrict decodes JSON, rejecting unknown fields to catch malformed or
// tampered client data during verification.
func decodeStrict(data []byte, v any) error {
	dec := json.NewDecoder(bytes.NewReader(data))
	dec.DisallowUnknownFields()
	return dec.Decode(v)
}
