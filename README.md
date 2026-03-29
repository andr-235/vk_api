# VK API Client for Go

[![Go Reference](https://pkg.go.dev/badge/github.com/andr-235/vk_api.svg)](https://pkg.go.dev/github.com/andr-235/vk_api)
[![Go Report Card](https://goreportcard.com/badge/github.com/andr-235/vk_api)](https://goreportcard.com/report/github.com/andr-235/vk_api)

Модульный клиент для VK API с поддержкой retry, middleware, rate limiting и типизированными ответами.

## Архитектура

Проект использует **слоистую архитектуру** с разделением ответственности:

```
pkg/
├── client/        # Основной клиент (Client, Caller интерфейс)
├── config/        # Конфигурация (Config, Builder, Options)
├── transport/     # HTTP транспорт (Transport, ErrorMapper)
├── retry/         # Retry-логика (Retryer, RetryPolicy, ExponentialBackoff)
├── middleware/    # Middleware (Interceptor, Logging, Metrics)
└── ratelimit/     # Rate limiting (RateLimiter, TokenBucket)

api/               # Высокоуровневые обёртки (users, groups, messages, wall)
internal/          # Внутренние утилиты (encode, testkit)
```

### Преимущества архитектуры

- **SRP** — каждый пакет отвечает за одну задачу
- **Тестируемость** — интерфейсы позволяют мокировать зависимости
- **Переиспользование** — пакеты независимы и могут использоваться отдельно
- **Расширяемость** — легко добавить middleware, кастомный транспорт, retry-логику

## Установка

```bash
go get github.com/andr-235/vk_api
```

## Быстрый старт

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/andr-235/vk_api/pkg/client"
    "github.com/andr-235/vk_api/pkg/config"
    "github.com/andr-235/vk_api/pkg/ratelimit"
    "github.com/andr-235/vk_api/api/users"
)

func main() {
    // Создание клиента с rate limiting
    c, err := client.NewBuilder().
        WithToken("your_access_token").
        WithVersion("5.199").
        WithRateLimiter(ratelimit.NewTokenBucketRateLimiter(3.0)). // 3 запроса/сек
        Build()
    if err != nil {
        log.Fatal(err)
    }

    // Вызов API
    var user []users.User
    err = c.Call(context.Background(), "users.get", users.GetParams{
        UserIDs: []string{"1"},
        Fields:  []string{"bdate"},
    }, &user)
    
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("User: %+v\n", user[0])
}
```

## Конфигурация

### Через Builder (рекомендуется)

```go
c, err := client.NewBuilder().
    WithToken("token").
    WithVersion("5.199").
    WithLang("ru").
    WithTestMode(false).
    WithRateLimiter(ratelimit.NewTokenBucketRateLimiter(3.0)).
    WithRetryer(retry.NewSimpleRetryer(3, retry.DefaultPolicy())).
    Build()
```

### Через New

```go
c := client.New(
    config.Config{
        Token:   "token",
        Version: "5.199",
    },
    client.WithRateLimiter(ratelimit.NewTokenBucketRateLimiter(3.0)),
    client.WithRetryer(retry.NewSimpleRetryer(3, retry.DefaultPolicy())),
)
```

## Retry-логика

```go
import "github.com/andr-235/vk_api/pkg/retry"

// Retry с экспоненциальной задержкой
retryer := retry.NewSimpleRetryer(
    3, // макс. 3 попытки
    &retry.ExponentialBackoff{
        Initial:    100 * time.Millisecond,
        Max:        10 * time.Second,
        Multiplier: 2.0,
        Jitter:     0.1, // 10% разброс
    },
)

c := client.New(cfg, client.WithRetryer(retryer))
```

## Middleware

### Логгирование

```go
import "github.com/andr-235/vk_api/pkg/middleware"

logger := middleware.NewLoggingInterceptor(
    func(ctx context.Context, msg string, fields map[string]any) {
        log.Printf("%s: %v", msg, fields)
    },
)

c := client.New(cfg, client.WithInterceptors(logger))
```

### Метрики

```go
metrics := middleware.NewMetricsInterceptor(
    func(ctx context.Context, method string, statusCode int, duration time.Duration) {
        // записать latency
    },
    func(ctx context.Context, method string, err error) {
        // записать error
    },
    func(ctx context.Context, method string) {
        // записать request count
    },
)

c := client.New(cfg, client.WithInterceptors(metrics))
```

## Rate Limiting

```go
import "github.com/andr-235/vk_api/pkg/ratelimit"

// 3 запроса в секунду (стандартный лимит VK API)
limiter := ratelimit.NewTokenBucketRateLimiter(3.0)

c := client.New(cfg, client.WithRateLimiter(limiter))

// Без ограничений (для тестов)
limiter := ratelimit.NewNoOpRateLimiter()
```

## Обработка ошибок

```go
import "github.com/andr-235/vk_api/pkg/transport"

err := c.Call(ctx, "users.get", params, &users)

switch e := err.(type) {
case *transport.AuthError:
    // Ошибка аутентификации — обновить токен
    log.Printf("Auth error: %v", e)

case *transport.CaptchaError:
    // Требуется CAPTCHA
    log.Printf("Captcha: %s", e.CaptchaImg)

case *transport.RateLimitError:
    // Превышен лимит
    log.Printf("Rate limit, retry after %d seconds", e.RetryAfter)

case *transport.PermissionError:
    // Нет прав доступа
    log.Printf("Permission denied: %v", e)

case *transport.VKError:
    // Другая ошибка VK API
    log.Printf("VK error %d: %s", e.Code, e.Message)

case *transport.HTTPError:
    // HTTP-ошибка
    log.Printf("HTTP error %d: %s", e.StatusCode, e.Body)
}
```

## Высокоуровневые обёртки

Пакет `api/` предоставляет типизированные обёртки для методов VK API:

```go
import (
    "github.com/andr-235/vk_api/api/users"
    "github.com/andr-235/vk_api/api/groups"
)

// Получить пользователей
usersList, err := users.Get(ctx, c, users.GetParams{
    UserIDs: []string{"1", "2"},
    Fields:  []string{"bdate", "city"},
})

// Получить группу
groupsList, err := groups.GetByID(ctx, c, groups.GetByIDParams{
    GroupIDs: []string{"vk"},
    Fields:   []string{"description"},
})
```

### Доступные модули

- `api/users` — пользователи (get, search, followers, subscriptions)
- `api/groups` — сообщества (getByID, getMembers, getBanned, адреса, callback-серверы)
- `api/messages` — сообщения (send)
- `api/wall` — записи на стене (get)

## Тестирование

### Моки для тестирования

```go
import "github.com/andr-235/vk_api/internal/testkit"

mock := testkit.NewMockCaller(
    testkit.WithCallFunc(func(ctx context.Context, method string, params, out any) error {
        // Вернуть тестовые данные
        return nil
    }),
)

// Использовать в тестах
users, err := users.Get(ctx, mock, users.GetParams{...})
```

### Integration тесты

```bash
# Запуск integration тестов (требуется VK_TOKEN)
go test -tags=integration ./...
```

## Производительность

### Benchmark

```bash
go test -bench=. -benchmem ./internal/encode/...
```

Результаты (AMD Ryzen 5 5600X):

```
BenchmarkValues-12                 1209015    1010 ns/op
BenchmarkValuesLarge-12             698469    1752 ns/op
BenchmarkEncodeMap-12              2602002     465 ns/op
BenchmarkCachePerformance-12       1358098     866 ns/op
```

## Лицензия

MIT License — см. [LICENSE](LICENSE)
