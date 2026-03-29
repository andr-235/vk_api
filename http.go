package vk

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

// RequestBuilder строит HTTP-запросы для VK API.
type RequestBuilder struct {
	baseURL     string
	version     string
	lang        string
	testMode    bool
	token       string
	tokenSource TokenSource
}

// NewRequestBuilder создаёт новый RequestBuilder.
func NewRequestBuilder(cfg Config) *RequestBuilder {
	return &RequestBuilder{
		baseURL:     cfg.BaseURL,
		version:     cfg.Version,
		lang:        cfg.Lang,
		testMode:    cfg.TestMode,
		token:       cfg.Token,
		tokenSource: cfg.TokenSource,
	}
}

// Build строит HTTP-запрос.
func (b *RequestBuilder) Build(ctx context.Context, method string, params url.Values) (*http.Request, error) {
	values := make(url.Values, len(params)+4)

	// Копируем параметры
	for k, v := range params {
		values[k] = v
	}

	// Добавляем обязательные параметры VK API
	values.Set("v", b.version)
	if b.lang != "" {
		values.Set("lang", b.lang)
	}
	if b.testMode {
		values.Set("test_mode", "1")
	}
	if b.token != "" && b.tokenSource == TokenInParams {
		values.Set("access_token", b.token)
	}

	// Строим URL
	reqURL := b.baseURL + method + "?" + values.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, reqURL, nil)
	if err != nil {
		return nil, fmt.Errorf("vk: build request: %w", err)
	}

	req.Header.Set("Accept", "application/json")

	if b.token != "" && b.tokenSource == TokenInHeader {
		req.Header.Set("Authorization", "Bearer "+b.token)
	}

	return req, nil
}

// ResponseHandler обрабатывает HTTP-ответы от VK API.
type ResponseHandler struct {
	errorMapper ErrorMapper
}

// ErrorMapper маппит VK API ошибки на типизированные ошибки Go.
type ErrorMapper interface {
	MapError(err *vkErrorEnvelope) error
}

// NewResponseHandler создаёт новый ResponseHandler.
func NewResponseHandler(errorMapper ErrorMapper) *ResponseHandler {
	return &ResponseHandler{errorMapper: errorMapper}
}

// Handle обрабатывает HTTP-ответ.
func (h *ResponseHandler) Handle(resp *http.Response) (*responseEnvelope, error) {
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("vk: read response: %w", err)
	}

	// Проверяем HTTP статус
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, ParseHTTPError(resp, string(body))
	}

	// Парсим ответ VK API
	var respEnv responseEnvelope
	if err := json.Unmarshal(body, &respEnv); err != nil {
		return nil, fmt.Errorf("vk: decode envelope: %w", err)
	}

	// Обрабатываем ошибку VK API
	if respEnv.Error != nil {
		return nil, h.errorMapper.MapError(respEnv.Error)
	}

	return &respEnv, nil
}

// DefaultErrorMapper реализует маппинг ошибок по умолчанию.
type DefaultErrorMapper struct{}

// MapError маппит ошибку VK API на типизированную ошибку Go.
func (m *DefaultErrorMapper) MapError(err *vkErrorEnvelope) error {
	vkErr := &VKError{
		Code:             err.Code,
		Message:          err.Message,
		CaptchaSID:       err.CaptchaSID,
		CaptchaImg:       err.CaptchaImg,
		RedirectURI:      err.RedirectURI,
		ConfirmationText: err.ConfirmationText,
	}

	// Возвращаем специфичную ошибку в зависимости от кода
	switch err.Code {
	case ErrorCodeAuthFailed:
		return &AuthError{Code: err.Code, Message: err.Message}
	case ErrorCodeCaptcha:
		return &CaptchaError{
			Code:       err.Code,
			Message:    err.Message,
			CaptchaSID: err.CaptchaSID,
			CaptchaImg: err.CaptchaImg,
		}
	case ErrorCodeTooManyRequests, ErrorCodeRateLimit:
		return &RateLimitError{Code: err.Code, Message: err.Message}
	case ErrorCodePermissionDenied, ErrorCodeAccessDenied:
		return &PermissionError{Code: err.Code, Message: err.Message}
	default:
		return vkErr
	}
}

// TransportConfig содержит конфигурацию транспорта.
type TransportConfig struct {
	Config      Config
	Builder     *RequestBuilder
	Handler     *ResponseHandler
	HTTPClient  Doer
	RateLimiter RateLimiter
}

// NewTransportConfig создаёт конфигурацию транспорта.
func NewTransportConfig(cfg Config) TransportConfig {
	return TransportConfig{
		Config:      cfg,
		Builder:     NewRequestBuilder(cfg),
		Handler:     NewResponseHandler(&DefaultErrorMapper{}),
		HTTPClient:  cfg.HTTPClient,
		RateLimiter: cfg.RateLimiter,
	}
}
