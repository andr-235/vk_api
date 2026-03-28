package vk

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestGroupsGetByID(t *testing.T) {
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
					"members_count": 1000
				}
			]
		}`))
	}))
	defer srv.Close()

	c := New(WithBaseURL(srv.URL))

	groups, err := c.GroupsGetByID(context.Background(), GroupsGetByIDParams{
		GroupIDs: []string{"vk_test"},
		Fields:   []string{"members_count"},
	})
	if err != nil {
		t.Fatalf("GroupsGetByID() error = %v", err)
	}

	if gotPath != "/groups.getById" {
		t.Fatalf("unexpected path: %s", gotPath)
	}
	if gotForm.Get("group_ids") != "vk_test" {
		t.Fatalf("unexpected group_ids: %q", gotForm.Get("group_ids"))
	}
	if gotForm.Get("fields") != "members_count" {
		t.Fatalf("unexpected fields: %q", gotForm.Get("fields"))
	}

	if len(groups) != 1 {
		t.Fatalf("unexpected groups len: %d", len(groups))
	}
	if groups[0].Name != "VK Test Group" {
		t.Fatalf("unexpected name: %q", groups[0].Name)
	}
	if groups[0].MembersCount != 1000 {
		t.Fatalf("unexpected members_count: %d", groups[0].MembersCount)
	}
}
