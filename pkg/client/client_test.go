package client_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/andr-235/vk_api/pkg/client"
	"github.com/andr-235/vk_api/pkg/config"
	"github.com/andr-235/vk_api/pkg/middleware"
	"github.com/andr-235/vk_api/pkg/ratelimit"
	"github.com/andr-235/vk_api/pkg/retry"
)

func TestListResponse(t *testing.T) {
	type Item struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	}

	jsonData := `{"count": 2, "items": [{"id": 1, "name": "First"}, {"id": 2, "name": "Second"}]}`

	var resp client.ListResponse[Item]
	err := json.Unmarshal([]byte(jsonData), &resp)
	if err != nil {
		t.Fatalf("Unmarshal error = %v", err)
	}

	if resp.Count != 2 {
		t.Errorf("Count = %d, want 2", resp.Count)
	}
	if len(resp.Items) != 2 {
		t.Fatalf("Items len = %d, want 2", len(resp.Items))
	}
	if resp.Items[0].ID != 1 {
		t.Errorf("Items[0].ID = %d, want 1", resp.Items[0].ID)
	}
	if resp.Items[1].Name != "Second" {
		t.Errorf("Items[1].Name = %q, want %q", resp.Items[1].Name, "Second")
	}
}

func TestNew(t *testing.T) {
	t.Run("default config", func(t *testing.T) {
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("New() panicked: %v", r)
			}
		}()

		c := client.New(config.DefaultConfig())
		if c == nil {
			t.Error("New() returned nil")
		}
	})

	t.Run("with token", func(t *testing.T) {
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("New() panicked: %v", r)
			}
		}()

		c := client.New(config.DefaultConfig(), client.WithToken("test-token"))
		if c == nil {
			t.Error("New() returned nil")
		}
	})

	t.Run("invalid config panics", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("New() should panic with invalid config")
			}
		}()

		_ = client.New(config.Config{})
	})
}

func TestNewBuilder(t *testing.T) {
	t.Run("build with default config", func(t *testing.T) {
		c, err := client.NewBuilder().Build()
		if err != nil {
			t.Fatalf("Build() error = %v", err)
		}
		if c == nil {
			t.Error("Build() returned nil")
		}
	})

	t.Run("build with token", func(t *testing.T) {
		c, err := client.NewBuilder().
			WithToken("test-token").
			WithVersion("5.199").
			Build()
		if err != nil {
			t.Fatalf("Build() error = %v", err)
		}
		if c == nil {
			t.Error("Build() returned nil")
		}
	})

	t.Run("build with invalid config", func(t *testing.T) {
		_, err := client.NewBuilder().
			WithVersion(""). // пустая версия
			Build()
		if err == nil {
			t.Error("Build() should return error with invalid config")
		}
	})

	t.Run("must build panics on error", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("MustBuild() should panic on error")
			}
		}()

		_ = client.NewBuilder().WithVersion("").MustBuild()
	})
}

func TestClientCall(t *testing.T) {
	t.Run("successful call", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"response": {"count": 1, "items": [{"id": 1}]}}`))
		}))
		defer srv.Close()

		c := client.New(config.Config{
			BaseURL: srv.URL + "/",
			Version: "5.199",
		})

		type Response struct {
			Count int `json:"count"`
			Items []struct {
				ID int `json:"id"`
			} `json:"items"`
		}

		var resp Response
		err := c.Call(context.Background(), "test.method", nil, &resp)
		if err != nil {
			t.Fatalf("Call() error = %v", err)
		}
		if resp.Count != 1 {
			t.Errorf("Count = %d, want 1", resp.Count)
		}
	})

	t.Run("call with nil out", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"response": {}}`))
		}))
		defer srv.Close()

		c := client.New(config.Config{
			BaseURL: srv.URL + "/",
			Version: "5.199",
		})

		err := c.Call(context.Background(), "test.method", nil, nil)
		if err != nil {
			t.Fatalf("Call() error = %v", err)
		}
	})

	t.Run("call with error response", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"error": {"error_code": 5, "error_msg": "Auth failed"}}`))
		}))
		defer srv.Close()

		c := client.New(config.Config{
			BaseURL: srv.URL + "/",
			Version: "5.199",
		})

		var resp interface{}
		err := c.Call(context.Background(), "test.method", nil, &resp)
		if err == nil {
			t.Fatal("Call() should return error")
		}
	})
}

func TestClientCallWithRawHandler(t *testing.T) {
	t.Run("handler called", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"response": {"key": "value"}}`))
		}))
		defer srv.Close()

		c := client.New(config.Config{
			BaseURL: srv.URL + "/",
			Version: "5.199",
		})

		var gotRaw json.RawMessage
		handler := func(raw json.RawMessage) error {
			gotRaw = raw
			return nil
		}

		err := c.CallWithRawHandler(context.Background(), "test.method", nil, handler)
		if err != nil {
			t.Fatalf("CallWithRawHandler() error = %v", err)
		}
		if len(gotRaw) == 0 {
			t.Error("Handler was not called")
		}
	})

	t.Run("handler error", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"response": {}}`))
		}))
		defer srv.Close()

		c := client.New(config.Config{
			BaseURL: srv.URL + "/",
			Version: "5.199",
		})

		handler := func(raw json.RawMessage) error {
			return errors.New("handler error")
		}

		err := c.CallWithRawHandler(context.Background(), "test.method", nil, handler)
		if err == nil {
			t.Error("CallWithRawHandler() should return handler error")
		}
	})
}

func TestClientWithOptions(t *testing.T) {
	t.Run("with rate limiter", func(t *testing.T) {
		limiter := ratelimit.NewNoOpRateLimiter()
		c := client.New(config.DefaultConfig(), client.WithRateLimiter(limiter))
		if c == nil {
			t.Error("New() returned nil")
		}
	})

	t.Run("with retryer", func(t *testing.T) {
		retryer := retry.NewNoRetryer()
		c := client.New(config.DefaultConfig(), client.WithRetryer(retryer))
		if c == nil {
			t.Error("New() returned nil")
		}
	})

	t.Run("with interceptors", func(t *testing.T) {
		interceptor := middleware.NewLoggingInterceptor(nil)
		c := client.New(config.DefaultConfig(), client.WithInterceptors(interceptor))
		if c == nil {
			t.Error("New() returned nil")
		}
	})

	t.Run("with custom http client", func(t *testing.T) {
		httpClient := &http.Client{Timeout: 10 * time.Second}
		c := client.New(config.DefaultConfig(), client.WithHTTPClient(httpClient))
		if c == nil {
			t.Error("New() returned nil")
		}
	})
}

func TestClientConfig(t *testing.T) {
	cfg := config.Config{
		Token:   "test-token",
		Version: "5.199",
		Lang:    "ru",
	}
	c := client.New(cfg)

	got := c.Config()
	if got.Token != cfg.Token {
		t.Errorf("Config().Token = %q, want %q", got.Token, cfg.Token)
	}
	if got.Version != cfg.Version {
		t.Errorf("Config().Version = %q, want %q", got.Version, cfg.Version)
	}
	if got.Lang != cfg.Lang {
		t.Errorf("Config().Lang = %q, want %q", got.Lang, cfg.Lang)
	}
}
