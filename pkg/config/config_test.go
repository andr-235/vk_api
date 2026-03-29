package config_test

import (
	"testing"

	"github.com/andr-235/vk_api/pkg/config"
)

func TestDefaultConfig(t *testing.T) {
	cfg := config.DefaultConfig()

	if cfg.Version != "5.199" {
		t.Errorf("Version = %q, want %q", cfg.Version, "5.199")
	}
	if cfg.BaseURL != "https://api.vk.ru/method/" {
		t.Errorf("BaseURL = %q, want %q", cfg.BaseURL, "https://api.vk.ru/method/")
	}
	if cfg.TokenSource != config.TokenInParams {
		t.Errorf("TokenSource = %v, want %v", cfg.TokenSource, config.TokenInParams)
	}
}

func TestConfigValidate(t *testing.T) {
	t.Run("valid config", func(t *testing.T) {
		cfg := config.DefaultConfig()
		err := cfg.Validate()
		if err != nil {
			t.Errorf("Validate() error = %v", err)
		}
	})

	t.Run("empty version", func(t *testing.T) {
		cfg := config.Config{}
		err := cfg.Validate()
		if err == nil {
			t.Error("Validate() should return error for empty version")
		}

		if err == nil {
			return
		}
		// Проверяем, что это ConfigError
		_ = err.(*config.ConfigError)
	})

	t.Run("invalid URL", func(t *testing.T) {
		cfg := config.Config{
			Version: "5.199",
			BaseURL: "://invalid-url",
		}
		err := cfg.Validate()
		if err == nil {
			t.Error("Validate() should return error for invalid URL")
		}
	})
}

func TestConfigNormalize(t *testing.T) {
	t.Run("adds trailing slash", func(t *testing.T) {
		cfg := config.Config{BaseURL: "https://api.vk.ru/method"}
		cfg.Normalize()
		if cfg.BaseURL != "https://api.vk.ru/method/" {
			t.Errorf("BaseURL = %q, want %q", cfg.BaseURL, "https://api.vk.ru/method/")
		}
	})

	t.Run("keeps trailing slash", func(t *testing.T) {
		cfg := config.Config{BaseURL: "https://api.vk.ru/method/"}
		cfg.Normalize()
		if cfg.BaseURL != "https://api.vk.ru/method/" {
			t.Errorf("BaseURL = %q, want %q", cfg.BaseURL, "https://api.vk.ru/method/")
		}
	})
}

func TestConfigError(t *testing.T) {
	err := &config.ConfigError{
		Field:   "version",
		Message: "version is required",
	}

	expected := "config: invalid version: version is required"
	if err.Error() != expected {
		t.Errorf("Error() = %q, want %q", err.Error(), expected)
	}
}

func TestTokenSource(t *testing.T) {
	if config.TokenInParams != 0 {
		t.Errorf("TokenInParams = %d, want 0", config.TokenInParams)
	}
	if config.TokenInHeader != 1 {
		t.Errorf("TokenInHeader = %d, want 1", config.TokenInHeader)
	}
}

func TestBuilder(t *testing.T) {
	t.Run("build with defaults", func(t *testing.T) {
		cfg, err := config.NewBuilder().Build()
		if err != nil {
			t.Fatalf("Build() error = %v", err)
		}
		if cfg.Version != "5.199" {
			t.Errorf("Version = %q, want %q", cfg.Version, "5.199")
		}
	})

	t.Run("build with token", func(t *testing.T) {
		cfg, err := config.NewBuilder().
			WithToken("test-token").
			WithVersion("5.200").
			WithLang("en").
			WithTestMode(true).
			WithBaseURL("https://custom.api.vk.com/").
			WithTokenSource(config.TokenInHeader).
			Build()
		if err != nil {
			t.Fatalf("Build() error = %v", err)
		}
		if cfg.Token != "test-token" {
			t.Errorf("Token = %q, want %q", cfg.Token, "test-token")
		}
		if cfg.Version != "5.200" {
			t.Errorf("Version = %q, want %q", cfg.Version, "5.200")
		}
		if cfg.Lang != "en" {
			t.Errorf("Lang = %q, want %q", cfg.Lang, "en")
		}
		if !cfg.TestMode {
			t.Error("TestMode = false, want true")
		}
		if cfg.BaseURL != "https://custom.api.vk.com/" {
			t.Errorf("BaseURL = %q, want %q", cfg.BaseURL, "https://custom.api.vk.com/")
		}
		if cfg.TokenSource != config.TokenInHeader {
			t.Errorf("TokenSource = %v, want %v", cfg.TokenSource, config.TokenInHeader)
		}
	})

	t.Run("build with invalid config", func(t *testing.T) {
		_, err := config.NewBuilder().
			WithVersion("").
			Build()
		if err == nil {
			t.Error("Build() should return error for invalid config")
		}
	})

	t.Run("must build", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("MustBuild() should panic on invalid config")
			}
		}()

		_ = config.NewBuilder().WithVersion("").MustBuild()
	})

	t.Run("must build success", func(t *testing.T) {
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("MustBuild() panicked: %v", r)
			}
		}()

		cfg := config.NewBuilder().MustBuild()
		if cfg.Version != "5.199" {
			t.Errorf("Version = %q, want %q", cfg.Version, "5.199")
		}
	})
}

func TestOptions(t *testing.T) {
	t.Run("WithToken", func(t *testing.T) {
		cfg := config.DefaultConfig()
		config.WithToken("test")(&cfg)
		if cfg.Token != "test" {
			t.Errorf("Token = %q, want %q", cfg.Token, "test")
		}
	})

	t.Run("WithVersion", func(t *testing.T) {
		cfg := config.DefaultConfig()
		config.WithVersion("5.200")(&cfg)
		if cfg.Version != "5.200" {
			t.Errorf("Version = %q, want %q", cfg.Version, "5.200")
		}
	})

	t.Run("WithLang", func(t *testing.T) {
		cfg := config.DefaultConfig()
		config.WithLang("en")(&cfg)
		if cfg.Lang != "en" {
			t.Errorf("Lang = %q, want %q", cfg.Lang, "en")
		}
	})

	t.Run("WithTestMode", func(t *testing.T) {
		cfg := config.DefaultConfig()
		config.WithTestMode(true)(&cfg)
		if !cfg.TestMode {
			t.Error("TestMode = false, want true")
		}
	})

	t.Run("WithBaseURL", func(t *testing.T) {
		cfg := config.DefaultConfig()
		config.WithBaseURL("https://custom.api.vk.com/")(&cfg)
		if cfg.BaseURL != "https://custom.api.vk.com/" {
			t.Errorf("BaseURL = %q, want %q", cfg.BaseURL, "https://custom.api.vk.com/")
		}
	})

	t.Run("WithTokenSource", func(t *testing.T) {
		cfg := config.DefaultConfig()
		config.WithTokenSource(config.TokenInHeader)(&cfg)
		if cfg.TokenSource != config.TokenInHeader {
			t.Errorf("TokenSource = %v, want %v", cfg.TokenSource, config.TokenInHeader)
		}
	})
}
