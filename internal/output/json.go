package output

import (
	"encoding/json"
	"io"
)

// RenderJSON serializes any data structure to indented JSON format.
func RenderJSON(w io.Writer, data interface{}) error {
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	return encoder.Encode(data)
}
