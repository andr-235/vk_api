package middleware

import (
	"context"
	"time"
)

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
