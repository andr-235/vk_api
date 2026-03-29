package transport_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/andr-235/vk_api/pkg/config"
	"github.com/andr-235/vk_api/pkg/transport"
)

func TestDefaultHTTPClient(t *testing.T) {
	hc := transport.DefaultHTTPClient()
	if hc == nil {
		t.Fatal("DefaultHTTPClient() returned nil")
	}
	if hc.Timeout != 30*time.Second {
		t.Errorf("Timeout = %v, want 30s", hc.Timeout)
	}
}

func TestTransportNew(t *testing.T) {
	cfg := config.DefaultConfig()
	cfg.Token = "test-token"

	hc := transport.DefaultHTTPClient()
	tr := transport.New(cfg, hc)

	if tr == nil {
		t.Fatal("New() returned nil")
	}
}

func TestTransportCall(t *testing.T) {
	t.Run("successful call", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"response": {"id": 123}}`))
		}))
		defer srv.Close()

		cfg := config.Config{
			BaseURL: srv.URL + "/",
			Version: "5.199",
		}
		tr := transport.New(cfg, transport.DefaultHTTPClient())

		var resp struct {
			ID int `json:"id"`
		}
		err := tr.Call(context.Background(), "test.method", url.Values{}, &resp)
		if err != nil {
			t.Fatalf("Call() error = %v", err)
		}
		if resp.ID != 123 {
			t.Errorf("ID = %d, want 123", resp.ID)
		}
	})

	t.Run("call with nil out", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"response": {}}`))
		}))
		defer srv.Close()

		cfg := config.Config{
			BaseURL: srv.URL + "/",
			Version: "5.199",
		}
		tr := transport.New(cfg, transport.DefaultHTTPClient())

		err := tr.Call(context.Background(), "test.method", url.Values{}, nil)
		if err != nil {
			t.Fatalf("Call() error = %v", err)
		}
	})

	t.Run("auth error", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"error": {"error_code": 5, "error_msg": "User authorization failed"}}`))
		}))
		defer srv.Close()

		cfg := config.Config{
			BaseURL: srv.URL + "/",
			Version: "5.199",
		}
		tr := transport.New(cfg, transport.DefaultHTTPClient())

		var resp interface{}
		err := tr.Call(context.Background(), "test.method", url.Values{}, &resp)
		if err == nil {
			t.Fatal("Call() should return error")
		}

		var authErr *transport.AuthError
		if !errors.As(err, &authErr) {
			t.Errorf("Expected AuthError, got %T: %v", err, err)
		}
	})

	t.Run("captcha error", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"error": {"error_code": 14, "error_msg": "Captcha needed", "captcha_sid": "abc123", "captcha_img": "http://example.com/captcha.png"}}`))
		}))
		defer srv.Close()

		cfg := config.Config{
			BaseURL: srv.URL + "/",
			Version: "5.199",
		}
		tr := transport.New(cfg, transport.DefaultHTTPClient())

		var resp interface{}
		err := tr.Call(context.Background(), "test.method", url.Values{}, &resp)
		if err == nil {
			t.Fatal("Call() should return error")
		}

		var captchaErr *transport.CaptchaError
		if !errors.As(err, &captchaErr) {
			t.Errorf("Expected CaptchaError, got %T: %v", err, err)
		}
	})

	t.Run("rate limit error", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"error": {"error_code": 6, "error_msg": "Too many requests"}}`))
		}))
		defer srv.Close()

		cfg := config.Config{
			BaseURL: srv.URL + "/",
			Version: "5.199",
		}
		tr := transport.New(cfg, transport.DefaultHTTPClient())

		var resp interface{}
		err := tr.Call(context.Background(), "test.method", url.Values{}, &resp)
		if err == nil {
			t.Fatal("Call() should return error")
		}

		var rateErr *transport.RateLimitError
		if !errors.As(err, &rateErr) {
			t.Errorf("Expected RateLimitError, got %T: %v", err, err)
		}
	})

	t.Run("permission error", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"error": {"error_code": 7, "error_msg": "Permission denied"}}`))
		}))
		defer srv.Close()

		cfg := config.Config{
			BaseURL: srv.URL + "/",
			Version: "5.199",
		}
		tr := transport.New(cfg, transport.DefaultHTTPClient())

		var resp interface{}
		err := tr.Call(context.Background(), "test.method", url.Values{}, &resp)
		if err == nil {
			t.Fatal("Call() should return error")
		}

		var permErr *transport.PermissionError
		if !errors.As(err, &permErr) {
			t.Errorf("Expected PermissionError, got %T: %v", err, err)
		}
	})
}

func TestTransportCallRaw(t *testing.T) {
	t.Run("handler called", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"response": {"key": "value"}}`))
		}))
		defer srv.Close()

		cfg := config.Config{
			BaseURL: srv.URL + "/",
			Version: "5.199",
		}
		tr := transport.New(cfg, transport.DefaultHTTPClient())

		var gotRaw json.RawMessage
		handler := func(raw json.RawMessage) error {
			gotRaw = raw
			return nil
		}

		err := tr.CallRaw(context.Background(), "test.method", url.Values{}, handler)
		if err != nil {
			t.Fatalf("CallRaw() error = %v", err)
		}
		if len(gotRaw) == 0 {
			t.Error("Handler was not called")
		}
	})
}

func TestHTTPError(t *testing.T) {
	err := &transport.HTTPError{
		StatusCode: 500,
		Body:       "Internal Server Error",
	}

	expected := "http error 500: Internal Server Error"
	if err.Error() != expected {
		t.Errorf("Error() = %q, want %q", err.Error(), expected)
	}
}

func TestVKError(t *testing.T) {
	err := &transport.VKError{
		Code:    1,
		Message: "Unknown error occurred",
	}

	expected := "vk api error 1: Unknown error occurred"
	if err.Error() != expected {
		t.Errorf("Error() = %q, want %q", err.Error(), expected)
	}
}

func TestAuthError(t *testing.T) {
	err := &transport.AuthError{
		Code:    5,
		Message: "User authorization failed",
	}

	expected := "auth error 5: User authorization failed"
	if err.Error() != expected {
		t.Errorf("Error() = %q, want %q", err.Error(), expected)
	}
}

func TestCaptchaError(t *testing.T) {
	err := &transport.CaptchaError{
		Code:       14,
		Message:    "Captcha needed",
		CaptchaSID: "abc123",
		CaptchaImg: "http://example.com/captcha.png",
	}

	expected := "captcha required: Captcha needed"
	if err.Error() != expected {
		t.Errorf("Error() = %q, want %q", err.Error(), expected)
	}
}

func TestRateLimitError(t *testing.T) {
	t.Run("with retry after", func(t *testing.T) {
		err := &transport.RateLimitError{
			Code:       6,
			Message:    "Too many requests",
			RetryAfter: 10,
		}

		expected := "rate limit exceeded, retry after 10 seconds"
		if err.Error() != expected {
			t.Errorf("Error() = %q, want %q", err.Error(), expected)
		}
	})

	t.Run("without retry after", func(t *testing.T) {
		err := &transport.RateLimitError{
			Code:    6,
			Message: "Too many requests",
		}

		expected := "rate limit exceeded: Too many requests"
		if err.Error() != expected {
			t.Errorf("Error() = %q, want %q", err.Error(), expected)
		}
	})
}

func TestPermissionError(t *testing.T) {
	err := &transport.PermissionError{
		Code:    7,
		Message: "Permission denied",
	}

	expected := "permission denied: Permission denied"
	if err.Error() != expected {
		t.Errorf("Error() = %q, want %q", err.Error(), expected)
	}
}
