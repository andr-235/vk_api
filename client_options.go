package vk

import "net/http"

type Option func(*Client)

type TokenSource int

const (
	TokenInParams TokenSource = iota
	TokenInHeader
)

func WithToken(token string) Option {
	return func(c *Client) { c.token = token }
}

func WithVersion(version string) Option {
	return func(c *Client) { c.version = version }
}

func WithLang(lang string) Option {
	return func(c *Client) { c.lang = lang }
}

func WithTestMode(enabled bool) Option {
	return func(c *Client) { c.testMode = enabled }
}

func WithBaseURL(baseURL string) Option {
	return func(c *Client) { c.baseURL = baseURL }
}

func WithHTTPClient(hc *http.Client) Option {
	return func(c *Client) { c.httpClient = hc }
}

func WithTokenSource(src TokenSource) Option {
	return func(c *Client) { c.tokenSource = src }
}
