package middleware

import (
	"context"
	"time"
)

// RequestInterceptor перехватывает запросы для логгирования, метрик, трассировки.
type RequestInterceptor interface {
	// InterceptRequest вызывается перед отправкой запроса.
	InterceptRequest(ctx context.Context, req *RequestContext) context.Context
	// InterceptResponse вызывается после получения ответа.
	InterceptResponse(ctx context.Context, req *RequestContext, resp *ResponseContext) context.Context
	// InterceptError вызывается при ошибке.
	InterceptError(ctx context.Context, req *RequestContext, err error) context.Context
}

// RequestContext содержит информацию о запросе.
type RequestContext struct {
	Method    string
	RequestID string
	StartTime time.Time
	Duration  time.Duration
	Attempt   int
}

// ResponseContext содержит информацию об ответе.
type ResponseContext struct {
	StatusCode int
	BodySize   int64
	Duration   time.Duration
}

// InterceptorChain позволяет объединять несколько interceptor'ов.
type InterceptorChain []RequestInterceptor

// InterceptRequest проходит по всем interceptor'ам.
func (c InterceptorChain) InterceptRequest(ctx context.Context, req *RequestContext) context.Context {
	for _, interceptor := range c {
		ctx = interceptor.InterceptRequest(ctx, req)
	}
	return ctx
}

// InterceptResponse проходит по всем interceptor'ам.
func (c InterceptorChain) InterceptResponse(ctx context.Context, req *RequestContext, resp *ResponseContext) context.Context {
	for _, interceptor := range c {
		ctx = interceptor.InterceptResponse(ctx, req, resp)
	}
	return ctx
}

// InterceptError проходит по всем interceptor'ам.
func (c InterceptorChain) InterceptError(ctx context.Context, req *RequestContext, err error) context.Context {
	for _, interceptor := range c {
		ctx = interceptor.InterceptError(ctx, req, err)
	}
	return ctx
}

// InterceptorFunc — функциональная обёртка для RequestInterceptor.
type InterceptorFunc struct {
	RequestFunc  func(ctx context.Context, req *RequestContext) context.Context
	ResponseFunc func(ctx context.Context, req *RequestContext, resp *ResponseContext) context.Context
	ErrorFunc    func(ctx context.Context, req *RequestContext, err error) context.Context
}

// NewInterceptorFunc создаёт interceptor из функций.
func NewInterceptorFunc(
	requestFunc func(ctx context.Context, req *RequestContext) context.Context,
	responseFunc func(ctx context.Context, req *RequestContext, resp *ResponseContext) context.Context,
	errorFunc func(ctx context.Context, req *RequestContext, err error) context.Context,
) *InterceptorFunc {
	return &InterceptorFunc{
		RequestFunc:  requestFunc,
		ResponseFunc: responseFunc,
		ErrorFunc:    errorFunc,
	}
}

// InterceptRequest вызывает RequestFunc.
func (f *InterceptorFunc) InterceptRequest(ctx context.Context, req *RequestContext) context.Context {
	if f.RequestFunc != nil {
		return f.RequestFunc(ctx, req)
	}
	return ctx
}

// InterceptResponse вызывает ResponseFunc.
func (f *InterceptorFunc) InterceptResponse(ctx context.Context, req *RequestContext, resp *ResponseContext) context.Context {
	if f.ResponseFunc != nil {
		return f.ResponseFunc(ctx, req, resp)
	}
	return ctx
}

// InterceptError вызывает ErrorFunc.
func (f *InterceptorFunc) InterceptError(ctx context.Context, req *RequestContext, err error) context.Context {
	if f.ErrorFunc != nil {
		return f.ErrorFunc(ctx, req, err)
	}
	return ctx
}
