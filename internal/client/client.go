package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/mgm702/odds-api-cli/internal/cache"
	"github.com/mgm702/odds-api-cli/internal/model"
)

const baseURL = "https://api.the-odds-api.com"

type Client struct {
	BaseURL     string
	APIKey      string
	HTTPClient  *http.Client
	Verbose     bool
	CacheConfig CacheConfig
	cacheStore  *cache.Store
}

func New(apiKey string) *Client {
	apiBase := baseURL
	if v := strings.TrimSpace(os.Getenv("ODDS_API_BASE_URL")); v != "" {
		apiBase = strings.TrimRight(v, "/")
	}

	return &Client{
		BaseURL:     apiBase,
		APIKey:      apiKey,
		HTTPClient:  &http.Client{Timeout: 30 * time.Second},
		CacheConfig: DefaultCacheConfig(),
	}
}

type Response struct {
	StatusCode int
	Quota      model.QuotaInfo
	Body       io.ReadCloser
}

func (c *Client) Get(ctx context.Context, path string, params url.Values) (*Response, error) {
	if params == nil {
		params = url.Values{}
	}
	params.Set("apiKey", c.APIKey)

	cachePath := c.BaseURL + path
	if c.shouldUseCache(path) {
		key := cache.RequestKey(http.MethodGet, cachePath, params)
		store, err := c.cacheStoreOrInit()
		if err != nil {
			if c.Verbose {
				fmt.Fprintf(os.Stderr, "Cache disabled (init error): %v\n", err)
			}
		} else if c.CacheConfig.Mode != CacheModeRefresh {
			entry, err := store.Get(key, c.CacheConfig.TTL)
			if err == nil {
				if c.Verbose {
					fmt.Fprintf(os.Stderr, "Cache hit: %s\n", path)
				}
				return &Response{
					StatusCode: entry.StatusCode,
					Quota:      entry.Quota,
					Body:       io.NopCloser(bytes.NewReader(entry.Body)),
				}, nil
			}
			if c.Verbose {
				if err == cache.ErrNotFound {
					fmt.Fprintf(os.Stderr, "Cache miss: %s\n", path)
				} else {
					fmt.Fprintf(os.Stderr, "Cache bypass (read error): %v\n", err)
				}
			}
		} else if c.Verbose {
			fmt.Fprintf(os.Stderr, "Cache refresh mode (skip read): %s\n", path)
		}
	}

	u := fmt.Sprintf("%s%s?%s", c.BaseURL, path, params.Encode())

	if c.Verbose {
		redacted := strings.Replace(u, c.APIKey, "***", 1)
		fmt.Fprintf(os.Stderr, "GET %s\n", redacted)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, fmt.Errorf("building request: %w", err)
	}

	start := time.Now()
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	elapsed := time.Since(start)

	quota := parseQuota(resp.Header)

	if c.Verbose {
		fmt.Fprintf(os.Stderr, "Status: %d (%s)\n", resp.StatusCode, elapsed.Round(time.Millisecond))
		fmt.Fprintf(os.Stderr, "Quota: used=%d remaining=%d last_cost=%d\n", quota.Used, quota.Remaining, quota.LastCost)
	}

	if resp.StatusCode != http.StatusOK {
		defer resp.Body.Close()
		body, _ := io.ReadAll(resp.Body)
		return nil, &APIError{
			StatusCode: resp.StatusCode,
			Message:    string(body),
			Quota:      quota,
		}
	}

	if c.shouldUseCache(path) {
		if store, err := c.cacheStoreOrInit(); err == nil {
			body, readErr := io.ReadAll(resp.Body)
			if readErr != nil {
				resp.Body.Close()
				return nil, fmt.Errorf("reading response body: %w", readErr)
			}
			resp.Body.Close()
			resp.Body = io.NopCloser(bytes.NewReader(body))

			if putErr := store.Put(cache.RequestKey(http.MethodGet, cachePath, params), cache.Entry{
				StatusCode: resp.StatusCode,
				Quota:      quota,
				Body:       body,
			}); putErr != nil && c.Verbose {
				fmt.Fprintf(os.Stderr, "Cache bypass (write error): %v\n", putErr)
			}
		}
	}

	return &Response{
		StatusCode: resp.StatusCode,
		Quota:      quota,
		Body:       resp.Body,
	}, nil
}

func (c *Client) shouldUseCache(path string) bool {
	if !c.CacheConfig.Enabled {
		return false
	}
	if c.CacheConfig.Mode == CacheModeOff {
		return false
	}
	if c.CacheConfig.Mode == CacheModeRefresh {
		return true
	}
	return cacheablePath(path)
}

func cacheablePath(path string) bool {
	switch {
	case path == "/v4/sports":
		return true
	case strings.Contains(path, "/events/") && strings.HasSuffix(path, "/odds"):
		return false
	case strings.Contains(path, "/events/") && strings.HasSuffix(path, "/markets"):
		return true
	case strings.HasSuffix(path, "/events"):
		return true
	case strings.HasSuffix(path, "/odds"):
		return true
	case strings.HasSuffix(path, "/scores"):
		return true
	case strings.HasSuffix(path, "/participants"):
		return true
	default:
		return strings.Contains(path, "/historical/")
	}
}

func (c *Client) GetQuotaOnly(ctx context.Context) (model.QuotaInfo, error) {
	params := url.Values{}
	params.Set("apiKey", c.APIKey)

	u := fmt.Sprintf("%s/v4/sports?%s", c.BaseURL, params.Encode())

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return model.QuotaInfo{}, fmt.Errorf("building request: %w", err)
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return model.QuotaInfo{}, fmt.Errorf("request failed: %w", err)
	}
	resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return model.QuotaInfo{}, &APIError{
			StatusCode: resp.StatusCode,
			Message:    "failed to fetch quota",
		}
	}

	return parseQuota(resp.Header), nil
}

func Decode[T any](resp *Response) (T, error) {
	var result T
	defer resp.Body.Close()
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return result, fmt.Errorf("decoding response: %w", err)
	}
	return result, nil
}

func parseQuota(h http.Header) model.QuotaInfo {
	return model.QuotaInfo{
		Remaining: headerInt(h, "X-Requests-Remaining"),
		Used:      headerInt(h, "X-Requests-Used"),
		LastCost:  headerInt(h, "X-Requests-Last"),
	}
}

func headerInt(h http.Header, key string) int {
	v := h.Get(key)
	if v == "" {
		return 0
	}
	n, _ := strconv.Atoi(v)
	return n
}

type APIError struct {
	StatusCode int
	Message    string
	Quota      model.QuotaInfo
}

func (e *APIError) Error() string {
	switch {
	case e.StatusCode == 401:
		return "unauthorized: invalid API key"
	case e.StatusCode == 422:
		return fmt.Sprintf("invalid request: %s", e.Message)
	case e.StatusCode == 429:
		return "rate limited: quota exceeded"
	case e.StatusCode >= 500:
		return fmt.Sprintf("server error (%d): %s", e.StatusCode, e.Message)
	default:
		return fmt.Sprintf("API error (%d): %s", e.StatusCode, e.Message)
	}
}

func (e *APIError) IsUserError() bool {
	return e.StatusCode == 401 || e.StatusCode == 422 || e.StatusCode == 429
}
