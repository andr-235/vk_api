package vk

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
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

		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Fatalf("read body: %v", err)
		}

		gotForm, err = url.ParseQuery(string(body))
		if err != nil {
			t.Fatalf("parse body: %v", err)
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
