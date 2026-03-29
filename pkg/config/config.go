// Package config предоставляет конфигурацию для VK API клиента.
//
// Пример использования:
//
//	cfg := config.DefaultConfig()
//	cfg.Token = "your-token"
//
//	// или через builder
//	cfg, err := config.NewBuilder().
//		WithToken("token").
//		WithVersion("5.199").
//		Build()
package config

import (
	"fmt"
	"net/url"
)

const (
	defaultBaseURL = "https://api.vk.ru/method/"
	defaultVersion = "5.199"
)

// Config содержит конфигурацию VK клиента.
type Config struct {
	// Token — токен доступа.
	Token string

	// Version — версия VK API (по умолчанию "5.199").
	Version string

	// Lang — язык ответов (по умолчанию пустой).
	Lang string

	// TestMode — включает тестовый режим API.
	TestMode bool

	// BaseURL — базовый URL API (по умолчанию "https://api.vk.ru/method/").
	BaseURL string

	// TokenSource — способ передачи токена.
	TokenSource TokenSource
}

// DefaultConfig возвращает конфигурацию по умолчанию.
func DefaultConfig() Config {
	return Config{
		Version:     defaultVersion,
		BaseURL:     defaultBaseURL,
		TokenSource: TokenInParams,
	}
}

// Validate проверяет валидность конфигурации.
//
// Возвращает ConfigError если конфигурация невалидна.
func (c *Config) Validate() error {
	if c.Version == "" {
		return &ConfigError{Field: "version", Message: "version is required"}
	}
	if c.BaseURL != "" {
		if _, err := url.Parse(c.BaseURL); err != nil {
			return &ConfigError{Field: "baseURL", Message: fmt.Sprintf("invalid URL: %v", err)}
		}
	}
	return nil
}

// normalize нормализует конфигурацию после валидации.
func (c *Config) normalize() {
	if c.BaseURL != "" && c.BaseURL[len(c.BaseURL)-1] != '/' {
		c.BaseURL += "/"
	}
}

// Normalize вызывает normalize и возвращает Config для fluent API.
func (c *Config) Normalize() *Config {
	c.normalize()
	return c
}

// ConfigError представляет ошибку конфигурации.
type ConfigError struct {
	Field   string
	Message string
}

// Error реализует интерфейс error.
func (e *ConfigError) Error() string {
	return fmt.Sprintf("config: invalid %s: %s", e.Field, e.Message)
}

// TokenSource определяет способ передачи токена в VK API.
type TokenSource int

const (
	// TokenInParams передаёт токен в параметрах запроса (access_token).
	TokenInParams TokenSource = iota

	// TokenInHeader передаёт токен в заголовке Authorization.
	TokenInHeader
)
