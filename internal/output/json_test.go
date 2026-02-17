package output

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/mgm702/odds-api-cli/internal/model"
)

func TestJSONWriter(t *testing.T) {
	var buf bytes.Buffer
	jw := NewJSONWriter(&buf)
	err := jw.Write([]model.Sport{
		{Key: "nfl", Title: "NFL", Active: true},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var result []model.Sport
	if err := json.Unmarshal(buf.Bytes(), &result); err != nil {
		t.Fatalf("output is not valid JSON: %v", err)
	}
	if len(result) != 1 || result[0].Key != "nfl" {
		t.Errorf("unexpected result: %+v", result)
	}
}

func TestJSONWriterQuota(t *testing.T) {
	var buf bytes.Buffer
	jw := NewJSONWriter(&buf)
	err := jw.Write(model.QuotaInfo{Used: 142, Remaining: 358, LastCost: 0})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var result model.QuotaInfo
	if err := json.Unmarshal(buf.Bytes(), &result); err != nil {
		t.Fatalf("output is not valid JSON: %v", err)
	}
	if result.Used != 142 {
		t.Errorf("expected used=142, got %d", result.Used)
	}
}
