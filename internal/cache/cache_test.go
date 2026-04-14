package cache

import (
	"errors"
	"net/url"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestRequestKey_StableAcrossParamOrder(t *testing.T) {
	q1 := url.Values{}
	q1.Set("regions", "us")
	q1.Set("markets", "h2h,spreads")
	q1.Set("apiKey", "secret-a")

	q2 := url.Values{}
	q2.Set("markets", "h2h,spreads")
	q2.Set("apiKey", "secret-b")
	q2.Set("regions", "us")

	k1 := RequestKey("get", "/v4/sports/nba/odds", q1)
	k2 := RequestKey("GET", "/v4/sports/nba/odds", q2)

	if k1 != k2 {
		t.Fatalf("expected same key, got %s vs %s", k1, k2)
	}
}

func TestRequestKey_ChangesWhenParamsChange(t *testing.T) {
	q1 := url.Values{}
	q1.Set("regions", "us")

	q2 := url.Values{}
	q2.Set("regions", "uk")

	k1 := RequestKey("GET", "/v4/sports/nba/odds", q1)
	k2 := RequestKey("GET", "/v4/sports/nba/odds", q2)

	if k1 == k2 {
		t.Fatal("expected different keys for different params")
	}
}

func TestStore_PutAndGet(t *testing.T) {
	dir := t.TempDir()
	s := NewStore(dir)
	key := "abc"
	body := []byte(`{"ok":true}`)

	err := s.Put(key, Entry{StatusCode: 200, Body: body})
	if err != nil {
		t.Fatalf("put failed: %v", err)
	}

	entry, err := s.Get(key, time.Minute)
	if err != nil {
		t.Fatalf("get failed: %v", err)
	}
	if string(entry.Body) != string(body) {
		t.Fatalf("unexpected body: %s", string(entry.Body))
	}
}

func TestStore_TTLExpired(t *testing.T) {
	dir := t.TempDir()
	s := NewStore(dir)
	key := "ttl"

	err := s.Put(key, Entry{StoredAt: time.Now().Add(-2 * time.Minute), StatusCode: 200, Body: []byte("[]")})
	if err != nil {
		t.Fatalf("put failed: %v", err)
	}

	_, err = s.Get(key, time.Second)
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("expected not found for expired key, got: %v", err)
	}
}

func TestResolveDir_EnvOverride(t *testing.T) {
	base := t.TempDir()
	t.Setenv("ODDS_CACHE_DIR", base)
	t.Setenv("XDG_CACHE_HOME", filepath.Join(base, "xdg"))

	dir, err := ResolveDir()
	if err != nil {
		t.Fatalf("resolve dir failed: %v", err)
	}
	if dir != base {
		t.Fatalf("expected %s, got %s", base, dir)
	}
}

func TestStore_CorruptFile(t *testing.T) {
	dir := t.TempDir()
	s := NewStore(dir)
	key := "bad"
	path := filepath.Join(dir, key+".json")
	if err := os.WriteFile(path, []byte("{not-json"), 0o644); err != nil {
		t.Fatalf("write corrupt file: %v", err)
	}

	_, err := s.Get(key, time.Minute)
	if err == nil {
		t.Fatal("expected decode error")
	}
}
