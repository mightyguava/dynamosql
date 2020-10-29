package testutil

import (
	"bytes"
	"encoding/json"
)

// MarshalJSON serializes the value into formatted JSON, without HTML escaping.
func MarshalJSON(v interface{}) string {
	buf := &bytes.Buffer{}
	enc := json.NewEncoder(buf)
	enc.SetEscapeHTML(false)
	enc.SetIndent("", "  ")
	if err := enc.Encode(v); err != nil {
		panic("error encoding JSON")
	}
	return buf.String()
}
