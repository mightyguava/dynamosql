package testutil

import (
	"bytes"
	"encoding/json"

	"github.com/alecthomas/repr"
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

// Repr pretty prints the value using alecthomas/repr
func Repr(v interface{}) string {
	return repr.String(v,
		repr.Indent("  "),
		repr.OmitEmpty(true),
		repr.IgnoreGoStringer())
}
