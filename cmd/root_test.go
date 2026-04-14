package cmd

import (
	"os"
	"testing"
	"time"
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

func TestCurrentCacheConfig_InvalidMode(t *testing.T) {
	cacheMode = "bad"
	defer func() { cacheMode = "smart" }()

	_, err := currentCacheConfig()
	if err == nil {
		t.Fatal("expected error for invalid cache mode")
	}
}

func TestCurrentCacheConfig_Valid(t *testing.T) {
	cacheEnabled = true
	cacheMode = "refresh"
	cacheTTL = 2 * time.Minute
	cacheDir = "/tmp/odds-cache"

	cfg, err := currentCacheConfig()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !cfg.Enabled {
		t.Fatal("expected cache enabled")
	}
	if cfg.Mode != "refresh" {
		t.Fatalf("expected refresh mode, got %s", cfg.Mode)
	}
	if cfg.TTL != 2*time.Minute {
		t.Fatalf("expected 2m TTL, got %s", cfg.TTL)
	}
	if cfg.Dir != "/tmp/odds-cache" {
		t.Fatalf("expected cache dir set, got %s", cfg.Dir)
	}
}
