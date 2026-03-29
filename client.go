package vk

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/andr-235/vk_api/internal/transport"
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

func (c *Client) transportConfig() transport.Config {
	tokenSource := transport.TokenInParams
	if c.tokenSource == TokenInHeader {
		tokenSource = transport.TokenInHeader
	}

	return transport.Config{
		BaseURL:     c.baseURL,
		Version:     c.version,
		Lang:        c.lang,
		TestMode:    c.testMode,
		Token:       c.token,
		TokenSource: tokenSource,
		HTTPClient:  c.httpClient,
	}
}

func (c *Client) Call(ctx context.Context, method string, params, out any) error {
	cfg := c.transportConfig()

	values, err := transport.EncodeValues(params)
	if err != nil {
		return fmt.Errorf("vk: encode params: %w", err)
	}

	values.Set("v", cfg.Version)

	if cfg.Lang != "" {
		values.Set("lang", cfg.Lang)
	}
	if cfg.TestMode {
		values.Set("test_mode", "1")
	}
	if cfg.Token != "" && cfg.TokenSource == transport.TokenInParams {
		values.Set("access_token", cfg.Token)
	}

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		cfg.BaseURL+method+"?"+values.Encode(),
		nil,
	)
	if err != nil {
		return fmt.Errorf("vk: build request: %w", err)
	}

	req.Header.Set("Accept", "application/json")

	if cfg.Token != "" && cfg.TokenSource == transport.TokenInHeader {
		req.Header.Set("Authorization", "Bearer "+cfg.Token)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("vk: do request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("vk: read response: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("vk: unexpected http status %d: %s", resp.StatusCode, string(body))
	}

	var respEnv responseEnvelope
	if err := json.Unmarshal(body, &respEnv); err != nil {
		return fmt.Errorf("vk: decode envelope: %w", err)
	}

	if respEnv.Error != nil {
		return &VKError{
			Code:             respEnv.Error.Code,
			Message:          respEnv.Error.Message,
			CaptchaSID:       respEnv.Error.CaptchaSID,
			CaptchaImg:       respEnv.Error.CaptchaImg,
			RedirectURI:      respEnv.Error.RedirectURI,
			ConfirmationText: respEnv.Error.ConfirmationText,
		}
	}

	if out == nil || len(respEnv.Response) == 0 {
		return nil
	}

	if err := json.Unmarshal(respEnv.Response, out); err != nil {
		return fmt.Errorf("vk: decode response payload: %w", err)
	}

	return nil
}

func (c *Client) CallWithRawHandler(ctx context.Context, method string, params any, handler func(json.RawMessage) error) error {
	cfg := c.transportConfig()

	values, err := transport.EncodeValues(params)
	if err != nil {
		return fmt.Errorf("vk: encode params: %w", err)
	}

	values.Set("v", cfg.Version)

	if cfg.Lang != "" {
		values.Set("lang", cfg.Lang)
	}
	if cfg.TestMode {
		values.Set("test_mode", "1")
	}
	if cfg.Token != "" && cfg.TokenSource == transport.TokenInParams {
		values.Set("access_token", cfg.Token)
	}

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		cfg.BaseURL+method+"?"+values.Encode(),
		nil,
	)
	if err != nil {
		return fmt.Errorf("vk: build request: %w", err)
	}

	req.Header.Set("Accept", "application/json")

	if cfg.Token != "" && cfg.TokenSource == transport.TokenInHeader {
		req.Header.Set("Authorization", "Bearer "+cfg.Token)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("vk: do request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("vk: read response: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("vk: unexpected http status %d: %s", resp.StatusCode, string(body))
	}

	var respEnv responseEnvelope
	if err := json.Unmarshal(body, &respEnv); err != nil {
		return fmt.Errorf("vk: decode envelope: %w", err)
	}

	if respEnv.Error != nil {
		return &VKError{
			Code:             respEnv.Error.Code,
			Message:          respEnv.Error.Message,
			CaptchaSID:       respEnv.Error.CaptchaSID,
			CaptchaImg:       respEnv.Error.CaptchaImg,
			RedirectURI:      respEnv.Error.RedirectURI,
			ConfirmationText: respEnv.Error.ConfirmationText,
		}
	}

	if handler == nil || len(respEnv.Response) == 0 {
		return nil
	}

	return handler(respEnv.Response)
}
