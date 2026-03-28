package vk

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	internalencode "github.com/andr-235/vk/internal/encode"
)

func (c *Client) Call(ctx context.Context, method string, params any, out any) error {
	endpoint, err := c.endpoint(method)
	if err != nil {
		return err
	}

	values, err := internalencode.Values(params)
	if err != nil {
		return fmt.Errorf("vk: encode params: %w", err)
	}

	if c.token != "" {
		values.Set("access_token", c.token)
	}
	if c.version != "" {
		values.Set("v", c.version)
	}
	if c.lang != "" {
		values.Set("lang", c.lang)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewBufferString(values.Encode()))
	if err != nil {
		return fmt.Errorf("vk: build request: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")

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

	var env responseEnvelope
	if err := json.Unmarshal(body, &env); err != nil {
		return fmt.Errorf("vk: decode envelope: %w", err)
	}

	if env.Error != nil {
		return newVKError(env.Error)
	}

	if out == nil || len(env.Response) == 0 {
		return nil
	}

	if err := json.Unmarshal(env.Response, out); err != nil {
		return fmt.Errorf("vk: decode response payload: %w", err)
	}

	return nil
}
