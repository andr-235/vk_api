package middleware_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/andr-235/vk_api/pkg/middleware"
)

func TestInterceptorChain(t *testing.T) {
	t.Run("empty chain", func(t *testing.T) {
		chain := middleware.InterceptorChain{}
		ctx := context.Background()
		req := &middleware.RequestContext{
			Method:    "test.method",
			RequestID: "123",
			StartTime: time.Now(),
		}

		got := chain.InterceptRequest(ctx, req)
		if got != ctx {
			t.Error("InterceptRequest should return same context")
		}

		resp := &middleware.ResponseContext{
			StatusCode: 200,
			Duration:   time.Millisecond,
		}
		got = chain.InterceptResponse(ctx, req, resp)
		if got != ctx {
			t.Error("InterceptResponse should return same context")
		}

		got = chain.InterceptError(ctx, req, errors.New("test error"))
		if got != ctx {
			t.Error("InterceptError should return same context")
		}
	})

	t.Run("single interceptor", func(t *testing.T) {
		called := false
		interceptor := middleware.NewInterceptorFunc(
			func(ctx context.Context, req *middleware.RequestContext) context.Context {
				called = true
				return ctx
			},
			nil,
			nil,
		)

		chain := middleware.InterceptorChain{interceptor}
		ctx := context.Background()
		req := &middleware.RequestContext{Method: "test"}

		_ = chain.InterceptRequest(ctx, req)
		if !called {
			t.Error("Interceptor was not called")
		}
	})

	t.Run("multiple interceptors", func(t *testing.T) {
		var calls []string

		interceptor1 := middleware.NewInterceptorFunc(
			func(ctx context.Context, req *middleware.RequestContext) context.Context {
				calls = append(calls, "req1")
				return ctx
			},
			func(ctx context.Context, req *middleware.RequestContext, resp *middleware.ResponseContext) context.Context {
				calls = append(calls, "resp1")
				return ctx
			},
			nil,
		)

		interceptor2 := middleware.NewInterceptorFunc(
			func(ctx context.Context, req *middleware.RequestContext) context.Context {
				calls = append(calls, "req2")
				return ctx
			},
			func(ctx context.Context, req *middleware.RequestContext, resp *middleware.ResponseContext) context.Context {
				calls = append(calls, "resp2")
				return ctx
			},
			nil,
		)

		chain := middleware.InterceptorChain{interceptor1, interceptor2}
		ctx := context.Background()
		req := &middleware.RequestContext{Method: "test"}
		resp := &middleware.ResponseContext{}

		_ = chain.InterceptRequest(ctx, req)
		if len(calls) != 2 || calls[0] != "req1" || calls[1] != "req2" {
			t.Errorf("calls = %v, want [req1, req2]", calls)
		}

		calls = nil
		_ = chain.InterceptResponse(ctx, req, resp)
		if len(calls) != 2 || calls[0] != "resp1" || calls[1] != "resp2" {
			t.Errorf("calls = %v, want [resp1, resp2]", calls)
		}
	})

	t.Run("panic recovery in request", func(t *testing.T) {
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("Should not panic: %v", r)
			}
		}()

		panicInterceptor := middleware.NewInterceptorFunc(
			func(ctx context.Context, req *middleware.RequestContext) context.Context {
				panic("test panic")
			},
			nil,
			nil,
		)

		chain := middleware.InterceptorChain{panicInterceptor}
		ctx := context.Background()
		req := &middleware.RequestContext{Method: "test"}

		_ = chain.InterceptRequest(ctx, req)
	})

	t.Run("panic recovery in response", func(t *testing.T) {
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("Should not panic: %v", r)
			}
		}()

		panicInterceptor := middleware.NewInterceptorFunc(
			nil,
			func(ctx context.Context, req *middleware.RequestContext, resp *middleware.ResponseContext) context.Context {
				panic("test panic")
			},
			nil,
		)

		chain := middleware.InterceptorChain{panicInterceptor}
		ctx := context.Background()
		req := &middleware.RequestContext{Method: "test"}
		resp := &middleware.ResponseContext{}

		_ = chain.InterceptResponse(ctx, req, resp)
	})

	t.Run("panic recovery in error", func(t *testing.T) {
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("Should not panic: %v", r)
			}
		}()

		panicInterceptor := middleware.NewInterceptorFunc(
			nil,
			nil,
			func(ctx context.Context, req *middleware.RequestContext, err error) context.Context {
				panic("test panic")
			},
		)

		chain := middleware.InterceptorChain{panicInterceptor}
		ctx := context.Background()
		req := &middleware.RequestContext{Method: "test"}

		_ = chain.InterceptError(ctx, req, errors.New("test"))
	})
}

