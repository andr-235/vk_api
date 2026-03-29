package client

import (
	"github.com/andr-235/vk_api/pkg/config"
	"github.com/andr-235/vk_api/pkg/middleware"
	"github.com/andr-235/vk_api/pkg/ratelimit"
	"github.com/andr-235/vk_api/pkg/retry"
)

// WithToken устанавливает токен доступа.
func WithToken(token string) Option {
	return func(o *options) {
		o.config.Token = token
	}
}

// WithVersion устанавливает версию API (по умолчанию "5.199").
func WithVersion(version string) Option {
	return func(o *options) {
		o.config.Version = version
	}
}

// WithLang устанавливает язык ответов (например, "ru", "en").
func WithLang(lang string) Option {
	return func(o *options) {
		o.config.Lang = lang
	}
}

// WithTestMode включает/выключает тестовый режим API.
func WithTestMode(enabled bool) Option {
	return func(o *options) {
		o.config.TestMode = enabled
	}
}

// WithBaseURL устанавливает базовый URL API.
// По умолчанию: "https://api.vk.ru/method/"
func WithBaseURL(baseURL string) Option {
	return func(o *options) {
		o.config.BaseURL = baseURL
	}
}

// WithTokenSource устанавливает способ передачи токена.
// TokenInParams (по умолчанию) — токен в параметрах запроса.
// TokenInHeader — токен в заголовке Authorization.
func WithTokenSource(src config.TokenSource) Option {
	return func(o *options) {
		o.config.TokenSource = src
	}
}

// WithInterceptors устанавливает interceptor'ы для middleware.
func WithInterceptors(interceptors ...middleware.RequestInterceptor) Option {
	return func(o *options) {
		o.interceptors = interceptors
	}
}

// WithRetryer устанавливает retryer для повторных попыток.
// По умолчанию используется NoRetryer (без повторных попыток).
func WithRetryer(retryer retry.Retryer) Option {
	return func(o *options) {
		o.retryer = retryer
	}
}

// WithRateLimiter устанавливает rate limiter для ограничения частоты запросов.
// По умолчанию используется NoOpRateLimiter (без ограничений).
func WithRateLimiter(limiter ratelimit.RateLimiter) Option {
	return func(o *options) {
		o.rateLimiter = limiter
	}
}

// WithHTTPClient устанавливает кастомный HTTP-клиент.
// По умолчанию используется DefaultHTTPClient().
func WithHTTPClient(hc Doer) Option {
	return func(o *options) {
		o.httpClient = hc
	}
}
