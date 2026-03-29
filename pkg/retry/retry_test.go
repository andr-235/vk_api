package retry

import (
	"context"
	"errors"
	"net"
	"net/url"
	"testing"
	"time"
)

func BenchmarkExponentialBackoff(b *testing.B) {
	backoff := DefaultPolicy()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = backoff.Backoff(i % 5)
	}
}

func BenchmarkSimpleRetryer(b *testing.B) {
	retryer := NewSimpleRetryer(3, DefaultPolicy())
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = retryer.Execute(ctx, func() error {
			return nil
		})
	}
}

func TestIsTemporaryError(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected bool
	}{
		{
			name:     "nil error",
			err:      nil,
			expected: false,
		},
		{
			name:     "net.OpError",
			err:      &net.OpError{Op: "dial", Net: "tcp", Err: errors.New("connection refused")},
			expected: true,
		},
		{
			name:     "url.Error",
			err:      &url.Error{Op: "Get", URL: "http://example.com", Err: errors.New("timeout")},
			expected: true,
		},
		{
			name:     "context.Canceled",
			err:      context.Canceled,
			expected: false,
		},
		{
			name:     "context.DeadlineExceeded",
			err:      context.DeadlineExceeded,
			expected: true, // deadline exceeded — временная ошибка
		},
		{
			name:     "generic error",
			err:      errors.New("some error"),
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isTemporaryError(tt.err)
			if got != tt.expected {
				t.Errorf("isTemporaryError() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestIsRetryableError(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected bool
	}{
		{
			name:     "nil error",
			err:      nil,
			expected: false,
		},
		{
			name:     "net.OpError",
			err:      &net.OpError{Op: "dial", Net: "tcp", Err: errors.New("connection refused")},
			expected: true,
		},
		{
			name:     "context.Canceled",
			err:      context.Canceled,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsRetryableError(tt.err)
			if got != tt.expected {
				t.Errorf("IsRetryableError() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestSimpleRetryer(t *testing.T) {
	t.Run("success on first try", func(t *testing.T) {
		retryer := NewSimpleRetryer(3, DefaultPolicy())
		attempts := 0

		err := retryer.Execute(context.Background(), func() error {
			attempts++
			return nil
		})

		if err != nil {
			t.Errorf("Execute() error = %v", err)
		}
		if attempts != 1 {
			t.Errorf("attempts = %d, want 1", attempts)
		}
	})

	t.Run("no retry for non-temporary error", func(t *testing.T) {
		retryer := NewSimpleRetryer(3, DefaultPolicy())
		attempts := 0

		err := retryer.Execute(context.Background(), func() error {
			attempts++
			return errors.New("non-temporary error")
		})

		if err == nil {
			t.Error("Execute() error = nil, want error")
		}
		if attempts != 1 {
			t.Errorf("attempts = %d, want 1", attempts)
		}
	})

	t.Run("context canceled", func(t *testing.T) {
		retryer := NewSimpleRetryer(3, &ExponentialBackoff{
			Initial:    time.Second,
			Max:        time.Second,
			Multiplier: 1,
			Jitter:     0,
		})

		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		err := retryer.Execute(ctx, func() error {
			return &net.OpError{Op: "dial", Net: "tcp", Err: errors.New("connection refused")}
		})

		if err != context.Canceled {
			t.Errorf("Execute() error = %v, want %v", err, context.Canceled)
		}
	})
}

func TestNoRetryer(t *testing.T) {
	retryer := NewNoRetryer()
	attempts := 0

	err := retryer.Execute(context.Background(), func() error {
		attempts++
		return errors.New("error")
	})

	if err == nil {
		t.Error("Execute() error = nil, want error")
	}
	if attempts != 1 {
		t.Errorf("attempts = %d, want 1", attempts)
	}
}
