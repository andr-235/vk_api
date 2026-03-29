package vk

import (
	"context"
	"encoding/json"
	"net/http"
)

type RequestParam struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// Caller определяет интерфейс для вызова методов VK API.
// Позволяет мокировать клиент для тестирования.
type Caller interface {
	Call(ctx context.Context, method string, params, out any) error
	CallWithRawHandler(ctx context.Context, method string, params any, handler func(json.RawMessage) error) error
}

// Doer определяет интерфейс для выполнения HTTP-запросов.
// Позволяет использовать кастомные HTTP-клиенты с retry-логикой, circuit breaker и т.д.
type Doer interface {
	Do(req *http.Request) (*http.Response, error)
}

// RateLimiter определяет интерфейс для ограничения частоты запросов.
type RateLimiter interface {
	Wait(ctx context.Context) error
}

// ListResponse — универсальный ответ с коллекцией объектов.
type ListResponse[T any] struct {
	Count int `json:"count"`
	Items []T `json:"items"`
}

// PaginatedResponse — ответ с коллекцией объектов и расширенной информацией.
type PaginatedResponse[T any] struct {
	Count int `json:"count"`
	Items []T `json:"items"`
}

// Coordinates представляет географические координаты.
type Coordinates struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

// TokenSource определяет способ передачи токена в VK API.
type TokenSource int

const (
	// TokenInParams передаёт токен в параметрах запроса (access_token).
	TokenInParams TokenSource = iota
	// TokenInHeader передаёт токен в заголовке Authorization.
	TokenInHeader
)
