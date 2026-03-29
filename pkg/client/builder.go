// Package client предоставляет основной клиент для взаимодействия с VK API.
package client

import (
	"github.com/andr-235/vk_api/pkg/config"
	"github.com/andr-235/vk_api/pkg/middleware"
	"github.com/andr-235/vk_api/pkg/ratelimit"
	"github.com/andr-235/vk_api/pkg/retry"
	"github.com/andr-235/vk_api/pkg/transport"
)

// Builder предоставляет fluent API для создания клиента.
//
// Пример использования:
//
//	client, err := client.NewBuilder().
//		WithToken("token").
//		WithVersion("5.199").
//		WithRateLimiter(ratelimit.NewTokenBucketRateLimiter(3.0)).
//		Build()
type Builder struct {
	config       config.Config
	interceptors []middleware.RequestInterceptor
	retryer      retry.Retryer
	rateLimiter  ratelimit.RateLimiter
	httpClient   Doer
}

// NewBuilder создаёт новый builder с конфигурацией по умолчанию.
func NewBuilder() *Builder {
	return &Builder{
		config:      config.DefaultConfig(),
		interceptors: middleware.InterceptorChain{},
		retryer:     retry.NewNoRetryer(),
		rateLimiter: ratelimit.NewNoOpRateLimiter(),
	}
}

// WithToken устанавливает токен доступа.
func (b *Builder) WithToken(token string) *Builder {
	b.config.Token = token
	return b
}

// WithVersion устанавливает версию API.
func (b *Builder) WithVersion(version string) *Builder {
	b.config.Version = version
	return b
}

// WithLang устанавливает язык ответов.
func (b *Builder) WithLang(lang string) *Builder {
	b.config.Lang = lang
	return b
}

// WithTestMode включает/выключает тестовый режим.
func (b *Builder) WithTestMode(enabled bool) *Builder {
	b.config.TestMode = enabled
	return b
}

// WithBaseURL устанавливает базовый URL API.
func (b *Builder) WithBaseURL(baseURL string) *Builder {
	b.config.BaseURL = baseURL
	return b
}

// WithTokenSource устанавливает способ передачи токена.
func (b *Builder) WithTokenSource(src config.TokenSource) *Builder {
	b.config.TokenSource = src
	return b
}

// WithInterceptors устанавливает interceptor'ы.
func (b *Builder) WithInterceptors(interceptors ...middleware.RequestInterceptor) *Builder {
	b.interceptors = interceptors
	return b
}

// WithRetryer устанавливает retryer.
func (b *Builder) WithRetryer(retryer retry.Retryer) *Builder {
	b.retryer = retryer
	return b
}

// WithRateLimiter устанавливает rate limiter.
func (b *Builder) WithRateLimiter(limiter ratelimit.RateLimiter) *Builder {
	b.rateLimiter = limiter
	return b
}

// WithHTTPClient устанавливает кастомный HTTP-клиент.
func (b *Builder) WithHTTPClient(hc Doer) *Builder {
	b.httpClient = hc
	return b
}

// Build создаёт новый Client.
//
// Возвращает ошибку, если конфигурация невалидна.
func (b *Builder) Build() (*Client, error) {
	if err := b.config.Validate(); err != nil {
		return nil, err
	}
	b.config.Normalize()

	opt := &options{
		interceptors: b.interceptors,
		retryer:      b.retryer,
		rateLimiter:  b.rateLimiter,
		httpClient:   b.httpClient,
	}

	// Создаём транспорт
	httpClient := opt.httpClient
	if httpClient == nil {
		httpClient = transport.DefaultHTTPClient()
	}

	trans := transport.New(b.config, httpClient)

	return &Client{
		config:       b.config,
		transport:    trans,
		interceptors: opt.interceptors,
		retryer:      opt.retryer,
		rateLimiter:  opt.rateLimiter,
	}, nil
}

// MustBuild создаёт Client или паникует при ошибке.
//
// Используйте, когда уверены в валидности конфигурации.
func (b *Builder) MustBuild() *Client {
	client, err := b.Build()
	if err != nil {
		panic(err)
	}
	return client
}
