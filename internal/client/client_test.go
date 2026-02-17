package client

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestGet_Success(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("apiKey") == "" {
			t.Error("expected apiKey param")
		}
		w.Header().Set("X-Requests-Remaining", "500")
		w.Header().Set("X-Requests-Used", "100")
		w.Header().Set("X-Requests-Last", "1")
		w.Write([]byte(`[{"key":"nfl"}]`))
	}))
	defer srv.Close()

	c := New("test-key")
	c.BaseURL = srv.URL

	resp, err := c.Get(context.Background(), "/v4/sports", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer resp.Body.Close()

	if resp.Quota.Remaining != 500 {
		t.Errorf("expected remaining=500, got %d", resp.Quota.Remaining)
	}
	if resp.Quota.Used != 100 {
		t.Errorf("expected used=100, got %d", resp.Quota.Used)
	}
	if resp.Quota.LastCost != 1 {
		t.Errorf("expected last_cost=1, got %d", resp.Quota.LastCost)
	}
}

func TestGet_Params(t *testing.T) {
	var gotParams url.Values
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotParams = r.URL.Query()
		w.Write([]byte(`[]`))
	}))
	defer srv.Close()

	c := New("test-key")
	c.BaseURL = srv.URL

	params := url.Values{}
	params.Set("regions", "us")
	params.Set("markets", "h2h")

	_, err := c.Get(context.Background(), "/v4/sports/nba/odds", params)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if gotParams.Get("regions") != "us" {
		t.Errorf("expected regions=us, got %s", gotParams.Get("regions"))
	}
	if gotParams.Get("markets") != "h2h" {
		t.Errorf("expected markets=h2h, got %s", gotParams.Get("markets"))
	}
	if gotParams.Get("apiKey") != "test-key" {
		t.Errorf("expected apiKey=test-key, got %s", gotParams.Get("apiKey"))
	}
}

func TestGet_401(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(401)
		w.Write([]byte(`unauthorized`))
	}))
	defer srv.Close()

	c := New("bad-key")
	c.BaseURL = srv.URL

	_, err := c.Get(context.Background(), "/v4/sports", nil)
	if err == nil {
		t.Fatal("expected error")
	}

	apiErr, ok := err.(*APIError)
	if !ok {
		t.Fatalf("expected APIError, got %T", err)
	}
	if !apiErr.IsUserError() {
		t.Error("expected user error")
	}
	if apiErr.StatusCode != 401 {
		t.Errorf("expected 401, got %d", apiErr.StatusCode)
	}
}

func TestGet_422(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(422)
		w.Write([]byte(`invalid params`))
	}))
	defer srv.Close()

	c := New("test-key")
	c.BaseURL = srv.URL

	_, err := c.Get(context.Background(), "/v4/sports", nil)
	apiErr, ok := err.(*APIError)
	if !ok {
		t.Fatalf("expected APIError, got %T", err)
	}
	if !apiErr.IsUserError() {
		t.Error("expected user error")
	}
}

func TestGet_429(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(429)
		w.Write([]byte(`rate limited`))
	}))
	defer srv.Close()

	c := New("test-key")
	c.BaseURL = srv.URL

	_, err := c.Get(context.Background(), "/v4/sports", nil)
	apiErr, ok := err.(*APIError)
	if !ok {
		t.Fatalf("expected APIError, got %T", err)
	}
	if !apiErr.IsUserError() {
		t.Error("expected user error for rate limit")
	}
}

func TestGet_500(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(500)
		w.Write([]byte(`server error`))
	}))
	defer srv.Close()

	c := New("test-key")
	c.BaseURL = srv.URL

	_, err := c.Get(context.Background(), "/v4/sports", nil)
	apiErr, ok := err.(*APIError)
	if !ok {
		t.Fatalf("expected APIError, got %T", err)
	}
	if apiErr.IsUserError() {
		t.Error("expected system error, not user error")
	}
}

func TestGetQuotaOnly(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Requests-Remaining", "358")
		w.Header().Set("X-Requests-Used", "142")
		w.Header().Set("X-Requests-Last", "0")
		w.Write([]byte(`[]`))
	}))
	defer srv.Close()

	c := New("test-key")
	c.BaseURL = srv.URL

	quota, err := c.GetQuotaOnly(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if quota.Remaining != 358 {
		t.Errorf("expected 358, got %d", quota.Remaining)
	}
	if quota.Used != 142 {
		t.Errorf("expected 142, got %d", quota.Used)
	}
	if quota.LastCost != 0 {
		t.Errorf("expected 0, got %d", quota.LastCost)
	}
}

func TestDecode(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`[{"key":"nfl","title":"NFL","group":"American Football","description":"","active":true,"has_outrights":false}]`))
	}))
	defer srv.Close()

	c := New("test-key")
	c.BaseURL = srv.URL

	resp, err := c.Get(context.Background(), "/v4/sports", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	type sport struct {
		Key   string `json:"key"`
		Title string `json:"title"`
	}
	result, err := Decode[[]sport](resp)
	if err != nil {
		t.Fatalf("decode error: %v", err)
	}
	if len(result) != 1 || result[0].Key != "nfl" {
		t.Errorf("unexpected result: %+v", result)
	}
}

func TestGet_Verbose(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`[]`))
	}))
	defer srv.Close()

	c := New("test-key")
	c.BaseURL = srv.URL
	c.Verbose = true

	resp, err := c.Get(context.Background(), "/v4/sports", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	resp.Body.Close()
}

func TestHeaderInt_Empty(t *testing.T) {
	h := http.Header{}
	if got := headerInt(h, "X-Missing"); got != 0 {
		t.Errorf("expected 0 for missing header, got %d", got)
	}
}

func TestHeaderInt_Invalid(t *testing.T) {
	h := http.Header{}
	h.Set("X-Bad", "notanumber")
	if got := headerInt(h, "X-Bad"); got != 0 {
		t.Errorf("expected 0 for invalid header, got %d", got)
	}
}
