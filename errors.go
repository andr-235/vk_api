package vk

import (
	"errors"
	"fmt"
	"net/http"
	"time"
)

const (
	ErrorCodeUnknown           = 1
	ErrorCodeAppDisabled       = 2
	ErrorCodeUnknownMethod     = 3
	ErrorCodeInvalidSignature  = 4
	ErrorCodeAuthFailed        = 5
	ErrorCodeTooManyRequests   = 6
	ErrorCodePermissionDenied  = 7
	ErrorCodeInvalidRequest    = 8
	ErrorCodeFloodControl      = 9
	ErrorCodeInternalServer    = 10
	ErrorCodeTestMode          = 11
	ErrorCodeCaptcha           = 14
	ErrorCodeAccessDenied      = 15
	ErrorCodeHTTPSRequired     = 16
	ErrorCodeValidationNeeded  = 17
	ErrorCodeUserDeleted       = 18
	ErrorCodeStandaloneOnly    = 20
	ErrorCodeMethodDisabled    = 23
	ErrorCodeConfirmation      = 24
	ErrorCodeGroupTokenInvalid = 27
	ErrorCodeAppTokenInvalid   = 28
	ErrorCodeRateLimit         = 29
	ErrorCodePrivateProfile    = 30

	ErrorCodeParamRequired = 100
	ErrorCodeInvalidAppID  = 101
	ErrorCodeInvalidUserID = 113
	ErrorCodeInvalidTime   = 150

	ErrorCodeAlbumAccessDenied = 200
	ErrorCodeAudioAccessDenied = 201
	ErrorCodeGroupAccessDenied = 203

	ErrorCodeAlbumFull = 300

	ErrorCodePaymentRequired = 500

	ErrorCodeAdsPermissionDenied = 600
	ErrorCodeAdsError            = 603
)

// VKError представляет ошибку VK API.
type VKError struct {
	Code          int
	Message       string
	RequestParams []RequestParam

	CaptchaSID  string
	CaptchaImg  string
	RedirectURI string

	ConfirmationText string
}

func (e *VKError) Error() string {
	if e == nil {
		return "<nil>"
	}
	if e.Message == "" {
		return fmt.Sprintf("vk api error %d", e.Code)
	}
	return fmt.Sprintf("vk api error %d: %s", e.Code, e.Message)
}

func (e *VKError) Unwrap() error {
	return nil
}

func (e *VKError) IsAuth() bool {
	return e != nil && e.Code == ErrorCodeAuthFailed
}

func (e *VKError) IsCaptcha() bool {
	return e != nil && e.Code == ErrorCodeCaptcha
}

func (e *VKError) IsRateLimit() bool {
	return e != nil && (e.Code == ErrorCodeTooManyRequests || e.Code == ErrorCodeRateLimit)
}

func (e *VKError) IsValidation() bool {
	return e != nil && e.Code == ErrorCodeValidationNeeded
}

func (e *VKError) IsPermission() bool {
	return e != nil && (e.Code == ErrorCodePermissionDenied || e.Code == ErrorCodeAccessDenied)
}

func (e *VKError) IsParam() bool {
	return e != nil && (e.Code == ErrorCodeInvalidRequest || e.Code == ErrorCodeParamRequired)
}

func (e *VKError) IsPrivate() bool {
	return e != nil && e.Code == ErrorCodePrivateProfile
}

func AsVKError(err error) (*VKError, bool) {
	var vkErr *VKError
	if errors.As(err, &vkErr) {
		return vkErr, true
	}
	return nil, false
}

// AuthError представляет ошибку аутентификации.
type AuthError struct {
	Code    int
	Message string
}

func (e *AuthError) Error() string {
	return fmt.Sprintf("auth error %d: %s", e.Code, e.Message)
}

func (e *AuthError) Unwrap() error {
	return &VKError{Code: e.Code, Message: e.Message}
}

// RateLimitError представляет ошибку превышения лимита запросов.
type RateLimitError struct {
	Code       int
	Message    string
	RetryAfter time.Duration
}

func (e *RateLimitError) Error() string {
	if e.RetryAfter > 0 {
		return fmt.Sprintf("rate limit exceeded, retry after %v", e.RetryAfter)
	}
	return fmt.Sprintf("rate limit exceeded: %s", e.Message)
}

func (e *RateLimitError) Unwrap() error {
	return &VKError{Code: e.Code, Message: e.Message}
}

// CaptchaError представляет ошибку CAPTCHA.
type CaptchaError struct {
	Code       int
	Message    string
	CaptchaSID string
	CaptchaImg string
}

func (e *CaptchaError) Error() string {
	return fmt.Sprintf("captcha required: %s", e.Message)
}

func (e *CaptchaError) Unwrap() error {
	return &VKError{Code: e.Code, Message: e.Message, CaptchaSID: e.CaptchaSID, CaptchaImg: e.CaptchaImg}
}

// PermissionError представляет ошибку доступа.
type PermissionError struct {
	Code    int
	Message string
}

func (e *PermissionError) Error() string {
	return fmt.Sprintf("permission denied: %s", e.Message)
}

func (e *PermissionError) Unwrap() error {
	return &VKError{Code: e.Code, Message: e.Message}
}

// ValidationError представляет ошибку валидации параметров.
type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	if e.Field != "" {
		return fmt.Sprintf("validation error for field %q: %s", e.Field, e.Message)
	}
	return fmt.Sprintf("validation error: %s", e.Message)
}

// HTTPError представляет HTTP-ошибку.
type HTTPError struct {
	StatusCode int
	Body       string
}

func (e *HTTPError) Error() string {
	return fmt.Sprintf("http error %d: %s", e.StatusCode, e.Body)
}

// IsAuth возвращает true, если ошибка связана с аутентификацией.
func IsAuth(err error) bool {
	if e, ok := err.(*AuthError); ok {
		return e != nil
	}
	if e, ok := err.(*VKError); ok {
		return e.IsAuth()
	}
	return false
}

// IsRateLimit возвращает true, если ошибка связана с rate limiting.
func IsRateLimit(err error) bool {
	_, ok := err.(*RateLimitError)
	if !ok {
		if e, ok := err.(*VKError); ok {
			return e.IsRateLimit()
		}
	}
	return ok
}

// IsCaptcha возвращает true, если ошибка требует прохождения CAPTCHA.
func IsCaptcha(err error) bool {
	_, ok := err.(*CaptchaError)
	if !ok {
		if e, ok := err.(*VKError); ok {
			return e.IsCaptcha()
		}
	}
	return ok
}

// IsPermission возвращает true, если ошибка связана с правами доступа.
func IsPermission(err error) bool {
	if _, ok := err.(*PermissionError); ok {
		return true
	}
	if e, ok := err.(*VKError); ok {
		return e.IsPermission()
	}
	return false
}

// ParseHTTPError создаёт HTTPError из http.Response.
func ParseHTTPError(resp *http.Response, body string) error {
	return &HTTPError{
		StatusCode: resp.StatusCode,
		Body:       body,
	}
}
