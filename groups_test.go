package vk

import (
	"context"
	"io"
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
					"id": 1,
					"name": "VK Test Group",
					"screen_name": "vk_test",
					"type": "group"
				}
			]
		}`))
	}))
	defer srv.Close()

	c := New(WithBaseURL(srv.URL))

	groups, err := c.GroupsGetByID(context.Background(), GroupsGetByIDParams{
		GroupIDs: []string{"vk_test"},
		Fields:   []string{"screen_name"},
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
	if gotForm.Get("fields") != "screen_name" {
		t.Fatalf("unexpected fields: %q", gotForm.Get("fields"))
	}
	if len(groups) != 1 {
		t.Fatalf("unexpected groups len: %d", len(groups))
	}
	if groups[0].ScreenName != "vk_test" {
		t.Fatalf("unexpected screen_name: %q", groups[0].ScreenName)
	}
}
