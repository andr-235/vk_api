package transport

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	internalencode "github.com/andr-235/vk_api/internal/encode"
)

type Doer interface {
	Do(req *http.Request) (*http.Response, error)
}

type TokenSource int

const (
	TokenInParams TokenSource = iota
	TokenInHeader
)

type Config struct {
	BaseURL     string
	Version     string
	Lang        string
	TestMode    bool
	Token       string
	TokenSource TokenSource
	HTTPClient  Doer
}

type ResponseEnvelope struct {
	Response json.RawMessage `json:"response"`
	Error    *APIError       `json:"error"`
}

type RequestParam struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type APIError struct {
	Code          int            `json:"error_code"`
	Message       string         `json:"error_msg"`
	RequestParams []RequestParam `json:"request_params,omitempty"`

	CaptchaSID  string `json:"captcha_sid,omitempty"`
	CaptchaImg  string `json:"captcha_img,omitempty"`
	RedirectURI string `json:"redirect_uri,omitempty"`

	ConfirmationText string `json:"confirmation_text,omitempty"`
}

func Call(ctx context.Context, cfg Config, method string, params any, out any) error {
	if cfg.Version == "" {
		return fmt.Errorf("vk: api version is required")
	}
	if cfg.BaseURL == "" {
		return fmt.Errorf("vk: base url is required")
	}
	if cfg.HTTPClient == nil {
		return fmt.Errorf("vk: http client is required")
	}
	if method == "" {
		return fmt.Errorf("vk: method is required")
	}

	values, err := internalencode.Values(params)
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
	if cfg.Token != "" && cfg.TokenSource == TokenInParams {
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

	if cfg.Token != "" && cfg.TokenSource == TokenInHeader {
		req.Header.Set("Authorization", "Bearer "+cfg.Token)
	}

	resp, err := cfg.HTTPClient.Do(req)
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

	var env ResponseEnvelope
	if err := json.Unmarshal(body, &env); err != nil {
		return fmt.Errorf("vk: decode envelope: %w", err)
	}

	if env.Error != nil {
		return env.Error
	}

	if out == nil || len(env.Response) == 0 {
		return nil
	}

	if err := json.Unmarshal(env.Response, out); err != nil {
		return fmt.Errorf("vk: decode response payload: %w", err)
	}

	return nil
}

func (e *APIError) Error() string {
	if e == nil {
		return "<nil>"
	}
	if e.Message == "" {
		return fmt.Sprintf("vk api error %d", e.Code)
	}
	return fmt.Sprintf("vk api error %d: %s", e.Code, e.Message)
}

func QueryFromParams(params any) (url.Values, error) {
	return internalencode.Values(params)
}
