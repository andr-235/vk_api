package testkit_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"testing"

	"github.com/andr-235/vk_api/internal/testkit"
)

func TestMockCaller(t *testing.T) {
	t.Run("call func called", func(t *testing.T) {
		called := false
		mock := testkit.NewMockCaller(
			testkit.WithCallFunc(func(ctx context.Context, method string, params, out any) error {
				called = true
				return nil
			}),
		)

		_ = mock.Call(context.Background(), "test.method", nil, nil)
		if !called {
			t.Error("CallFunc was not called")
		}
	})

	t.Run("call func returns error", func(t *testing.T) {
		testErr := errors.New("test error")
		mock := testkit.NewMockCaller(
			testkit.WithCallFunc(func(ctx context.Context, method string, params, out any) error {
				return testErr
			}),
		)

		err := mock.Call(context.Background(), "test.method", nil, nil)
		if err != testErr {
			t.Errorf("Call() error = %v, want %v", err, testErr)
		}
	})

	t.Run("call func decodes response", func(t *testing.T) {
		type Response struct {
			ID int `json:"id"`
		}

		mock := testkit.NewMockCaller(
			testkit.WithCallFunc(func(ctx context.Context, method string, params, out any) error {
				resp := &Response{ID: 123}
				data, _ := json.Marshal(resp)
				return json.Unmarshal(data, out)
			}),
		)

		var got Response
		err := mock.Call(context.Background(), "test.method", nil, &got)
		if err != nil {
			t.Fatalf("Call() error = %v", err)
		}
		if got.ID != 123 {
			t.Errorf("ID = %d, want 123", got.ID)
		}
	})

	t.Run("nil call func", func(t *testing.T) {
		mock := &testkit.MockCaller{}
		err := mock.Call(context.Background(), "test.method", nil, nil)
		if err != nil {
			t.Errorf("Call() error = %v, want nil", err)
		}
	})

	t.Run("call with raw handler func called", func(t *testing.T) {
		called := false
		mock := testkit.NewMockCaller(
			testkit.WithCallWithRawHandlerFunc(func(ctx context.Context, method string, params any, handler func(json.RawMessage) error) error {
				called = true
				return handler(json.RawMessage(`{"test": true}`))
			}),
		)

		var handlerCalled bool
		handler := func(raw json.RawMessage) error {
			handlerCalled = true
			return nil
		}

		_ = mock.CallWithRawHandler(context.Background(), "test.method", nil, handler)
		if !called {
			t.Error("CallWithRawHandlerFunc was not called")
		}
		if !handlerCalled {
			t.Error("Handler was not called")
		}
	})

	t.Run("nil call with raw handler func", func(t *testing.T) {
		mock := &testkit.MockCaller{}
		handler := func(raw json.RawMessage) error { return nil }
		err := mock.CallWithRawHandler(context.Background(), "test.method", nil, handler)
		if err != nil {
			t.Errorf("CallWithRawHandler() error = %v, want nil", err)
		}
	})
}

func TestMockDoer(t *testing.T) {
	t.Run("do func called", func(t *testing.T) {
		called := false
		mock := testkit.NewMockDoer(func(req *http.Request) (*http.Response, error) {
			called = true
			return nil, nil
		})

		_, _ = mock.Do(nil)
		if !called {
			t.Error("DoFunc was not called")
		}
	})

	t.Run("do func returns response", func(t *testing.T) {
		mock := testkit.NewMockDoer(func(req *http.Request) (*http.Response, error) {
			return nil, nil
		})

		resp, err := mock.Do(nil)
		if err != nil {
			t.Errorf("Do() error = %v", err)
		}
		if resp != nil {
			t.Errorf("resp = %v, want nil", resp)
		}
	})

	t.Run("do func returns error", func(t *testing.T) {
		testErr := errors.New("test error")
		mock := testkit.NewMockDoer(func(req *http.Request) (*http.Response, error) {
			return nil, testErr
		})

		resp, err := mock.Do(nil)
		if err != testErr {
			t.Errorf("Do() error = %v, want %v", err, testErr)
		}
		if resp != nil {
			t.Errorf("resp = %v, want nil", resp)
		}
	})
}

func TestMockRateLimiter(t *testing.T) {
	t.Run("wait func called", func(t *testing.T) {
		called := false
		mock := testkit.NewMockRateLimiter(func(ctx context.Context) error {
			called = true
			return nil
		})

		err := mock.Wait(context.Background())
		if err != nil {
			t.Errorf("Wait() error = %v", err)
		}
		if !called {
			t.Error("WaitFunc was not called")
		}
	})

	t.Run("wait func returns error", func(t *testing.T) {
		testErr := errors.New("rate limit exceeded")
		mock := testkit.NewMockRateLimiter(func(ctx context.Context) error {
			return testErr
		})

		err := mock.Wait(context.Background())
		if err != testErr {
			t.Errorf("Wait() error = %v, want %v", err, testErr)
		}
	})

	t.Run("nil wait func", func(t *testing.T) {
		mock := &testkit.MockRateLimiter{}
		err := mock.Wait(context.Background())
		if err != nil {
			t.Errorf("Wait() error = %v, want nil", err)
		}
	})
}

func TestResponder(t *testing.T) {
	t.Run("responder with data", func(t *testing.T) {
		type Response struct {
			ID int `json:"id"`
		}

		resp := testkit.NewResponder(&Response{ID: 42}, nil)
		callFunc := resp.ToCallFunc()

		var got Response
		err := callFunc(context.Background(), "test.method", nil, &got)
		if err != nil {
			t.Fatalf("ToCallFunc() error = %v", err)
		}
		if got.ID != 42 {
			t.Errorf("ID = %d, want 42", got.ID)
		}
	})

	t.Run("responder with error", func(t *testing.T) {
		testErr := errors.New("responder error")
		resp := testkit.NewResponder(nil, testErr)
		callFunc := resp.ToCallFunc()

		err := callFunc(context.Background(), "test.method", nil, nil)
		if err != testErr {
			t.Errorf("ToCallFunc() error = %v, want %v", err, testErr)
		}
	})

	t.Run("responder with nil data", func(t *testing.T) {
		resp := testkit.NewResponder(nil, nil)
		callFunc := resp.ToCallFunc()

		err := callFunc(context.Background(), "test.method", nil, nil)
		if err != nil {
			t.Errorf("ToCallFunc() error = %v, want nil", err)
		}
	})
}
