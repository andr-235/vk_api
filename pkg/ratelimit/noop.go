package ratelimit

import "context"

// NoOpRateLimiter — rate limiter без ограничений.
type NoOpRateLimiter struct{}

// NewNoOpRateLimiter создаёт rate limiter без ограничений.
func NewNoOpRateLimiter() *NoOpRateLimiter {
	return &NoOpRateLimiter{}
}

// Wait всегда возвращает nil без задержки.
func (r *NoOpRateLimiter) Wait(ctx context.Context) error {
	return nil
}
