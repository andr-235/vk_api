package vk

import (
	"context"
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
		gotForm = r.URL.Query()

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
			"response": [
				{
					"id": 743784474,
					"first_name": "Персик",
					"last_name": "Рыжий",
					"bdate": "21.12.2000",
					"can_access_closed": true,
					"is_closed": false
				}
			]
		}`))
	}))
	defer srv.Close()

	c := New(WithBaseURL(srv.URL))

	users, err := c.UsersGet(context.Background(), UsersGetParams{
		UserIDs: []string{"743784474"},
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
	if !users[0].CanAccessClosed {
		t.Fatal("expected can_access_closed=true")
	}
	if users[0].IsClosed {
		t.Fatal("expected is_closed=false")
	}
}

func TestUsersGet_WithScreenName(t *testing.T) {
	var gotForm url.Values

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotForm = r.URL.Query()

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
			"response": [
				{
					"id": 1,
					"first_name": "Павел",
					"last_name": "Дуров",
					"screen_name": "durov",
					"can_access_closed": true,
					"is_closed": false
				}
			]
		}`))
	}))
	defer srv.Close()

	c := New(WithBaseURL(srv.URL))

	users, err := c.UsersGet(context.Background(), UsersGetParams{
		UserIDs: []string{"durov"},
		Fields:  []string{"screen_name"},
	})
	if err != nil {
		t.Fatalf("UsersGet() error = %v", err)
	}

	if gotForm.Get("user_ids") != "durov" {
		t.Fatalf("unexpected user_ids: %q", gotForm.Get("user_ids"))
	}
	if gotForm.Get("fields") != "screen_name" {
		t.Fatalf("unexpected fields: %q", gotForm.Get("fields"))
	}
	if len(users) != 1 {
		t.Fatalf("unexpected users len: %d", len(users))
	}
	if users[0].ScreenName != "durov" {
		t.Fatalf("unexpected screen_name: %q", users[0].ScreenName)
	}
}
