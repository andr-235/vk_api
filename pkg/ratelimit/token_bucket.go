package ratelimit

import (
	"context"
	"sync"
	"time"
)

// TokenBucketRateLimiter реализует алгоритм token bucket для rate limiting.
// Потокобезопасен.
type TokenBucketRateLimiter struct {
	rate       float64 // токенов в секунду
	tokens     float64 // текущее количество токенов
	maxTokens  float64 // максимальное количество токенов
	lastUpdate time.Time
	mu         sync.Mutex
}

// NewTokenBucketRateLimiter создаёт новый rate limiter с заданной частотой.
// rate — количество запросов в секунду.
func NewTokenBucketRateLimiter(rate float64) *TokenBucketRateLimiter {
	return &TokenBucketRateLimiter{
		rate:       rate,
		tokens:     rate,
		maxTokens:  rate,
		lastUpdate: time.Now(),
	}
}

// Wait блокирует выполнение до получения токена.
func (r *TokenBucketRateLimiter) Wait(ctx context.Context) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	now := time.Now()
	elapsed := now.Sub(r.lastUpdate).Seconds()
	r.lastUpdate = now

	// Добавляем токены за прошедшее время
	r.tokens += elapsed * r.rate
	if r.tokens > r.maxTokens {
		r.tokens = r.maxTokens
	}

	// Если нет токенов, ждём
	if r.tokens < 1 {
		waitDuration := time.Duration((1 - r.tokens) / r.rate * float64(time.Second))
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(waitDuration):
		}
		r.tokens = 0
	} else {
		r.tokens--
	}

	return nil
}
