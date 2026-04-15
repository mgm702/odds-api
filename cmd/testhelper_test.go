package cmd

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/mgm702/odds-api-cli/internal/client"
)

func setupTestServer(t *testing.T, handler http.HandlerFunc) (*httptest.Server, *client.Client) {
	t.Helper()
	srv := httptest.NewServer(handler)
	t.Cleanup(srv.Close)
	c := client.New("test-key")
	c.BaseURL = srv.URL
	return srv, c
}

func withQuotaHeaders(w http.ResponseWriter) {
	w.Header().Set("X-Requests-Remaining", "500")
	w.Header().Set("X-Requests-Used", "100")
	w.Header().Set("X-Requests-Last", "1")
}

func fixtureHandler(fixturePath string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		withQuotaHeaders(w)
		data, err := os.ReadFile(fixturePath)
		if err != nil {
			w.WriteHeader(500)
			return
		}
		_, _ = w.Write(data)
	}
}
