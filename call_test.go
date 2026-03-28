package vk

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestClientCall_Success(t *testing.T) {
	var gotPath string
	var gotForm url.Values

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path

		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Fatalf("read body: %v", err)
		}

		gotForm, err = url.ParseQuery(string(body))
		if err != nil {
			t.Fatalf("parse form: %v", err)
		}

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
		Fields:  []string{"photo_100", "city"},
	}, &out)
	if err != nil {
		t.Fatalf("Call() error = %v", err)
	}

	if gotPath != "/users.get" {
		t.Fatalf("unexpected path: %s", gotPath)
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
	if gotForm.Get("fields") != "photo_100,city" {
		t.Fatalf("unexpected fields: %q", gotForm.Get("fields"))
	}

	if len(out) != 1 {
		t.Fatalf("unexpected users len: %d", len(out))
	}
	if out[0].FirstName != "Pavel" {
		t.Fatalf("unexpected first_name: %q", out[0].FirstName)
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
	if vkErr.Code != 5 {
		t.Fatalf("unexpected error code: %d", vkErr.Code)
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
