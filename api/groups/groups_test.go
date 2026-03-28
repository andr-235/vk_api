package groups

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	vk "github.com/andr-235/vk_api"
)

func TestGetByID(t *testing.T) {
	var gotPath string
	var gotForm url.Values

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		gotForm = r.URL.Query()

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
			"response": [
				{
					"id": 1,
					"name": "VK Test Group",
					"screen_name": "vk_test",
					"type": "group",
					"is_closed": 0,
					"members_count": 1000,
					"photo_100": "https://example.com/group_100.jpg",
					"description": "test group"
				}
			]
		}`))
	}))
	defer srv.Close()

	client := vk.New(vk.WithBaseURL(srv.URL))

	items, err := GetByID(context.Background(), client, GetByIDParams{
		GroupIDs: []string{"vk_test"},
		Fields:   []string{"members_count", "photo_100", "description"},
	})
	if err != nil {
		t.Fatalf("GetByID() error = %v", err)
	}

	if gotPath != "/groups.getById" {
		t.Fatalf("unexpected path: %s", gotPath)
	}
	if gotForm.Get("group_ids") != "vk_test" {
		t.Fatalf("unexpected group_ids: %q", gotForm.Get("group_ids"))
	}
	if gotForm.Get("fields") != "members_count,photo_100,description" {
		t.Fatalf("unexpected fields: %q", gotForm.Get("fields"))
	}

	if len(items) != 1 {
		t.Fatalf("unexpected groups len: %d", len(items))
	}
	if items[0].Name != "VK Test Group" {
		t.Fatalf("unexpected name: %q", items[0].Name)
	}
	if items[0].MembersCount != 1000 {
		t.Fatalf("unexpected members_count: %d", items[0].MembersCount)
	}
	if items[0].Photo100 != "https://example.com/group_100.jpg" {
		t.Fatalf("unexpected photo_100: %q", items[0].Photo100)
	}
	if items[0].Description != "test group" {
		t.Fatalf("unexpected description: %q", items[0].Description)
	}
}

func TestGetMembers(t *testing.T) {
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
						"photo_100": "https://example.com/u1.jpg",
						"online": 1,
						"can_access_closed": true,
						"is_closed": false
					},
					{
						"id": 20,
						"first_name": "Petr",
						"last_name": "Petrov",
						"sex": 2,
						"can_access_closed": true,
						"is_closed": false
					}
				]
			}
		}`))
	}))
	defer srv.Close()

	client := vk.New(vk.WithBaseURL(srv.URL))

	resp, err := GetMembers(context.Background(), client, GetMembersParams{
		GroupID: "vk_test",
		Sort:    "id_asc",
		Offset:  10,
		Count:   2,
		Fields:  []string{"photo_100", "online", "sex"},
		Filter:  "friends",
	})
	if err != nil {
		t.Fatalf("GetMembers() error = %v", err)
	}

	if gotPath != "/groups.getMembers" {
		t.Fatalf("unexpected path: %s", gotPath)
	}
	if gotForm.Get("group_id") != "vk_test" {
		t.Fatalf("unexpected group_id: %q", gotForm.Get("group_id"))
	}
	if gotForm.Get("sort") != "id_asc" {
		t.Fatalf("unexpected sort: %q", gotForm.Get("sort"))
	}
	if gotForm.Get("offset") != "10" {
		t.Fatalf("unexpected offset: %q", gotForm.Get("offset"))
	}
	if gotForm.Get("count") != "2" {
		t.Fatalf("unexpected count: %q", gotForm.Get("count"))
	}
	if gotForm.Get("fields") != "photo_100,online,sex" {
		t.Fatalf("unexpected fields: %q", gotForm.Get("fields"))
	}
	if gotForm.Get("filter") != "friends" {
		t.Fatalf("unexpected filter: %q", gotForm.Get("filter"))
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
		t.Fatalf("unexpected first_name: %q", resp.Items[0].FirstName)
	}
	if resp.Items[0].Photo100 != "https://example.com/u1.jpg" {
		t.Fatalf("unexpected photo_100: %q", resp.Items[0].Photo100)
	}
	if resp.Items[1].Sex != 2 {
		t.Fatalf("unexpected sex: %d", resp.Items[1].Sex)
	}
}
