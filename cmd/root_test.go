package cmd

import (
	"os"
	"testing"
)

func TestGetAPIKey_FromFlag(t *testing.T) {
	apiKey = "flag-key"
	defer func() { apiKey = "" }()

	key, err := getAPIKey()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if key != "flag-key" {
		t.Errorf("expected flag-key, got %s", key)
	}
}

func TestGetAPIKey_FromEnv(t *testing.T) {
	apiKey = ""
	os.Setenv("ODDS_API_KEY", "env-key")
	defer os.Unsetenv("ODDS_API_KEY")

	key, err := getAPIKey()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if key != "env-key" {
		t.Errorf("expected env-key, got %s", key)
	}
}

func TestGetAPIKey_Missing(t *testing.T) {
	apiKey = ""
	os.Unsetenv("ODDS_API_KEY")

	_, err := getAPIKey()
	if err == nil {
		t.Error("expected error for missing API key")
	}
}

func TestUseColor_Default(t *testing.T) {
	noColor = false
	os.Unsetenv("NO_COLOR")
	if !useColor() {
		t.Error("expected color by default")
	}
}

func TestUseColor_Flag(t *testing.T) {
	noColor = true
	defer func() { noColor = false }()
	if useColor() {
		t.Error("expected no color with --no-color flag")
	}
}

func TestUseColor_Env(t *testing.T) {
	noColor = false
	os.Setenv("NO_COLOR", "1")
	defer os.Unsetenv("NO_COLOR")
	if useColor() {
		t.Error("expected no color with NO_COLOR env")
	}
}
