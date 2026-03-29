package config

// TokenSource определяет способ передачи токена в VK API.
type TokenSource int

const (
	// TokenInParams передаёт токен в параметрах запроса (access_token).
	TokenInParams TokenSource = iota

	// TokenInHeader передаёт токен в заголовке Authorization.
	TokenInHeader
)

// Option конфигурирует Config.
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

// WithTokenSource устанавливает способ передачи токена.
func WithTokenSource(src TokenSource) Option {
	return func(c *Config) { c.TokenSource = src }
}
