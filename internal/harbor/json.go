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

// Origin constructs a same-origin client context that faithfully reports the
// given origin (the honest path).
func Origin(origin string) clientOrigin {
	return clientOrigin{origin: origin, crossOrigin: false}
}

// CrossOrigin constructs a client context that reports a possibly different
// origin with the crossOrigin flag set, for modeling iframe / mismatch cases.
func CrossOrigin(origin string) clientOrigin {
	return clientOrigin{origin: origin, crossOrigin: true}
}
