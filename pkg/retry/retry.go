package retry

import (
	"context"
)

// Retryer выполняет запрос с повторными попытками.
type Retryer interface {
	// Execute выполняет функцию с повторными попытками согласно политике.
	Execute(ctx context.Context, fn func() error) error
}

// RetryPolicy определяет политику повторных попыток.
type RetryPolicy interface {
	// ShouldRetry возвращает true, если запрос следует повторить.
	ShouldRetry(err error, attempt int) bool
	// Backoff возвращает задержку перед следующей попыткой.
	Backoff(attempt int) int
}
