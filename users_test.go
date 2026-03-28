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

func TestUsersGetFollowers(t *testing.T) {
	var gotPath string
	var gotForm url.Values

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		gotForm = r.URL.Query()

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
			"response": {
				"count": 2,
				"items": [
					{
						"id": 10,
						"first_name": "Ivan",
						"last_name": "Ivanov",
						"screen_name": "ivanov",
						"can_access_closed": true,
						"is_closed": false
					},
					{
						"id": 20,
						"first_name": "Petr",
						"last_name": "Petrov",
						"photo_100": "https://example.com/p100.jpg",
						"can_access_closed": true,
						"is_closed": false
					}
				]
			}
		}`))
	}))
	defer srv.Close()

	c := New(WithBaseURL(srv.URL))

	resp, err := c.UsersGetFollowers(context.Background(), UsersGetFollowersParams{
		UserID:   1,
		Offset:   5,
		Count:    2,
		Fields:   []string{"screen_name", "photo_100"},
		NameCase: "nom",
	})
	if err != nil {
		t.Fatalf("UsersGetFollowers() error = %v", err)
	}

	if gotPath != "/users.getFollowers" {
		t.Fatalf("unexpected path: %s", gotPath)
	}
	if gotForm.Get("user_id") != "1" {
		t.Fatalf("unexpected user_id: %q", gotForm.Get("user_id"))
	}
	if gotForm.Get("offset") != "5" {
		t.Fatalf("unexpected offset: %q", gotForm.Get("offset"))
	}
	if gotForm.Get("count") != "2" {
		t.Fatalf("unexpected count: %q", gotForm.Get("count"))
	}
	if gotForm.Get("fields") != "screen_name,photo_100" {
		t.Fatalf("unexpected fields: %q", gotForm.Get("fields"))
	}
	if gotForm.Get("name_case") != "nom" {
		t.Fatalf("unexpected name_case: %q", gotForm.Get("name_case"))
	}

	if resp == nil {
		t.Fatal("expected response, got nil")
	}
	if resp.Count != 2 {
		t.Fatalf("unexpected count: %d", resp.Count)
	}
	if len(resp.Items) != 2 {
		t.Fatalf("unexpected items len: %d", len(resp.Items))
	}
	if resp.Items[0].ID != 10 {
		t.Fatalf("unexpected first item id: %d", resp.Items[0].ID)
	}
	if resp.Items[0].ScreenName != "ivanov" {
		t.Fatalf("unexpected first item screen_name: %q", resp.Items[0].ScreenName)
	}
	if resp.Items[1].Photo100 != "https://example.com/p100.jpg" {
		t.Fatalf("unexpected second item photo_100: %q", resp.Items[1].Photo100)
	}
}

func TestUsersGetSubscriptions(t *testing.T) {
	var gotPath string
	var gotForm url.Values

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		gotForm = r.URL.Query()

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
			"response": {
				"users": {
					"count": 2,
					"items": [10, 20]
				},
				"groups": {
					"count": 1,
					"items": [30]
				}
			}
		}`))
	}))
	defer srv.Close()

	c := New(WithBaseURL(srv.URL))

	resp, err := c.UsersGetSubscriptions(context.Background(), UsersGetSubscriptionsParams{
		UserID: 1,
	})
	if err != nil {
		t.Fatalf("UsersGetSubscriptions() error = %v", err)
	}

	if gotPath != "/users.getSubscriptions" {
		t.Fatalf("unexpected path: %s", gotPath)
	}
	if gotForm.Get("user_id") != "1" {
		t.Fatalf("unexpected user_id: %q", gotForm.Get("user_id"))
	}
	if gotForm.Get("extended") != "" && gotForm.Get("extended") != "0" {
		t.Fatalf("unexpected extended: %q", gotForm.Get("extended"))
	}

	if resp == nil {
		t.Fatal("expected response, got nil")
	}
	if resp.Users.Count != 2 {
		t.Fatalf("unexpected users count: %d", resp.Users.Count)
	}
	if len(resp.Users.Items) != 2 {
		t.Fatalf("unexpected users items len: %d", len(resp.Users.Items))
	}
	if resp.Groups.Count != 1 {
		t.Fatalf("unexpected groups count: %d", resp.Groups.Count)
	}
	if len(resp.Groups.Items) != 1 {
		t.Fatalf("unexpected groups items len: %d", len(resp.Groups.Items))
	}
	if resp.Groups.Items[0] != 30 {
		t.Fatalf("unexpected first group id: %d", resp.Groups.Items[0])
	}
}

func TestUsersGetSubscriptionsExtended(t *testing.T) {
	var gotPath string
	var gotForm url.Values

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		gotForm = r.URL.Query()

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
			"response": {
				"count": 2,
				"items": [
					{
						"type": "profile",
						"id": 10,
						"first_name": "Ivan",
						"last_name": "Ivanov",
						"screen_name": "ivanov",
						"can_access_closed": true,
						"is_closed": false
					},
					{
						"type": "page",
						"id": 20,
						"name": "VK Test Page",
						"screen_name": "vk_test_page",
						"type": "page"
					}
				]
			}
		}`))
	}))
	defer srv.Close()

	c := New(WithBaseURL(srv.URL))

	resp, err := c.UsersGetSubscriptionsExtended(context.Background(), UsersGetSubscriptionsParams{
		UserID: 1,
		Offset: 5,
		Count:  2,
		Fields: []string{"screen_name"},
	})
	if err != nil {
		t.Fatalf("UsersGetSubscriptionsExtended() error = %v", err)
	}

	if gotPath != "/users.getSubscriptions" {
		t.Fatalf("unexpected path: %s", gotPath)
	}
	if gotForm.Get("user_id") != "1" {
		t.Fatalf("unexpected user_id: %q", gotForm.Get("user_id"))
	}
	if gotForm.Get("extended") != "1" {
		t.Fatalf("unexpected extended: %q", gotForm.Get("extended"))
	}
	if gotForm.Get("offset") != "5" {
		t.Fatalf("unexpected offset: %q", gotForm.Get("offset"))
	}
	if gotForm.Get("count") != "2" {
		t.Fatalf("unexpected count: %q", gotForm.Get("count"))
	}
	if gotForm.Get("fields") != "screen_name" {
		t.Fatalf("unexpected fields: %q", gotForm.Get("fields"))
	}

	if resp == nil {
		t.Fatal("expected response, got nil")
	}
	if resp.Count != 2 {
		t.Fatalf("unexpected count: %d", resp.Count)
	}
	if len(resp.Items) != 2 {
		t.Fatalf("unexpected items len: %d", len(resp.Items))
	}
	if resp.Items[0].FirstName != "Ivan" {
		t.Fatalf("unexpected first item first_name: %q", resp.Items[0].FirstName)
	}
	if resp.Items[1].Name != "VK Test Page" {
		t.Fatalf("unexpected second item name: %q", resp.Items[1].Name)
	}
}
