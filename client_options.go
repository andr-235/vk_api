package vk

// Option конфигурирует Client.
type Option func(*Config)

// WithToken устанавливает токен доступа.
func WithToken(token string) Option {
	return func(c *Config) { c.Token = token }
}

// WithVersion устанавливает версию API.
func WithVersion(version string) Option {
	return func(c *Config) { c.Version = version }
}

// WithLang устанавливает язык ответов.
func WithLang(lang string) Option {
	return func(c *Config) { c.Lang = lang }
}

// WithTestMode включает/выключает тестовый режим.
func WithTestMode(enabled bool) Option {
	return func(c *Config) { c.TestMode = enabled }
}

// WithBaseURL устанавливает базовый URL API.
func WithBaseURL(baseURL string) Option {
	return func(c *Config) { c.BaseURL = baseURL }
}

// WithHTTPClient устанавливает кастомный HTTP-клиент.
func WithHTTPClient(hc Doer) Option {
	return func(c *Config) { c.HTTPClient = hc }
}

// WithTokenSource устанавливает способ передачи токена.
func WithTokenSource(src TokenSource) Option {
	return func(c *Config) { c.TokenSource = src }
}

// WithRateLimiter устанавливает ограничитель частоты запросов.
func WithRateLimiter(limiter RateLimiter) Option {
	return func(c *Config) { c.RateLimiter = limiter }
}
