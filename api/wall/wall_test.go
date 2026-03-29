package wall

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	vk "github.com/andr-235/vk_api"
)

func TestGet(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("expected method %s, got %s", http.MethodPost, r.Method)
		}

		if got := r.URL.Path; got != "/wall.get" {
			t.Fatalf("expected path %q, got %q", "/wall.get", got)
		}

		q := r.URL.Query()

		if got := q.Get("owner_id"); got != "1" {
			t.Fatalf("expected owner_id=1, got %q", got)
		}

		if got := q.Get("offset"); got != "5" {
			t.Fatalf("expected offset=5, got %q", got)
		}

		if got := q.Get("count"); got != "10" {
			t.Fatalf("expected count=10, got %q", got)
		}

		if got := q.Get("v"); got != "5.199" {
			t.Fatalf("expected v=5.199, got %q", got)
		}

		if got := q.Get("access_token"); got != "test-token" {
			t.Fatalf("expected access_token=test-token, got %q", got)
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
			"response": {
				"count": 1,
				"items": [
					{
						"id": 10,
						"owner_id": 1,
						"from_id": 2,
						"date": 1710000000,
						"text": "hello world"
					}
				]
			}
		}`))
	}))
	defer server.Close()

	client := vk.New(
		vk.WithBaseURL(server.URL+"/"),
		vk.WithToken("test-token"),
		vk.WithVersion("5.199"),
	)

	resp, err := Get(context.Background(), client, WallGetParams{
		OwnerID: 1,
		Offset:  5,
		Count:   10,
	})
	if err != nil {
		t.Fatalf("Get returned error: %v", err)
	}

	if resp == nil {
		t.Fatal("expected response, got nil")
	}

	if resp.Count != 1 {
		t.Fatalf("expected count=1, got %d", resp.Count)
	}

	if len(resp.Items) != 1 {
		t.Fatalf("expected 1 item, got %d", len(resp.Items))
	}

	post := resp.Items[0]

	if post.ID != 10 {
		t.Fatalf("expected id=10, got %d", post.ID)
	}

	if post.OwnerID != 1 {
		t.Fatalf("expected owner_id=1, got %d", post.OwnerID)
	}

	if post.FromID != 2 {
		t.Fatalf("expected from_id=2, got %d", post.FromID)
	}

	if post.Date != 1710000000 {
		t.Fatalf("expected date=1710000000, got %d", post.Date)
	}

	if post.Text != "hello world" {
		t.Fatalf("expected text=%q, got %q", "hello world", post.Text)
	}
}

func TestGet_APIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
			"error": {
				"error_code": 5,
				"error_msg": "User authorization failed"
			}
		}`))
	}))
	defer server.Close()

	client := vk.New(
		vk.WithBaseURL(server.URL+"/"),
		vk.WithToken("test-token"),
		vk.WithVersion("5.199"),
	)

	resp, err := Get(context.Background(), client, WallGetParams{
		OwnerID: 1,
	})

	if resp != nil {
		t.Fatalf("expected nil response, got %#v", resp)
	}

	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if err.Error() != "vk api error 5: User authorization failed" {
		t.Fatalf("unexpected error: %v", err)
	}
}
