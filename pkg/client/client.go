package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"sync"
	"sync/atomic"
	"time"

	"github.com/andr-235/vk_api/internal/encode"
	"github.com/andr-235/vk_api/pkg/config"
	"github.com/andr-235/vk_api/pkg/middleware"
	"github.com/andr-235/vk_api/pkg/ratelimit"
	"github.com/andr-235/vk_api/pkg/retry"
	"github.com/andr-235/vk_api/pkg/transport"
)

// requestIDCounter — глобальный счётчик для уникальных request ID
var requestIDCounter atomic.Uint64
var requestIDOnce sync.Once

// initRequestIDCounter инициализирует счётчик случайным значением
func initRequestIDCounter() {
	requestIDCounter.Store(uint64(time.Now().UnixNano()))
}

// Caller определяет интерфейс для вызова методов VK API.
type Caller interface {
	Call(ctx context.Context, method string, params, out any) error
	CallWithRawHandler(ctx context.Context, method string, params any, handler func(json.RawMessage) error) error
}

// ListResponse — универсальный ответ с коллекцией объектов.
type ListResponse[T any] struct {
	Count int `json:"count"`
	Items []T `json:"items"`
}

// Doer определяет интерфейс для выполнения HTTP-запросов.
// Экспортируется из pkg/transport для удобства.
type Doer = transport.Doer

// Client — основной клиент VK API.
type Client struct {
	config       config.Config
	transport    *transport.Transport
	interceptors middleware.InterceptorChain
	retryer      retry.Retryer
	rateLimiter  ratelimit.RateLimiter
}

// Option конфигурирует Client.
type Option func(*options)

type options struct {
	config       config.Config
	interceptors middleware.InterceptorChain
	retryer      retry.Retryer
	rateLimiter  ratelimit.RateLimiter
	httpClient   Doer
}

// New создаёт новый Client с заданными опциями.
func New(cfg config.Config, opts ...Option) *Client {
	opt := &options{
		config:       cfg,
		interceptors: middleware.InterceptorChain{},
		retryer:      retry.NewNoRetryer(),
		rateLimiter:  ratelimit.NewNoOpRateLimiter(),
		httpClient:   nil,
	}

	for _, o := range opts {
		o(opt)
	}

	// Валидация конфигурации
	if err := opt.config.Validate(); err != nil {
		panic(fmt.Sprintf("vk/client: invalid config: %v", err))
	}
	opt.config.Normalize()

	// Создаём транспорт
	httpClient := opt.httpClient
	if httpClient == nil {
		httpClient = transport.DefaultHTTPClient()
	}

	trans := transport.New(opt.config, httpClient)

	return &Client{
		config:       opt.config,
		transport:    trans,
		interceptors: opt.interceptors,
		retryer:      opt.retryer,
		rateLimiter:  opt.rateLimiter,
	}
}

// Call вызывает метод VK API и декодирует ответ в out.
func (c *Client) Call(ctx context.Context, method string, params, out any) error {
	return c.callWithRetry(ctx, method, params, func(ctx context.Context, values url.Values) error {
		return c.transport.Call(ctx, method, values, out)
	})
}

// CallWithRawHandler вызывает метод VK API с обработчиком сырого ответа.
func (c *Client) CallWithRawHandler(ctx context.Context, method string, params any, handler func(json.RawMessage) error) error {
	return c.callWithRetry(ctx, method, params, func(ctx context.Context, values url.Values) error {
		return c.transport.CallRaw(ctx, method, values, handler)
	})
}

// callWithRetry выполняет запрос с повторными попытками.
func (c *Client) callWithRetry(ctx context.Context, method string, params any, fn func(context.Context, url.Values) error) error {
	// Кодируем параметры один раз
	values, err := encode.Values(params)
	if err != nil {
		return fmt.Errorf("client: encode params: %w", err)
	}

	// Создаём контекст запроса
	reqCtx := &middleware.RequestContext{
		Method:    method,
		RequestID: generateRequestID(),
		StartTime: generateTime(),
		Attempt:   1,
	}

	// Вызываем interceptor'ы перед запросом
	ctx = c.interceptors.InterceptRequest(ctx, reqCtx)

	// Выполняем запрос с retry
	execFn := func() error {
		// Ждём разрешения от rate limiter
		if c.rateLimiter != nil {
			if err := c.rateLimiter.Wait(ctx); err != nil {
				return fmt.Errorf("client: rate limiter wait: %w", err)
			}
		}

		reqCtx.StartTime = generateTime()
		err := fn(ctx, values)
		reqCtx.Duration = timeSince(reqCtx.StartTime)

		return err
	}

	err = c.retryer.Execute(ctx, execFn)

	if err != nil {
		c.interceptors.InterceptError(ctx, reqCtx, err)
		return err
	}

	// Вызываем interceptor'ы после успешного ответа
	c.interceptors.InterceptResponse(ctx, reqCtx, &middleware.ResponseContext{
		StatusCode: 200,
		Duration:   reqCtx.Duration,
	})

	return nil
}

// Config возвращает конфигурацию клиента.
func (c *Client) Config() config.Config {
	return c.config
}

// generateRequestID генерирует уникальный ID запроса.
func generateRequestID() string {
	requestIDOnce.Do(initRequestIDCounter)
	return fmt.Sprintf("%d", requestIDCounter.Add(1))
}

func generateTime() time.Time {
	return time.Now()
}

func timeSince(t time.Time) time.Duration {
	return time.Since(t)
}
