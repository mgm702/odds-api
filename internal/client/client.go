package client

import (
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

	"github.com/mgm702/odds-api-cli/internal/model"
)

const baseURL = "https://api.the-odds-api.com"

type Client struct {
	BaseURL    string
	APIKey     string
	HTTPClient *http.Client
	Verbose    bool
}

func New(apiKey string) *Client {
	return &Client{
		BaseURL:    baseURL,
		APIKey:     apiKey,
		HTTPClient: &http.Client{Timeout: 30 * time.Second},
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

	return &Response{
		StatusCode: resp.StatusCode,
		Quota:      quota,
		Body:       resp.Body,
	}, nil
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
