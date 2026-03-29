package retry

import (
	"context"
	"time"
)

// SimpleRetryer реализует Retryer с фиксированным числом попыток.
type SimpleRetryer struct {
	MaxRetries int
	Policy     RetryPolicy
}

// NewSimpleRetryer создаёт новый retryer.
func NewSimpleRetryer(maxRetries int, policy RetryPolicy) *SimpleRetryer {
	if policy == nil {
		policy = DefaultPolicy()
	}
	return &SimpleRetryer{
		MaxRetries: maxRetries,
		Policy:     policy,
	}
}

// Execute выполняет функцию с повторными попытками.
func (r *SimpleRetryer) Execute(ctx context.Context, fn func() error) error {
	var lastErr error

	for attempt := 0; attempt <= r.MaxRetries; attempt++ {
		// Проверяем контекст перед попыткой
		if err := ctx.Err(); err != nil {
			return err
		}

		// Выполняем функцию
		lastErr = fn()
		if lastErr == nil {
			return nil
		}

		// Проверяем, стоит ли повторять
		if !r.Policy.ShouldRetry(lastErr, attempt+1) {
			return lastErr
		}

		// Ждём перед следующей попыткой
		backoffMs := r.Policy.Backoff(attempt + 1)
		if backoffMs > 0 {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(time.Duration(backoffMs) * time.Millisecond):
			}
		}
	}

	return lastErr
}

// NoRetryer — retryer без повторных попыток.
type NoRetryer struct{}

// NewNoRetryer создаёт retryer без повторных попыток.
func NewNoRetryer() *NoRetryer {
	return &NoRetryer{}
}

// Execute просто выполняет функцию без повторных попыток.
func (r *NoRetryer) Execute(ctx context.Context, fn func() error) error {
	return fn()
}
