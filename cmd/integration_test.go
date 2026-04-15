package cmd

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"strings"
	"sync/atomic"
	"testing"
)

func buildBinary(t *testing.T) string {
	t.Helper()
	binary := t.TempDir() + "/odds"
	cmd := exec.Command("go", "build", "-o", binary, "..")
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		t.Fatalf("failed to build binary: %v", err)
	}
	return binary
}

func TestIntegration_Version(t *testing.T) {
	binary := buildBinary(t)
	out, err := exec.Command(binary, "--version").CombinedOutput()
	if err != nil {
		t.Fatalf("version failed: %v\n%s", err, out)
	}
	if !strings.Contains(string(out), "odds version") {
		t.Errorf("expected version output, got: %s", out)
	}
}

func TestIntegration_Help(t *testing.T) {
	binary := buildBinary(t)
	out, err := exec.Command(binary, "--help").CombinedOutput()
	if err != nil {
		t.Fatalf("help failed: %v\n%s", err, out)
	}
	output := string(out)
	if !strings.Contains(output, "sports") {
		t.Error("expected sports in help")
	}
	if !strings.Contains(output, "credits") {
		t.Error("expected credits in help")
	}
	if !strings.Contains(output, "watch") {
		t.Error("expected watch in help")
	}
	if !strings.Contains(output, "historical") {
		t.Error("expected historical in help")
	}
}

func TestIntegration_SportsJSON(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Requests-Remaining", "500")
		w.Header().Set("X-Requests-Used", "0")
		w.Header().Set("X-Requests-Last", "0")
		_, _ = w.Write([]byte(`[{"key":"nfl","group":"American Football","title":"NFL","description":"","active":true,"has_outrights":false}]`))
	}))
	defer srv.Close()

	binary := buildBinary(t)

	cmd := exec.Command(binary, "sports", "--json")
	cmd.Env = append(os.Environ(), "ODDS_API_BASE_URL="+srv.URL)

	// Can't override base URL from env in current impl, so test the binary runs and produces output.
	// We test with real API key if available, otherwise skip.
	if os.Getenv("ODDS_API_KEY") == "" {
		t.Skip("ODDS_API_KEY not set, skipping integration test")
	}

	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("sports --json failed: %v\n%s", err, out)
	}

	var result []map[string]interface{}
	if err := json.Unmarshal(out, &result); err != nil {
		t.Fatalf("output is not valid JSON: %v\n%s", err, out)
	}
	if len(result) == 0 {
		t.Error("expected at least one sport")
	}
}

func TestIntegration_CreditsJSON(t *testing.T) {
	if os.Getenv("ODDS_API_KEY") == "" {
		t.Skip("ODDS_API_KEY not set, skipping integration test")
	}

	binary := buildBinary(t)
	out, err := exec.Command(binary, "credits", "--json").CombinedOutput()
	if err != nil {
		t.Fatalf("credits --json failed: %v\n%s", err, out)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(out, &result); err != nil {
		t.Fatalf("output is not valid JSON: %v\n%s", err, out)
	}
	if _, ok := result["remaining"]; !ok {
		t.Error("expected remaining field in credits output")
	}
}

func TestIntegration_MissingAPIKey(t *testing.T) {
	binary := buildBinary(t)
	cmd := exec.Command(binary, "sports")
	cmd.Env = []string{"PATH=" + os.Getenv("PATH"), "HOME=" + os.Getenv("HOME")}
	out, err := cmd.CombinedOutput()
	if err == nil {
		t.Fatal("expected error for missing API key")
	}
	if !strings.Contains(string(out), "API key required") {
		t.Errorf("expected API key error message, got: %s", out)
	}
}

func TestIntegration_InvalidCommand(t *testing.T) {
	binary := buildBinary(t)
	cmd := exec.Command(binary, "nonexistent")
	out, err := cmd.CombinedOutput()
	if err == nil {
		t.Fatal("expected error for invalid command")
	}
	if !strings.Contains(string(out), "unknown command") {
		t.Errorf("expected unknown command message, got: %s", out)
	}
}

func TestIntegration_OddsMissingRegions(t *testing.T) {
	if os.Getenv("ODDS_API_KEY") == "" {
		t.Skip("ODDS_API_KEY not set")
	}

	binary := buildBinary(t)
	out, err := exec.Command(binary, "odds", "basketball_nba").CombinedOutput()
	if err == nil {
		t.Fatal("expected error for missing --regions")
	}
	if !strings.Contains(string(out), "required flag") {
		t.Errorf("expected required flag error, got: %s", out)
	}
}

func TestIntegration_SubcommandHelp(t *testing.T) {
	binary := buildBinary(t)

	commands := []string{"sports", "events", "odds", "scores", "credits", "watch", "historical"}
	for _, subcmd := range commands {
		out, err := exec.Command(binary, subcmd, "--help").CombinedOutput()
		if err != nil {
			t.Errorf("%s --help failed: %v", subcmd, err)
		}
		if len(out) == 0 {
			t.Errorf("%s --help produced no output", subcmd)
		}
	}
}

func TestIntegration_HistoricalHelp(t *testing.T) {
	binary := buildBinary(t)

	subcmds := []string{"odds", "events", "event-odds"}
	for _, subcmd := range subcmds {
		out, err := exec.Command(binary, "historical", subcmd, "--help").CombinedOutput()
		if err != nil {
			t.Errorf("historical %s --help failed: %v", subcmd, err)
		}
		if !strings.Contains(string(out), "date") {
			t.Errorf("historical %s help should mention --date flag", subcmd)
		}
	}
}

func TestIntegration_CacheReuseWithMockServer(t *testing.T) {
	var hits int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&hits, 1)
		w.Header().Set("X-Requests-Remaining", "500")
		w.Header().Set("X-Requests-Used", "0")
		w.Header().Set("X-Requests-Last", "0")
		_, _ = w.Write([]byte(`[{"key":"nfl","group":"American Football","title":"NFL","description":"","active":true,"has_outrights":false}]`))
	}))
	defer srv.Close()

	binary := buildBinary(t)
	cacheDir := t.TempDir()
	baseEnv := append(os.Environ(), "ODDS_API_BASE_URL="+srv.URL)

	cmd1 := exec.Command(binary, "sports", "--json", "--api-key", "test-key", "--cache", "--cache-dir", cacheDir)
	cmd1.Env = baseEnv
	out1, err := cmd1.CombinedOutput()
	if err != nil {
		t.Fatalf("first sports call failed: %v\n%s", err, out1)
	}

	cmd2 := exec.Command(binary, "sports", "--json", "--api-key", "test-key", "--cache", "--cache-dir", cacheDir)
	cmd2.Env = baseEnv
	out2, err := cmd2.CombinedOutput()
	if err != nil {
		t.Fatalf("second sports call failed: %v\n%s", err, out2)
	}

	if atomic.LoadInt32(&hits) != 1 {
		t.Fatalf("expected 1 server hit due to cache reuse, got %d", atomic.LoadInt32(&hits))
	}

	var first []map[string]any
	var second []map[string]any
	if err := json.Unmarshal(out1, &first); err != nil {
		t.Fatalf("first output invalid json: %v", err)
	}
	if err := json.Unmarshal(out2, &second); err != nil {
		t.Fatalf("second output invalid json: %v", err)
	}
	if len(first) != len(second) {
		t.Fatalf("expected same output lengths, got %d and %d", len(first), len(second))
	}
}
