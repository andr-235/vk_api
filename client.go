package vk

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"
)

const (
	defaultBaseURL = "https://api.vk.ru/method/"
	defaultVersion = "5.199"
)

type Client struct {
	token       string
	version     string
	lang        string
	testMode    bool
	baseURL     string
	httpClient  *http.Client
	tokenSource TokenSource
}

func New(opts ...Option) *Client {
	c := &Client{
		version: defaultVersion,
		baseURL: defaultBaseURL,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		tokenSource: TokenInParams,
	}

	for _, opt := range opts {
		opt(c)
	}

	if c.httpClient == nil {
		c.httpClient = http.DefaultClient
	}
	if c.baseURL == "" {
		c.baseURL = defaultBaseURL
	}
	if !strings.HasSuffix(c.baseURL, "/") {
		c.baseURL += "/"
	}
	if c.version == "" {
		c.version = defaultVersion
	}

	return c
}

func (c *Client) endpoint(method string) (string, error) {
	if strings.TrimSpace(method) == "" {
		return "", errors.New("vk: method is required")
	}
	return fmt.Sprintf("%s%s", c.baseURL, method), nil
}
