package middleware

import (
	"context"
)

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
