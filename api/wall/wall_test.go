package vk

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestWallGet(t *testing.T) {
	var gotPath string
	var gotForm url.Values

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		gotForm = r.URL.Query()

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
			"response": {
				"count": 1,
				"items": [
					{
						"id": 10,
						"owner_id": -1,
						"from_id": -1,
						"date": 1711111111,
						"text": "hello from wall"
					}
				]
			}
		}`))
	}))
	defer srv.Close()

	c := New(WithBaseURL(srv.URL))

	resp, err := c.WallGet(context.Background(), WallGetParams{
		OwnerID: -1,
		Count:   1,
	})
	if err != nil {
		t.Fatalf("WallGet() error = %v", err)
	}

	if gotPath != "/wall.get" {
		t.Fatalf("unexpected path: %s", gotPath)
	}
	if gotForm.Get("owner_id") != "-1" {
		t.Fatalf("unexpected owner_id: %q", gotForm.Get("owner_id"))
	}
	if gotForm.Get("count") != "1" {
		t.Fatalf("unexpected count: %q", gotForm.Get("count"))
	}
	if resp == nil {
		t.Fatal("expected response, got nil")
	}
	if resp.Count != 1 {
		t.Fatalf("unexpected count: %d", resp.Count)
	}
	if len(resp.Items) != 1 {
		t.Fatalf("unexpected items len: %d", len(resp.Items))
	}
	if resp.Items[0].Text != "hello from wall" {
		t.Fatalf("unexpected text: %q", resp.Items[0].Text)
	}
}
