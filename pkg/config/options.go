// Package config предоставляет конфигурацию для VK API клиента.
package config

// Option конфигурирует Config.
type Option func(*Config)

// WithToken устанавливает токен доступа.
func WithToken(token string) Option {
	return func(c *Config) { c.Token = token }
}

// WithVersion устанавливает версию API (по умолчанию "5.199").
func WithVersion(version string) Option {
	return func(c *Config) { c.Version = version }
}

// WithLang устанавливает язык ответов (например, "ru", "en").
func WithLang(lang string) Option {
	return func(c *Config) { c.Lang = lang }
}

// WithTestMode включает/выключает тестовый режим API.
func WithTestMode(enabled bool) Option {
	return func(c *Config) { c.TestMode = enabled }
}

// WithBaseURL устанавливает базовый URL API.
// По умолчанию: "https://api.vk.ru/method/"
func WithBaseURL(baseURL string) Option {
	return func(c *Config) { c.BaseURL = baseURL }
}

// WithTokenSource устанавливает способ передачи токена.
// TokenInParams (по умолчанию) — токен в параметрах запроса.
// TokenInHeader — токен в заголовке Authorization.
func WithTokenSource(src TokenSource) Option {
	return func(c *Config) { c.TokenSource = src }
}
