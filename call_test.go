package vk

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestClientCall_Success(t *testing.T) {
	var gotPath string
	var gotMethod string
	var gotForm url.Values
	var gotAuth string

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		gotMethod = r.Method
		gotAuth = r.Header.Get("Authorization")
		gotForm = r.URL.Query()

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"response":[{"id":1,"first_name":"Pavel","last_name":"Durov"}]}`))
	}))
	defer srv.Close()

	c := New(
		WithToken("token123"),
		WithVersion("5.199"),
		WithLang("ru"),
		WithBaseURL(srv.URL),
	)

	var out []User
	err := c.Call(context.Background(), "users.get", UsersGetParams{
		UserIDs: []int{1, 2},
		Fields:  []string{"bdate", "photo_100"},
	}, &out)
	if err != nil {
		t.Fatalf("Call() error = %v", err)
	}

	if gotMethod != http.MethodPost {
		t.Fatalf("unexpected method: %s", gotMethod)
	}
	if gotPath != "/users.get" {
		t.Fatalf("unexpected path: %s", gotPath)
	}
	if gotAuth != "" {
		t.Fatalf("unexpected authorization header: %q", gotAuth)
	}
	if gotForm.Get("access_token") != "token123" {
		t.Fatalf("unexpected access_token: %q", gotForm.Get("access_token"))
	}
	if gotForm.Get("v") != "5.199" {
		t.Fatalf("unexpected version: %q", gotForm.Get("v"))
	}
	if gotForm.Get("lang") != "ru" {
		t.Fatalf("unexpected lang: %q", gotForm.Get("lang"))
	}
	if gotForm.Get("user_ids") != "1,2" {
		t.Fatalf("unexpected user_ids: %q", gotForm.Get("user_ids"))
	}
	if gotForm.Get("fields") != "bdate,photo_100" {
		t.Fatalf("unexpected fields: %q", gotForm.Get("fields"))
	}

	if len(out) != 1 {
		t.Fatalf("unexpected users len: %d", len(out))
	}
	if out[0].ID != 1 {
		t.Fatalf("unexpected user id: %d", out[0].ID)
	}
	if out[0].FirstName != "Pavel" {
		t.Fatalf("unexpected first_name: %q", out[0].FirstName)
	}
}

func TestClientCall_TokenInHeader(t *testing.T) {
	var gotAuth string
	var gotForm url.Values

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotAuth = r.Header.Get("Authorization")
		gotForm = r.URL.Query()

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"response":[{"id":1,"first_name":"Test","last_name":"User"}]}`))
	}))
	defer srv.Close()

	c := New(
		WithToken("secret-token"),
		WithTokenSource(TokenInHeader),
		WithBaseURL(srv.URL),
	)

	var out []User
	err := c.Call(context.Background(), "users.get", nil, &out)
	if err != nil {
		t.Fatalf("Call() error = %v", err)
	}

	if gotAuth != "Bearer secret-token" {
		t.Fatalf("unexpected authorization header: %q", gotAuth)
	}
	if gotForm.Get("access_token") != "" {
		t.Fatalf("access_token must not be in query when header auth is used, got %q", gotForm.Get("access_token"))
	}
	if gotForm.Get("v") != "5.199" {
		t.Fatalf("unexpected version: %q", gotForm.Get("v"))
	}
}

func TestClientCall_VKError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
			"error": {
				"error_code": 5,
				"error_msg": "User authorization failed: invalid access_token.",
				"request_params": [{"key":"method","value":"users.get"}]
			}
		}`))
	}))
	defer srv.Close()

	c := New(WithBaseURL(srv.URL))

	var out []User
	err := c.Call(context.Background(), "users.get", nil, &out)
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	vkErr, ok := err.(*VKError)
	if !ok {
		t.Fatalf("expected *VKError, got %T", err)
	}
	if !vkErr.IsAuth() {
		t.Fatal("expected auth error")
	}
	if vkErr.Code != ErrorCodeAuthFailed {
		t.Fatalf("unexpected error code: %d", vkErr.Code)
	}
	if len(vkErr.RequestParams) != 1 {
		t.Fatalf("unexpected request_params len: %d", len(vkErr.RequestParams))
	}
}

func TestClientCall_RateLimitError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
			"error": {
				"error_code": 6,
				"error_msg": "Too many requests per second."
			}
		}`))
	}))
	defer srv.Close()

	c := New(WithBaseURL(srv.URL))

	err := c.Call(context.Background(), "users.get", nil, nil)
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	vkErr, ok := AsVKError(err)
	if !ok {
		t.Fatalf("expected VKError, got %T", err)
	}
	if !vkErr.IsRateLimit() {
		t.Fatal("expected rate-limit error")
	}
}

func TestClientCall_HTTPError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "bad gateway", http.StatusBadGateway)
	}))
	defer srv.Close()

	c := New(WithBaseURL(srv.URL))

	err := c.Call(context.Background(), "users.get", nil, nil)
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if !strings.Contains(err.Error(), "unexpected http status 502") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestClientCall_TestMode(t *testing.T) {
	var gotForm url.Values

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotForm = r.URL.Query()

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"response":[{"id":1,"first_name":"Mode","last_name":"Test"}]}`))
	}))
	defer srv.Close()

	c := New(
		WithBaseURL(srv.URL),
		WithTestMode(true),
	)

	var out []User
	err := c.Call(context.Background(), "users.get", nil, &out)
	if err != nil {
		t.Fatalf("Call() error = %v", err)
	}

	if gotForm.Get("test_mode") != "1" {
		t.Fatalf("unexpected test_mode: %q", gotForm.Get("test_mode"))
	}
}
