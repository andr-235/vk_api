package testkit

import (
	"context"
	"encoding/json"
	"net/http"
)

// MockCaller — мок интерфейса vk.Caller для тестирования.
type MockCaller struct {
	CallFunc               func(ctx context.Context, method string, params, out any) error
	CallWithRawHandlerFunc func(ctx context.Context, method string, params any, handler func(json.RawMessage) error) error
}

// Call реализует интерфейс vk.Caller.
func (m *MockCaller) Call(ctx context.Context, method string, params, out any) error {
	if m.CallFunc != nil {
		return m.CallFunc(ctx, method, params, out)
	}
	return nil
}

// CallWithRawHandler реализует интерфейс vk.Caller.
func (m *MockCaller) CallWithRawHandler(ctx context.Context, method string, params any, handler func(json.RawMessage) error) error {
	if m.CallWithRawHandlerFunc != nil {
		return m.CallWithRawHandlerFunc(ctx, method, params, handler)
	}
	return nil
}

// MockDoer — мок интерфейса vk.Doer для тестирования HTTP-клиента.
type MockDoer struct {
	DoFunc func(req *http.Request) (*http.Response, error)
}

// Do реализует интерфейс vk.Doer.
func (m *MockDoer) Do(req *http.Request) (*http.Response, error) {
	if m.DoFunc != nil {
		return m.DoFunc(req)
	}
	return nil, nil
}

// MockRateLimiter — мок интерфейса vk.RateLimiter для тестирования.
type MockRateLimiter struct {
	WaitFunc func(ctx context.Context) error
}

// Wait реализует интерфейс vk.RateLimiter.
func (m *MockRateLimiter) Wait(ctx context.Context) error {
	if m.WaitFunc != nil {
		return m.WaitFunc(ctx)
	}
	return nil
}

// NewMockCaller создаёт новый MockCaller с заданными функциями.
func NewMockCaller(opts ...MockCallerOption) *MockCaller {
	mc := &MockCaller{}
	for _, opt := range opts {
		opt(mc)
	}
	return mc
}

// MockCallerOption — опция для настройки MockCaller.
type MockCallerOption func(*MockCaller)

// WithCallFunc устанавливает функцию для вызова Call.
func WithCallFunc(fn func(ctx context.Context, method string, params, out any) error) MockCallerOption {
	return func(mc *MockCaller) {
		mc.CallFunc = fn
	}
}

// WithCallWithRawHandlerFunc устанавливает функцию для вызова CallWithRawHandler.
func WithCallWithRawHandlerFunc(fn func(ctx context.Context, method string, params any, handler func(json.RawMessage) error) error) MockCallerOption {
	return func(mc *MockCaller) {
		mc.CallWithRawHandlerFunc = fn
	}
}

// NewMockDoer создаёт новый MockDoer с заданной функцией.
func NewMockDoer(doFunc func(req *http.Request) (*http.Response, error)) *MockDoer {
	return &MockDoer{
		DoFunc: doFunc,
	}
}

// NewMockRateLimiter создаёт новый MockRateLimiter с заданной функцией.
func NewMockRateLimiter(waitFunc func(ctx context.Context) error) *MockRateLimiter {
	return &MockRateLimiter{
		WaitFunc: waitFunc,
	}
}

// Responder помогает создавать ответы для моков.
type Responder struct {
	data any
	err  error
}

// NewResponder создаёт новый Responder.
func NewResponder(data any, err error) *Responder {
	return &Responder{data: data, err: err}
}

// ToCallFunc возвращает функцию для использования с MockCaller.CallFunc.
func (r *Responder) ToCallFunc() func(ctx context.Context, method string, params, out any) error {
	return func(ctx context.Context, method string, params, out any) error {
		if r.err != nil {
			return r.err
		}
		if r.data != nil && out != nil {
			// Копируем данные в out
			dataBytes, err := json.Marshal(r.data)
			if err != nil {
				return err
			}
			return json.Unmarshal(dataBytes, out)
		}
		return nil
	}
}
