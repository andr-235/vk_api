package vk

import (
	"context"
	"net/http"
	"sync"
	"time"
)

type Option func(*Client)

type TokenSource int

const (
	TokenInParams TokenSource = iota
	TokenInHeader
)

func WithToken(token string) Option {
	return func(c *Client) { c.token = token }
}

func WithVersion(version string) Option {
	return func(c *Client) { c.version = version }
}

func WithLang(lang string) Option {
	return func(c *Client) { c.lang = lang }
}

func WithTestMode(enabled bool) Option {
	return func(c *Client) { c.testMode = enabled }
}

func WithBaseURL(baseURL string) Option {
	return func(c *Client) { c.baseURL = baseURL }
}

func WithHTTPClient(hc Doer) Option {
	return func(c *Client) { c.httpClient = hc }
}

func WithTokenSource(src TokenSource) Option {
	return func(c *Client) { c.tokenSource = src }
}

// WithRateLimiter устанавливает ограничитель частоты запросов.
func WithRateLimiter(limiter RateLimiter) Option {
	return func(c *Client) { c.rateLimiter = limiter }
}

// DefaultHTTPClient возвращает HTTP-клиент с оптимальными настройками
// для работы с VK API.
func DefaultHTTPClient() *http.Client {
	return &http.Client{
		Timeout: 30 * time.Second,
		Transport: &http.Transport{
			MaxIdleConns:        100,
			MaxIdleConnsPerHost: 10,
			IdleConnTimeout:     90 * time.Second,
		},
	}
}

// TokenBucketRateLimiter реализует алгоритм token bucket для rate limiting.
type TokenBucketRateLimiter struct {
	rate       float64   // токенов в секунду
	tokens     float64   // текущее количество токенов
	maxTokens  float64   // максимальное количество токенов
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
