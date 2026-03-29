package retry

import (
	"context"
	"errors"
	"math"
	"math/rand"
	"time"
)

// ExponentialBackoff реализует экспоненциальную задержку с jitter.
type ExponentialBackoff struct {
	Initial    time.Duration // Начальная задержка
	Max        time.Duration // Максимальная задержка
	Multiplier float64       // Множитель
	Jitter     float64       // Разброс (0.0-1.0)
}

// DefaultPolicy возвращает политику по умолчанию.
func DefaultPolicy() *ExponentialBackoff {
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
func (b *ExponentialBackoff) Backoff(attempt int) int {
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
		jitter := (rand.Float64()*2 - 1) * jitterRange
		delay += jitter
	}

	return int(delay)
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
	type httpError interface {
		Error() string
		StatusCode() int
	}
	if httpErr, ok := err.(httpError); ok {
		return httpErr.StatusCode() >= 500 && httpErr.StatusCode() < 600
	}

	// Rate limit с Retry-After
	type rateLimitError interface {
		Error() string
		RetryAfter() int
	}
	if rateErr, ok := err.(rateLimitError); ok {
		return rateErr.RetryAfter() > 0
	}

	return false
}

// isNetworkError определяет сетевые ошибки.
func isNetworkError(err error) bool {
	if err == nil {
		return false
	}

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
	return len(s) >= len(substr) && findSubstring(s, substr)
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

// IsTemporaryError проверяет, является ли ошибка временной.
func IsTemporaryError(err error) bool {
	return isTemporaryError(err)
}

// IsRetryableError проверяет, можно ли повторить запрос.
func IsRetryableError(err error) bool {
	if err == nil {
		return false
	}

	// Контекст отменён
	if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
		return false
	}

	return isTemporaryError(err)
}
