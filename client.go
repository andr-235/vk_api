package vk

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/andr-235/vk_api/internal/encode"
	"github.com/andr-235/vk_api/internal/transport"
)

const (
	defaultBaseURL = "https://api.vk.ru/method/"
	defaultVersion = "5.199"
)

// Client — основной клиент VK API.
// Потокобезопасен после создания.
type Client struct {
	config        Config
	transport     TransportConfig
	interceptors  InterceptorChain
	retryer       *SimpleRetryer
}

// Проверка, что Client реализует интерфейс Caller
var _ Caller = (*Client)(nil)

// New создаёт новый Client с заданными опциями.
func New(opts ...Option) *Client {
	cfg := DefaultConfig()
	for _, opt := range opts {
		opt(&cfg)
	}

	// Валидация и нормализация
	if err := cfg.Validate(); err != nil {
		// Для обратной совместимости паникуем при ошибке валидации
		// Используйте NewBuilder().Build() для обработки ошибок
		panic(fmt.Sprintf("vk: invalid config: %v", err))
	}
	cfg.normalize()

	return newClient(cfg)
}

// newClient создаёт Client из валидной конфигурации.
func newClient(cfg Config) *Client {
	return &Client{
		config:       cfg,
		transport:    NewTransportConfig(cfg),
		interceptors: InterceptorChain{},
		retryer:      NewSimpleRetryer(0, DefaultRetryPolicy()),
	}
}

// WithInterceptors устанавливает interceptor'ы для клиента.
func (c *Client) WithInterceptors(interceptors ...RequestInterceptor) *Client {
	c.interceptors = interceptors
	return c
}

// WithRetryer устанавливает retryer для клиента.
func (c *Client) WithRetryer(retryer *SimpleRetryer) *Client {
	c.retryer = retryer
	return c
}

// Call вызывает метод VK API и декодирует ответ в out.
func (c *Client) Call(ctx context.Context, method string, params, out any) error {
	return c.callWithRetry(ctx, method, params, func(ctx context.Context) error {
		respEnv, err := c.doRequest(ctx, method, params)
		if err != nil {
			return err
		}

		if out == nil || len(respEnv.Response) == 0 {
			return nil
		}

		if err := json.Unmarshal(respEnv.Response, out); err != nil {
			return fmt.Errorf("vk: decode response payload: %w", err)
		}

		return nil
	})
}

// CallWithRawHandler вызывает метод VK API с обработчиком сырого ответа.
func (c *Client) CallWithRawHandler(ctx context.Context, method string, params any, handler func(json.RawMessage) error) error {
	return c.callWithRetry(ctx, method, params, func(ctx context.Context) error {
		respEnv, err := c.doRequest(ctx, method, params)
		if err != nil {
			return err
		}

		if handler == nil || len(respEnv.Response) == 0 {
			return nil
		}

		return handler(respEnv.Response)
	})
}

// callWithRetry выполняет запрос с повторными попытками.
func (c *Client) callWithRetry(ctx context.Context, method string, params any, fn func(context.Context) error) error {
	if c.retryer == nil {
		return fn(ctx)
	}
	return c.retryer.Execute(ctx, func() error {
		return fn(ctx)
	})
}

// doRequest выполняет один запрос к VK API.
func (c *Client) doRequest(ctx context.Context, method string, params any) (*responseEnvelope, error) {
	// Ждём разрешения от rate limiter
	if c.transport.RateLimiter != nil {
		if err := c.transport.RateLimiter.Wait(ctx); err != nil {
			return nil, fmt.Errorf("vk: rate limiter wait: %w", err)
		}
	}

	// Кодируем параметры
	values, err := encode.Values(params)
	if err != nil {
		return nil, fmt.Errorf("vk: encode params: %w", err)
	}

	// Создаём контекст запроса
	reqCtx := &RequestContext{
		Method:    method,
		RequestID: generateRequestID(),
		StartTime: time.Now(),
		Attempt:   1,
	}

	// Вызываем interceptor'ы перед запросом
	ctx = c.interceptors.InterceptRequest(ctx, reqCtx)

	// Строим HTTP-запрос
	req, err := c.transport.Builder.Build(ctx, method, values)
	if err != nil {
		c.interceptors.InterceptError(ctx, reqCtx, err)
		return nil, err
	}

	// Выполняем запрос
	resp, err := c.transport.HTTPClient.Do(req)
	reqCtx.Duration = time.Since(reqCtx.StartTime)
	if err != nil {
		c.interceptors.InterceptError(ctx, reqCtx, err)
		return nil, err
	}
	defer resp.Body.Close()

	// Обрабатываем ответ
	respEnv, err := c.transport.Handler.Handle(resp)
	respCtx := &ResponseContext{
		StatusCode: resp.StatusCode,
		BodySize:   int(resp.ContentLength),
		Duration:   reqCtx.Duration,
	}

	if err != nil {
		c.interceptors.InterceptError(ctx, reqCtx, err)
		return nil, err
	}

	// Вызываем interceptor'ы после успешного ответа
	c.interceptors.InterceptResponse(ctx, reqCtx, respCtx)

	return respEnv, nil
}

// Config возвращает конфигурацию клиента.
func (c *Client) Config() Config {
	return c.config
}

// transportConfig возвращает конфигурацию транспорта (для обратной совместимости).
func (c *Client) transportConfig() transport.Config {
	tokenSource := transport.TokenInParams
	if c.config.TokenSource == TokenInHeader {
		tokenSource = transport.TokenInHeader
	}

	return transport.Config{
		BaseURL:     c.config.BaseURL,
		Version:     c.config.Version,
		Lang:        c.config.Lang,
		TestMode:    c.config.TestMode,
		Token:       c.config.Token,
		TokenSource: tokenSource,
		HTTPClient:  c.config.HTTPClient,
	}
}

// generateRequestID генерирует уникальный ID запроса.
func generateRequestID() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}
