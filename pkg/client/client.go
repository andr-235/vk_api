// Package client предоставляет основной клиент для взаимодействия с VK API.
//
// Клиент поддерживает:
//   - Конфигурацию через pkg/config
//   - HTTP транспорт через pkg/transport
//   - Retry логику через pkg/retry
//   - Middleware через pkg/middleware
//   - Rate limiting через pkg/ratelimit
//
// Пример использования:
//
//	client := client.New(config.DefaultConfig(),
//		client.WithToken("your-token"),
//		client.WithRateLimiter(ratelimit.NewTokenBucketRateLimiter(3.0)),
//	)
//
//	err := client.Call(ctx, "users.get", params, &users)
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
//
// Этот интерфейс позволяет мокировать клиент для тестирования.
type Caller interface {
	// Call вызывает метод VK API и декодирует ответ в out.
	Call(ctx context.Context, method string, params, out any) error
	// CallWithRawHandler вызывает метод VK API с обработчиком сырого ответа.
	CallWithRawHandler(ctx context.Context, method string, params any, handler func(json.RawMessage) error) error
}

// ListResponse — универсальный ответ с коллекцией объектов.
//
// Используется для методов, возвращающих списки (друзья, подписчики, и т.д.).
type ListResponse[T any] struct {
	Count int `json:"count"` // Общее количество объектов
	Items []T `json:"items"` // Массив объектов
}

// Doer определяет интерфейс для выполнения HTTP-запросов.
// Экспортируется из pkg/transport для удобства.
type Doer = transport.Doer

// Client — основной клиент VK API.
//
// Потокобезопасен после создания. Все методы могут вызываться из разных горутин.
//
// Пример создания:
//
//	cfg := config.DefaultConfig()
//	cfg.Token = "your-token"
//	client := client.New(cfg,
//		client.WithRateLimiter(ratelimit.NewTokenBucketRateLimiter(3.0)),
//		client.WithRetryer(retry.NewSimpleRetryer(3, retry.DefaultPolicy())),
//	)
type Client struct {
	config       config.Config
	transport    *transport.Transport
	interceptors middleware.InterceptorChain
	retryer      retry.Retryer
	rateLimiter  ratelimit.RateLimiter
}

// Option конфигурирует Client.
type Option func(*options)

// options содержит опции для создания Client.
type options struct {
	config       config.Config
	interceptors middleware.InterceptorChain
	retryer      retry.Retryer
	rateLimiter  ratelimit.RateLimiter
	httpClient   Doer
}

// New создаёт новый Client с заданными опциями.
//
// Если конфигурация невалидна, функция паникует.
// Для обработки ошибок используйте NewBuilder().Build().
//
// Пример:
//
//	client := client.New(config.DefaultConfig(),
//		client.WithToken("token"),
//		client.WithVersion("5.199"),
//	)
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
//
// method — имя метода VK API (например, "users.get").
// params — структура параметров с тегами `url:"name"`.
// out — указатель на структуру или слайс для декодирования ответа.
//
// Пример:
//
//	var users []User
//	err := client.Call(ctx, "users.get", params, &users)
func (c *Client) Call(ctx context.Context, method string, params, out any) error {
	return c.callWithRetry(ctx, method, params, func(ctx context.Context, values url.Values) error {
		return c.transport.Call(ctx, method, values, out)
	})
}

// CallWithRawHandler вызывает метод VK API с обработчиком сырого ответа.
//
// handler получает сырой JSON ответа и может обработать его самостоятельно.
//
// Пример:
//
//	err := client.CallWithRawHandler(ctx, "users.get", params, func(raw json.RawMessage) error {
//		// кастомная обработка
//	})
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
