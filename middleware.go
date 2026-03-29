package vk

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
	Method      string
	Params      map[string]string
	Headers     map[string]string
	RequestID   string
	StartTime   time.Time
	Duration    time.Duration
	Attempt     int
}

// ResponseContext содержит информацию об ответе.
type ResponseContext struct {
	StatusCode int
	BodySize   int
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

// LoggingInterceptor логирует запросы и ответы.
type LoggingInterceptor struct {
	LogFunc func(ctx context.Context, msg string, fields map[string]any)
}

// NewLoggingInterceptor создаёт logging interceptor.
func NewLoggingInterceptor(logFunc func(ctx context.Context, msg string, fields map[string]any)) *LoggingInterceptor {
	return &LoggingInterceptor{LogFunc: logFunc}
}

// InterceptRequest логирует начало запроса.
func (l *LoggingInterceptor) InterceptRequest(ctx context.Context, req *RequestContext) context.Context {
	if l.LogFunc != nil {
		l.LogFunc(ctx, "VK request started", map[string]any{
			"request_id": req.RequestID,
			"method":     req.Method,
			"attempt":    req.Attempt,
		})
	}
	return ctx
}

// InterceptResponse логирует завершение запроса.
func (l *LoggingInterceptor) InterceptResponse(ctx context.Context, req *RequestContext, resp *ResponseContext) context.Context {
	if l.LogFunc != nil {
		l.LogFunc(ctx, "VK request completed", map[string]any{
			"request_id":  req.RequestID,
			"method":      req.Method,
			"status_code": resp.StatusCode,
			"body_size":   resp.BodySize,
			"duration_ms": resp.Duration.Milliseconds(),
		})
	}
	return ctx
}

// InterceptError логирует ошибку запроса.
func (l *LoggingInterceptor) InterceptError(ctx context.Context, req *RequestContext, err error) context.Context {
	if l.LogFunc != nil {
		l.LogFunc(ctx, "VK request failed", map[string]any{
			"request_id": req.RequestID,
			"method":     req.Method,
			"attempt":    req.Attempt,
			"error":      err.Error(),
		})
	}
	return ctx
}

// MetricsInterceptor собирает метрики запросов.
type MetricsInterceptor struct {
	RecordLatency  func(ctx context.Context, method string, statusCode int, duration time.Duration)
	RecordError    func(ctx context.Context, method string, err error)
	RecordRequests func(ctx context.Context, method string)
}

// NewMetricsInterceptor создаёт metrics interceptor.
func NewMetricsInterceptor(
	recordLatency func(ctx context.Context, method string, statusCode int, duration time.Duration),
	recordError func(ctx context.Context, method string, err error),
	recordRequests func(ctx context.Context, method string),
) *MetricsInterceptor {
	return &MetricsInterceptor{
		RecordLatency:  recordLatency,
		RecordError:    recordError,
		RecordRequests: recordRequests,
	}
}

// InterceptRequest регистрирует начало запроса.
func (m *MetricsInterceptor) InterceptRequest(ctx context.Context, req *RequestContext) context.Context {
	if m.RecordRequests != nil {
		m.RecordRequests(ctx, req.Method)
	}
	return ctx
}

// InterceptResponse регистрирует успешный запрос.
func (m *MetricsInterceptor) InterceptResponse(ctx context.Context, req *RequestContext, resp *ResponseContext) context.Context {
	if m.RecordLatency != nil {
		m.RecordLatency(ctx, req.Method, resp.StatusCode, resp.Duration)
	}
	return ctx
}

// InterceptError регистрирует ошибку.
func (m *MetricsInterceptor) InterceptError(ctx context.Context, req *RequestContext, err error) context.Context {
	if m.RecordError != nil {
		m.RecordError(ctx, req.Method, err)
	}
	return ctx
}
