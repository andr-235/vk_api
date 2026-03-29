package ratelimit

import (
	"context"
)

// RateLimiter определяет интерфейс для ограничения частоты запросов.
type RateLimiter interface {
	// Wait блокирует выполнение до получения разрешения.
	// Возвращает ctx.Err() если контекст отменён.
	Wait(ctx context.Context) error
}
