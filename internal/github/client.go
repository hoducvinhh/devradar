package github

import (
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

type Client struct {
	baseURL    string
	token      string
	httpClient *http.Client
	rateLimit  *RateLimit
}

type Option func(*Client)

func WithToken(token string) Option {
	return func(c *Client) { c.token = strings.TrimSpace(token) }
}

func WithBaseURL(url string) Option {
	return func(c *Client) { c.baseURL = strings.TrimRight(url, "/") }
}

func WithHTTPClient(httpClient *http.Client) Option {
	return func(c *Client) { c.httpClient = httpClient }
}
