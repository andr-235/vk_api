package vk

import (
	"context"
	"math"
	"math/rand"
	"time"
)

// RetryPolicy определяет политику повторных попыток.
type RetryPolicy interface {
	// ShouldRetry возвращает true, если запрос следует повторить.
	ShouldRetry(err error, attempt int) bool
	// Backoff возвращает задержку перед следующей попыткой.
	Backoff(attempt int) time.Duration
}

// Retryer выполняет запрос с повторными попытками.
type Retryer interface {
	Execute(ctx context.Context, fn func() error) error
}

// ExponentialBackoff реализует экспоненциальную задержку с jitter.
type ExponentialBackoff struct {
	Initial time.Duration
	Max     time.Duration
	Multiplier float64
	Jitter    float64
}

// DefaultRetryPolicy возвращает политику по умолчанию.
func DefaultRetryPolicy() *ExponentialBackoff {
	return &ExponentialBackoff{
		Initial:    100 * time.Millisecond,
		Max:        10 * time.Second,
		Multiplier: 2.0,
		Jitter:     0.1,
	}
}

// ShouldRetry определяет, стоит ли повторять запрос.
func (b *ExponentialBackoff) ShouldRetry(err error, attempt int) bool {
	if err == nil {
		return false
	}
	
	// Повторяем только при временных ошибках
	return isTemporaryError(err) && attempt < 5
}

// Backoff вычисляет задержку с экспоненциальным ростом и jitter.
func (b *ExponentialBackoff) Backoff(attempt int) time.Duration {
	if attempt <= 0 {
		return 0
	}

	// Экспоненциальный рост
	delay := float64(b.Initial) * math.Pow(b.Multiplier, float64(attempt-1))
	if delay > float64(b.Max) {
		delay = float64(b.Max)
	}

	// Добавляем jitter (случайную вариацию)
	if b.Jitter > 0 {
		jitterRange := delay * b.Jitter
		jitter := (rand.Float64() * 2 - 1) * jitterRange
		delay += jitter
	}

	return time.Duration(delay)
}

// isTemporaryError возвращает true для временных ошибок.
func isTemporaryError(err error) bool {
	if err == nil {
		return false
	}

	// Сетевые ошибки
	if isNetworkError(err) {
		return true
	}

	// HTTP 5xx ошибки
	if httpErr, ok := err.(*HTTPError); ok {
		return httpErr.StatusCode >= 500 && httpErr.StatusCode < 600
	}

	// Rate limit с Retry-After
	if rateErr, ok := err.(*RateLimitError); ok {
		return rateErr.RetryAfter > 0
	}

	return false
}

// isNetworkError определяет сетевые ошибки.
func isNetworkError(err error) bool {
	if err == nil {
		return false
	}

	// Проверяем на контекстные ошибки
	select {
	case <-context.Background().Done():
		return false
	default:
	}

	// Базовая проверка на сетевые ошибки
	errStr := err.Error()
	return containsAny(errStr, []string{
		"timeout",
		"connection refused",
		"connection reset",
		"no such host",
		"network is unreachable",
		"i/o timeout",
	})
}

func containsAny(s string, substrs []string) bool {
	for _, sub := range substrs {
		if contains(s, sub) {
			return true
		}
	}
	return false
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && findSubstring(s, substr))
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

// SimpleRetryer реализует простой retryer с фиксированным числом попыток.
type SimpleRetryer struct {
	MaxRetries int
	Policy     RetryPolicy
}

// NewSimpleRetryer создаёт новый retryer.
func NewSimpleRetryer(maxRetries int, policy RetryPolicy) *SimpleRetryer {
	return &SimpleRetryer{
		MaxRetries: maxRetries,
		Policy:     policy,
	}
}

// Execute выполняет функцию с повторными попытками.
func (r *SimpleRetryer) Execute(ctx context.Context, fn func() error) error {
	var lastErr error

	for attempt := 0; attempt <= r.MaxRetries; attempt++ {
		if err := ctx.Err(); err != nil {
			return err
		}

		lastErr = fn()
		if lastErr == nil {
			return nil
		}

		if !r.Policy.ShouldRetry(lastErr, attempt+1) {
			return lastErr
		}

		// Ждём перед следующей попыткой
		backoff := r.Policy.Backoff(attempt + 1)
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(backoff):
		}
	}

	return lastErr
}
