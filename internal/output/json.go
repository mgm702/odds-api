package output

import (
	"encoding/json"
	"fmt"
	"io"
)

type JSONWriter struct {
	Out io.Writer
}

func NewJSONWriter(out io.Writer) *JSONWriter {
	return &JSONWriter{Out: out}
}

func (j *JSONWriter) Write(v any) error {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return fmt.Errorf("marshaling JSON: %w", err)
	}
	_, err = fmt.Fprintln(j.Out, string(data))
	return err
}
