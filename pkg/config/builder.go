package config

// Builder предоставляет fluent API для создания конфигурации.
type Builder struct {
	config Config
}

// NewBuilder создаёт новый builder с конфигурацией по умолчанию.
func NewBuilder() *Builder {
	return &Builder{
		config: DefaultConfig(),
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
func (b *Builder) WithTokenSource(src TokenSource) *Builder {
	b.config.TokenSource = src
	return b
}

// Build создаёт Config с заданной конфигурацией.
func (b *Builder) Build() (Config, error) {
	cfg := b.config
	if err := cfg.Validate(); err != nil {
		return Config{}, err
	}
	cfg.normalize()
	return cfg, nil
}

// MustBuild создаёт Config или паникует при ошибке.
func (b *Builder) MustBuild() Config {
	cfg, err := b.Build()
	if err != nil {
		panic(err)
	}
	return cfg
}
