package vk

import "fmt"

const (
	ErrorCodeAuth      = 5
	ErrorCodeCaptcha   = 14
	ErrorCodeRateLimit = 29
)

type RequestParam struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type VKError struct {
	Code          int
	Message       string
	RequestParams []RequestParam
	CaptchaSID    string
	CaptchaImg    string
	RedirectURI   string
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

func (e *VKError) IsAuth() bool {
	return e != nil && e.Code == ErrorCodeAuth
}

func (e *VKError) IsRateLimit() bool {
	return e != nil && e.Code == ErrorCodeRateLimit
}

func (e *VKError) IsCaptcha() bool {
	return e != nil && e.Code == ErrorCodeCaptcha
}

func newVKError(src *apiError) *VKError {
	if src == nil {
		return nil
	}

	return &VKError{
		Code:          src.Code,
		Message:       src.Message,
		RequestParams: src.RequestParams,
		CaptchaSID:    src.CaptchaSID,
		CaptchaImg:    src.CaptchaImg,
		RedirectURI:   src.RedirectURI,
	}
}