func TestLoggingInterceptor(t *testing.T) {
	t.Run("logs request", func(t *testing.T) {
		var loggedMsg string
		var loggedFields map[string]any

		logFunc := func(ctx context.Context, msg string, fields map[string]any) {
			loggedMsg = msg
			loggedFields = fields
		}

		interceptor := middleware.NewLoggingInterceptor(logFunc)
		ctx := context.Background()
		req := &middleware.RequestContext{
			Method:    "users.get",
			RequestID: "123",
			Attempt:   1,
		}

		_ = interceptor.InterceptRequest(ctx, req)
		if loggedMsg != "VK request started" {
			t.Errorf("loggedMsg = %q, want %q", loggedMsg, "VK request started")
		}
		if loggedFields["request_id"] != "123" {
			t.Errorf("request_id = %v, want 123", loggedFields["request_id"])
		}
	})

	t.Run("logs response", func(t *testing.T) {
		var loggedMsg string
		var loggedFields map[string]any

		logFunc := func(ctx context.Context, msg string, fields map[string]any) {
			loggedMsg = msg
			loggedFields = fields
		}

		interceptor := middleware.NewLoggingInterceptor(logFunc)
		ctx := context.Background()
		req := &middleware.RequestContext{
			Method:    "users.get",
			RequestID: "123",
		}
		resp := &middleware.ResponseContext{
			StatusCode: 200,
			BodySize:   1024,
			Duration:   100 * time.Millisecond,
		}

		_ = interceptor.InterceptResponse(ctx, req, resp)
		if loggedMsg != "VK request completed" {
			t.Errorf("loggedMsg = %q, want %q", loggedMsg, "VK request completed")
		}
		if loggedFields["status_code"] != 200 {
			t.Errorf("status_code = %v, want 200", loggedFields["status_code"])
		}
	})

	t.Run("logs error", func(t *testing.T) {
		var loggedMsg string
		var loggedFields map[string]any

		logFunc := func(ctx context.Context, msg string, fields map[string]any) {
			loggedMsg = msg
			loggedFields = fields
		}

		interceptor := middleware.NewLoggingInterceptor(logFunc)
		ctx := context.Background()
		req := &middleware.RequestContext{
			Method:    "users.get",
			RequestID: "123",
		}
		testErr := errors.New("test error")

		_ = interceptor.InterceptError(ctx, req, testErr)
		if loggedMsg != "VK request failed" {
			t.Errorf("loggedMsg = %q, want %q", loggedMsg, "VK request failed")
		}
		if loggedFields["error"] != "test error" {
			t.Errorf("error = %v, want test error", loggedFields["error"])
		}
	})

	t.Run("nil log func", func(t *testing.T) {
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("Should not panic: %v", r)
			}
		}()

		interceptor := middleware.NewLoggingInterceptor(nil)
		ctx := context.Background()
		req := &middleware.RequestContext{Method: "test"}
		resp := &middleware.ResponseContext{}

		_ = interceptor.InterceptRequest(ctx, req)
		_ = interceptor.InterceptResponse(ctx, req, resp)
		_ = interceptor.InterceptError(ctx, req, errors.New("test"))
	})
}

func TestMetricsInterceptor(t *testing.T) {
	t.Run("records request", func(t *testing.T) {
		var recordedMethod string

		recordFunc := func(ctx context.Context, method string) {
			recordedMethod = method
		}

		interceptor := middleware.NewMetricsInterceptor(nil, nil, recordFunc)
		ctx := context.Background()
		req := &middleware.RequestContext{Method: "users.get"}

		_ = interceptor.InterceptRequest(ctx, req)
		if recordedMethod != "users.get" {
			t.Errorf("recordedMethod = %q, want %q", recordedMethod, "users.get")
		}
	})

	t.Run("records latency", func(t *testing.T) {
		var recordedMethod string
		var recordedStatusCode int
		var recordedDuration time.Duration

		recordFunc := func(ctx context.Context, method string, statusCode int, duration time.Duration) {
			recordedMethod = method
			recordedStatusCode = statusCode
			recordedDuration = duration
		}

		interceptor := middleware.NewMetricsInterceptor(recordFunc, nil, nil)
		ctx := context.Background()
		req := &middleware.RequestContext{Method: "users.get"}
		resp := &middleware.ResponseContext{
			StatusCode: 200,
			Duration:   100 * time.Millisecond,
		}

		_ = interceptor.InterceptResponse(ctx, req, resp)
		if recordedMethod != "users.get" {
			t.Errorf("recordedMethod = %q, want %q", recordedMethod, "users.get")
		}
		if recordedStatusCode != 200 {
			t.Errorf("recordedStatusCode = %d, want 200", recordedStatusCode)
		}
		if recordedDuration != 100*time.Millisecond {
			t.Errorf("recordedDuration = %v, want 100ms", recordedDuration)
		}
	})

	t.Run("records error", func(t *testing.T) {
		var recordedMethod string
		var recordedErr error

		recordFunc := func(ctx context.Context, method string, err error) {
			recordedMethod = method
			recordedErr = err
		}

		interceptor := middleware.NewMetricsInterceptor(nil, recordFunc, nil)
		ctx := context.Background()
		req := &middleware.RequestContext{Method: "users.get"}
		testErr := errors.New("test error")

		_ = interceptor.InterceptError(ctx, req, testErr)
		if recordedMethod != "users.get" {
			t.Errorf("recordedMethod = %q, want %q", recordedMethod, "users.get")
		}
		if recordedErr.Error() != "test error" {
			t.Errorf("recordedErr = %v, want test error", recordedErr)
		}
	})

	t.Run("nil funcs", func(t *testing.T) {
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("Should not panic: %v", r)
			}
		}()

		interceptor := middleware.NewMetricsInterceptor(nil, nil, nil)
		ctx := context.Background()
		req := &middleware.RequestContext{Method: "test"}
		resp := &middleware.ResponseContext{}

		_ = interceptor.InterceptRequest(ctx, req)
		_ = interceptor.InterceptResponse(ctx, req, resp)
		_ = interceptor.InterceptError(ctx, req, errors.New("test"))
	})
}
