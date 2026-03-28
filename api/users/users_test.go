package users

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	vk "github.com/andr-235/vk_api"
)

func TestGet(t *testing.T) {
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

	client := vk.New(vk.WithBaseURL(srv.URL))

	users, err := Get(context.Background(), client, GetParams{
		UserIDs: []string{"743784474"},
		Fields:  []string{"bdate"},
	})
	if err != nil {
		t.Fatalf("Get() error = %v", err)
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

func TestGet_WithScreenName(t *testing.T) {
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

	client := vk.New(vk.WithBaseURL(srv.URL))

	users, err := Get(context.Background(), client, GetParams{
		UserIDs: []string{"durov"},
		Fields:  []string{"screen_name"},
	})
	if err != nil {
		t.Fatalf("Get() error = %v", err)
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

func TestGetFollowers(t *testing.T) {
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

	client := vk.New(vk.WithBaseURL(srv.URL))

	resp, err := GetFollowers(context.Background(), client, GetFollowersParams{
		UserID:   1,
		Offset:   5,
		Count:    2,
		Fields:   []string{"screen_name", "photo_100"},
		NameCase: "nom",
	})
	if err != nil {
		t.Fatalf("GetFollowers() error = %v", err)
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

func TestGetSubscriptions(t *testing.T) {
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

	client := vk.New(vk.WithBaseURL(srv.URL))

	resp, err := GetSubscriptions(context.Background(), client, GetSubscriptionsParams{
		UserID: 1,
	})
	if err != nil {
		t.Fatalf("GetSubscriptions() error = %v", err)
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

func TestGetSubscriptionsExtended(t *testing.T) {
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
						"is_member": 1,
						"members_count": 123
					}
				]
			}
		}`))
	}))
	defer srv.Close()

	client := vk.New(vk.WithBaseURL(srv.URL))

	resp, err := GetSubscriptionsExtended(context.Background(), client, GetSubscriptionsParams{
		UserID: 1,
		Offset: 5,
		Count:  2,
		Fields: []string{"screen_name"},
	})
	if err != nil {
		t.Fatalf("GetSubscriptionsExtended() error = %v", err)
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
	if resp.Items[1].MembersCount != 123 {
		t.Fatalf("unexpected second item members_count: %d", resp.Items[1].MembersCount)
	}
}

func TestSearch(t *testing.T) {
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
						"first_name": "Vasya",
						"last_name": "Babich",
						"screen_name": "vasya",
						"online": 1,
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

	client := vk.New(vk.WithBaseURL(srv.URL))

	resp, err := Search(context.Background(), client, SearchParams{
		Q:        "Вася Бабич",
		Sort:     0,
		Offset:   10,
		Count:    2,
		Fields:   []string{"screen_name", "photo_100", "online"},
		City:     1,
		Country:  1,
		Sex:      2,
		AgeFrom:  18,
		AgeTo:    35,
		Online:   true,
		HasPhoto: true,
		FromList: []string{"friends", "subscriptions"},
	})
	if err != nil {
		t.Fatalf("Search() error = %v", err)
	}

	if gotPath != "/users.search" {
		t.Fatalf("unexpected path: %s", gotPath)
	}
	if gotForm.Get("q") != "Вася Бабич" {
		t.Fatalf("unexpected q: %q", gotForm.Get("q"))
	}
	if gotForm.Get("sort") != "0" {
		t.Fatalf("unexpected sort: %q", gotForm.Get("sort"))
	}
	if gotForm.Get("offset") != "10" {
		t.Fatalf("unexpected offset: %q", gotForm.Get("offset"))
	}
	if gotForm.Get("count") != "2" {
		t.Fatalf("unexpected count: %q", gotForm.Get("count"))
	}
	if gotForm.Get("fields") != "screen_name,photo_100,online" {
		t.Fatalf("unexpected fields: %q", gotForm.Get("fields"))
	}
	if gotForm.Get("city") != "1" {
		t.Fatalf("unexpected city: %q", gotForm.Get("city"))
	}
	if gotForm.Get("country") != "1" {
		t.Fatalf("unexpected country: %q", gotForm.Get("country"))
	}
	if gotForm.Get("sex") != "2" {
		t.Fatalf("unexpected sex: %q", gotForm.Get("sex"))
	}
	if gotForm.Get("age_from") != "18" {
		t.Fatalf("unexpected age_from: %q", gotForm.Get("age_from"))
	}
	if gotForm.Get("age_to") != "35" {
		t.Fatalf("unexpected age_to: %q", gotForm.Get("age_to"))
	}
	if gotForm.Get("online") != "1" {
		t.Fatalf("unexpected online: %q", gotForm.Get("online"))
	}
	if gotForm.Get("has_photo") != "1" {
		t.Fatalf("unexpected has_photo: %q", gotForm.Get("has_photo"))
	}
	if gotForm.Get("from_list") != "friends,subscriptions" {
		t.Fatalf("unexpected from_list: %q", gotForm.Get("from_list"))
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
	if resp.Items[0].FirstName != "Vasya" {
		t.Fatalf("unexpected first item first_name: %q", resp.Items[0].FirstName)
	}
	if resp.Items[0].Online != 1 {
		t.Fatalf("unexpected first item online: %d", resp.Items[0].Online)
	}
	if resp.Items[1].Photo100 != "https://example.com/p100.jpg" {
		t.Fatalf("unexpected second item photo_100: %q", resp.Items[1].Photo100)
	}
}
