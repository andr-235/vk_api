package transport

import "fmt"

// VKError представляет ошибку VK API.
type VKError struct {
	Code             int
	Message          string
	RequestParams    []RequestParam
	CaptchaSID       string
	CaptchaImg       string
	RedirectURI      string
	ConfirmationText string
}

// Error реализует интерфейс error.
func (e *VKError) Error() string {
	if e.Message == "" {
		return fmt.Sprintf("vk api error %d", e.Code)
	}
	return fmt.Sprintf("vk api error %d: %s", e.Code, e.Message)
}

// AuthError представляет ошибку аутентификации.
type AuthError struct {
	Code    int
	Message string
}

// Error реализует интерфейс error.
func (e *AuthError) Error() string {
	return fmt.Sprintf("auth error %d: %s", e.Code, e.Message)
}

// CaptchaError представляет ошибку CAPTCHA.
type CaptchaError struct {
	Code       int
	Message    string
	CaptchaSID string
	CaptchaImg string
}

// Error реализует интерфейс error.
func (e *CaptchaError) Error() string {
	return fmt.Sprintf("captcha required: %s", e.Message)
}

// RateLimitError представляет ошибку превышения лимита запросов.
type RateLimitError struct {
	Code       int
	Message    string
	RetryAfter int
}

// Error реализует интерфейс error.
func (e *RateLimitError) Error() string {
	if e.RetryAfter > 0 {
		return fmt.Sprintf("rate limit exceeded, retry after %d seconds", e.RetryAfter)
	}
	return fmt.Sprintf("rate limit exceeded: %s", e.Message)
}

// PermissionError представляет ошибку доступа.
type PermissionError struct {
	Code    int
	Message string
}

// Error реализует интерфейс error.
func (e *PermissionError) Error() string {
	return fmt.Sprintf("permission denied: %s", e.Message)
}

// MapError маппит ошибку VK API на типизированную ошибку Go.
func MapError(err *vkErrorEnvelope) error {
	vkErr := &VKError{
		Code:             err.Code,
		Message:          err.Message,
		RequestParams:    err.RequestParams,
		CaptchaSID:       err.CaptchaSID,
		CaptchaImg:       err.CaptchaImg,
		RedirectURI:      err.RedirectURI,
		ConfirmationText: err.ConfirmationText,
	}

	// Возвращаем специфичную ошибку в зависимости от кода
	switch err.Code {
	case 5: // ErrorCodeAuthFailed
		return &AuthError{Code: err.Code, Message: err.Message}
	case 14: // ErrorCodeCaptcha
		return &CaptchaError{
			Code:       err.Code,
			Message:    err.Message,
			CaptchaSID: err.CaptchaSID,
			CaptchaImg: err.CaptchaImg,
		}
	case 6, 29: // ErrorCodeTooManyRequests, ErrorCodeRateLimit
		return &RateLimitError{Code: err.Code, Message: err.Message}
	case 7, 15: // ErrorCodePermissionDenied, ErrorCodeAccessDenied
		return &PermissionError{Code: err.Code, Message: err.Message}
	default:
		return vkErr
	}
}
