package vk

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestUsersGet(t *testing.T) {
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
			t.Fatalf("parse body: %v", err)
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
			"response": [
				{
					"id": 743784474,
					"first_name": "Персик",
					"last_name": "Рыжий",
					"bdate": "21.12.2000"
				}
			]
		}`))
	}))
	defer srv.Close()

	c := New(WithBaseURL(srv.URL))

	users, err := c.UsersGet(context.Background(), UsersGetParams{
		UserIDs: []int{743784474},
		Fields:  []string{"bdate"},
	})
	if err != nil {
		t.Fatalf("UsersGet() error = %v", err)
	}

	if gotPath != "/users.get" {
		t.Fatalf("unexpected path: %s", gotPath)
	}
	if gotForm.Get("user_ids") != "743784474" {
		t.Fatalf("unexpected user_ids: %q", gotForm.Get("user_ids"))
	}
	if gotForm.Get("fields") != "bdate" {
		t.Fatalf("unexpected fields: %q", gotForm.Get("fields"))
	}
	if len(users) != 1 {
		t.Fatalf("unexpected users len: %d", len(users))
	}
	if users[0].FirstName != "Персик" {
		t.Fatalf("unexpected first_name: %q", users[0].FirstName)
	}
	if users[0].BDate != "21.12.2000" {
		t.Fatalf("unexpected bdate: %q", users[0].BDate)
	}
}
