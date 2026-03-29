package vk

import (
	"fmt"
	"net/url"
	"time"
)

// Config содержит конфигурацию VK клиента.
type Config struct {
	Token       string
	Version     string
	Lang        string
	TestMode    bool
	BaseURL     string
	HTTPClient  Doer
	TokenSource TokenSource
	RateLimiter RateLimiter
}

// DefaultConfig возвращает конфигурацию по умолчанию.
func DefaultConfig() Config {
	return Config{
		Version:     defaultVersion,
		BaseURL:     defaultBaseURL,
		HTTPClient:  DefaultHTTPClient(),
		TokenSource: TokenInParams,
	}
}

// Validate проверяет валидность конфигурации.
func (c *Config) Validate() error {
	if c.Version == "" {
		return &ValidationError{Field: "version", Message: "version is required"}
	}
	if c.BaseURL != "" {
		if _, err := url.Parse(c.BaseURL); err != nil {
			return &ValidationError{Field: "baseURL", Message: fmt.Sprintf("invalid URL: %v", err)}
		}
	}
	if c.HTTPClient == nil {
		return &ValidationError{Field: "httpClient", Message: "httpClient is required"}
	}
	return nil
}

// normalize нормализует конфигурацию после валидации.
func (c *Config) normalize() {
	if c.BaseURL != "" && c.BaseURL[len(c.BaseURL)-1] != '/' {
		c.BaseURL += "/"
	}
}

// ClientBuilder предоставляет fluent API для создания клиента.
type ClientBuilder struct {
	config Config
}

// NewBuilder создаёт новый builder с конфигурацией по умолчанию.
func NewBuilder() *ClientBuilder {
	return &ClientBuilder{
		config: DefaultConfig(),
	}
}

// WithToken устанавливает токен доступа.
func (b *ClientBuilder) WithToken(token string) *ClientBuilder {
	b.config.Token = token
	return b
}

// WithVersion устанавливает версию API.
func (b *ClientBuilder) WithVersion(version string) *ClientBuilder {
	b.config.Version = version
	return b
}

// WithLang устанавливает язык ответов.
func (b *ClientBuilder) WithLang(lang string) *ClientBuilder {
	b.config.Lang = lang
	return b
}

// WithTestMode включает/выключает тестовый режим.
func (b *ClientBuilder) WithTestMode(enabled bool) *ClientBuilder {
	b.config.TestMode = enabled
	return b
}

// WithBaseURL устанавливает базовый URL API.
func (b *ClientBuilder) WithBaseURL(baseURL string) *ClientBuilder {
	b.config.BaseURL = baseURL
	return b
}

// WithHTTPClient устанавливает кастомный HTTP-клиент.
func (b *ClientBuilder) WithHTTPClient(hc Doer) *ClientBuilder {
	b.config.HTTPClient = hc
	return b
}

// WithTokenSource устанавливает способ передачи токена.
func (b *ClientBuilder) WithTokenSource(src TokenSource) *ClientBuilder {
	b.config.TokenSource = src
	return b
}

// WithRateLimiter устанавливает ограничитель частоты запросов.
func (b *ClientBuilder) WithRateLimiter(limiter RateLimiter) *ClientBuilder {
	b.config.RateLimiter = limiter
	return b
}

// Build создаёт новый Client с заданной конфигурацией.
func (b *ClientBuilder) Build() (*Client, error) {
	cfg := b.config
	if err := cfg.Validate(); err != nil {
		return nil, err
	}
	cfg.normalize()

	return newClient(cfg), nil
}

// MustBuild создаёт Client или паникует при ошибке.
func (b *ClientBuilder) MustBuild() *Client {
	client, err := b.Build()
	if err != nil {
		panic(err)
	}
	return client
}

// RequestConfig содержит параметры одного запроса.
type RequestConfig struct {
	Method      string
	Params      map[string]string
	Headers     map[string]string
	Timeout     time.Duration
	MaxRetries  int
	RetryPolicy RetryPolicy
}

// DefaultRequestConfig возвращает конфигурацию запроса по умолчанию.
func DefaultRequestConfig() RequestConfig {
	return RequestConfig{
		Headers:     make(map[string]string),
		Params:      make(map[string]string),
		MaxRetries:  0,
		RetryPolicy: DefaultRetryPolicy(),
	}
}
