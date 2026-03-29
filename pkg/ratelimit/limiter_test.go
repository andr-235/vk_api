package ratelimit_test

import (
	"context"
	"testing"
	"time"

	"github.com/andr-235/vk_api/pkg/ratelimit"
)

func TestNoOpRateLimiter(t *testing.T) {
	limiter := ratelimit.NewNoOpRateLimiter()

	ctx := context.Background()
	err := limiter.Wait(ctx)
	if err != nil {
		t.Errorf("Wait() error = %v, want nil", err)
	}
}

func TestTokenBucketRateLimiter(t *testing.T) {
	t.Run("allows requests within limit", func(t *testing.T) {
		limiter := ratelimit.NewTokenBucketRateLimiter(10.0) // 10 requests/sec

		ctx := context.Background()
		for i := 0; i < 5; i++ {
			err := limiter.Wait(ctx)
			if err != nil {
				t.Fatalf("Wait() error = %v", err)
			}
		}
	})

	t.Run("rate limits excess requests", func(t *testing.T) {
		limiter := ratelimit.NewTokenBucketRateLimiter(10.0) // 10 requests/sec

		ctx := context.Background()
		// Exhaust tokens
		for i := 0; i < 10; i++ {
			_ = limiter.Wait(ctx)
		}

		// Next request should wait
		start := time.Now()
		err := limiter.Wait(ctx)
		elapsed := time.Since(start)

		if err != nil {
			t.Fatalf("Wait() error = %v", err)
		}
		if elapsed < 50*time.Millisecond {
			t.Errorf("Wait() should have waited, elapsed = %v", elapsed)
		}
	})

	t.Run("context canceled", func(t *testing.T) {
		limiter := ratelimit.NewTokenBucketRateLimiter(1.0) // 1 request/sec

		// Exhaust tokens
		_ = limiter.Wait(context.Background())

		ctx, cancel := context.WithCancel(context.Background())
		cancel() // Cancel immediately

		err := limiter.Wait(ctx)
		if err != context.Canceled {
			t.Errorf("Wait() error = %v, want %v", err, context.Canceled)
		}
	})

	t.Run("context deadline exceeded", func(t *testing.T) {
		limiter := ratelimit.NewTokenBucketRateLimiter(0.5) // 0.5 requests/sec

		// Exhaust tokens
		_ = limiter.Wait(context.Background())

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
		defer cancel()

		err := limiter.Wait(ctx)
		if err != context.DeadlineExceeded {
			t.Errorf("Wait() error = %v, want %v", err, context.DeadlineExceeded)
		}
	})

	t.Run("tokens refill over time", func(t *testing.T) {
		limiter := ratelimit.NewTokenBucketRateLimiter(10.0) // 10 requests/sec

		ctx := context.Background()
		// Exhaust tokens
		for i := 0; i < 10; i++ {
			_ = limiter.Wait(ctx)
		}

		// Wait for tokens to refill
		time.Sleep(200 * time.Millisecond)

		// Should be able to make requests again
		start := time.Now()
		for i := 0; i < 2; i++ {
			err := limiter.Wait(ctx)
			if err != nil {
				t.Fatalf("Wait() error = %v", err)
			}
		}
		elapsed := time.Since(start)

		if elapsed > 50*time.Millisecond {
			t.Errorf("Should not have waited, elapsed = %v", elapsed)
		}
	})
}

func BenchmarkNoOpRateLimiter(b *testing.B) {
	limiter := ratelimit.NewNoOpRateLimiter()
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = limiter.Wait(ctx)
	}
}

func BenchmarkTokenBucketRateLimiter(b *testing.B) {
	limiter := ratelimit.NewTokenBucketRateLimiter(1000.0) // High rate for benchmark
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = limiter.Wait(ctx)
	}
}
