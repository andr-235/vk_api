package vk

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestMessagesSend(t *testing.T) {
	var gotPath string
	var gotForm url.Values

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		gotForm = r.URL.Query()

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"response": 12345}`))
	}))
	defer srv.Close()

	c := New(WithBaseURL(srv.URL))

	messageID, err := c.MessagesSend(context.Background(), MessagesSendParams{
		UserID:   1,
		RandomID: 42,
		Message:  "hello",
	})
	if err != nil {
		t.Fatalf("MessagesSend() error = %v", err)
	}

	if gotPath != "/messages.send" {
		t.Fatalf("unexpected path: %s", gotPath)
	}
	if gotForm.Get("user_id") != "1" {
		t.Fatalf("unexpected user_id: %q", gotForm.Get("user_id"))
	}
	if gotForm.Get("random_id") != "42" {
		t.Fatalf("unexpected random_id: %q", gotForm.Get("random_id"))
	}
	if gotForm.Get("message") != "hello" {
		t.Fatalf("unexpected message: %q", gotForm.Get("message"))
	}
	if messageID != 12345 {
		t.Fatalf("unexpected message id: %d", messageID)
	}
}
