package github

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

const (
	DefaultBaseURL = "https://api.github.com"
	UserAgent      = "DevRadar-CLI/1.0"
	DefaultTimeout = 10 * time.Second
	MaxRetries     = 3
)

// Client handles communication with the GitHub REST API.
type Client struct {
	baseURL    string
	token      string
	httpClient *http.Client
	rateLimit  *RateLimit
}

// Option allows configuring the GitHub client.
type Option func(*Client)

// WithToken sets the GitHub Personal Access Token.
func WithToken(token string) Option {
	return func(c *Client) {
		c.token = strings.TrimSpace(token)
	}
}

// WithBaseURL overrides the default GitHub API base URL.
func WithBaseURL(url string) Option {
	return func(c *Client) {
		c.baseURL = strings.TrimRight(url, "/")
	}
}

// WithHTTPClient overrides the default HTTP client.
func WithHTTPClient(httpClient *http.Client) Option {
	return func(c *Client) {
		c.httpClient = httpClient
	}
}

// NewClient creates a new GitHub API client.
func NewClient(opts ...Option) *Client {
	c := &Client{
		baseURL: DefaultBaseURL,
		httpClient: &http.Client{
			Timeout: DefaultTimeout,
		},
		rateLimit: &RateLimit{},
	}

	for _, opt := range opts {
		opt(c)
	}

	return c
}

// RateLimit returns the current rate limit tracker.
func (c *Client) RateLimit() *RateLimit {
	return c.rateLimit
}

// HasToken returns true if client is authenticated with a token.
func (c *Client) HasToken() bool {
	return c.token != ""
}

// Do executes an HTTP request with retry logic and rate limit tracking.
func (c *Client) Do(req *http.Request) (*http.Response, error) {
	req.Header.Set("Accept", "application/vnd.github.v3+json")
	req.Header.Set("User-Agent", UserAgent)

	if c.token != "" {
		req.Header.Set("Authorization", "Bearer "+c.token)
	}

	var resp *http.Response
	var err error

	backoff := 500 * time.Millisecond

	for attempt := 0; attempt < MaxRetries; attempt++ {
		resp, err = c.httpClient.Do(req)
		if err == nil {
			c.rateLimit.UpdateFromHeaders(resp.Header)

			// If server error 5xx, retry with backoff
			if resp.StatusCode >= 500 && attempt < MaxRetries-1 {
				_ = resp.Body.Close()
				time.Sleep(backoff)
				backoff *= 2
				continue
			}

			// Handle Rate Limit Exceeded (403)
			if resp.StatusCode == http.StatusForbidden && strings.Contains(strings.ToLower(resp.Header.Get("X-RateLimit-Remaining")), "0") {
				_ = resp.Body.Close()
				return nil, fmt.Errorf("GitHub API rate limit exceeded. Please configure a personal access token using 'devradar auth set-token <token>'")
			}

			// Handle Unauthorized (401)
			if resp.StatusCode == http.StatusUnauthorized {
				_ = resp.Body.Close()
				return nil, fmt.Errorf("invalid GitHub token (HTTP 401 Unauthorized). Please check your token or run 'devradar auth set-token'")
			}

			return resp, nil
		}

		// On connection/timeout error, retry if attempts remain
		if attempt < MaxRetries-1 {
			time.Sleep(backoff)
			backoff *= 2
		}
	}

	return nil, fmt.Errorf("request failed after %d attempts: %w", MaxRetries, err)
}

// Get performs a GET request and unmarshals the JSON response into target.
func (c *Client) Get(ctx context.Context, path string, target interface{}) error {
	fullURL := path
	if !strings.HasPrefix(path, "http://") && !strings.HasPrefix(path, "https://") {
		fullURL = c.baseURL + path
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fullURL, nil)
	if err != nil {
		return fmt.Errorf("creating request: %w", err)
	}

	resp, err := c.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("github api error (%d): %s", resp.StatusCode, string(bodyBytes))
	}

	if target != nil {
		if err := json.NewDecoder(resp.Body).Decode(target); err != nil {
			return fmt.Errorf("decoding response json: %w", err)
		}
	}

	return nil
}
