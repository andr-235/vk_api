package client

import (
	"github.com/andr-235/vk_api/pkg/config"
	"github.com/andr-235/vk_api/pkg/middleware"
	"github.com/andr-235/vk_api/pkg/ratelimit"
	"github.com/andr-235/vk_api/pkg/retry"
	"github.com/andr-235/vk_api/pkg/transport"
)

// WithToken устанавливает токен доступа.
func WithToken(token string) Option {
	return func(o *options) {
		o.config.Token = token
	}
}

// WithVersion устанавливает версию API.
func WithVersion(version string) Option {
	return func(o *options) {
		o.config.Version = version
	}
}

// WithLang устанавливает язык ответов.
func WithLang(lang string) Option {
	return func(o *options) {
		o.config.Lang = lang
	}
}

// WithTestMode включает/выключает тестовый режим.
func WithTestMode(enabled bool) Option {
	return func(o *options) {
		o.config.TestMode = enabled
	}
}

// WithBaseURL устанавливает базовый URL API.
func WithBaseURL(baseURL string) Option {
	return func(o *options) {
		o.config.BaseURL = baseURL
	}
}

// WithTokenSource устанавливает способ передачи токена.
func WithTokenSource(src config.TokenSource) Option {
	return func(o *options) {
		o.config.TokenSource = src
	}
}

// WithInterceptors устанавливает interceptor'ы.
func WithInterceptors(interceptors ...middleware.RequestInterceptor) Option {
	return func(o *options) {
		o.interceptors = interceptors
	}
}

// WithRetryer устанавливает retryer.
func WithRetryer(retryer retry.Retryer) Option {
	return func(o *options) {
		o.retryer = retryer
	}
}

// WithRateLimiter устанавливает rate limiter.
func WithRateLimiter(limiter ratelimit.RateLimiter) Option {
	return func(o *options) {
		o.rateLimiter = limiter
	}
}

// WithHTTPClient устанавливает кастомный HTTP-клиент.
func WithHTTPClient(hc transport.Doer) Option {
	return func(o *options) {
		o.httpClient = hc
	}
}
