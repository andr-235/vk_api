package vk

import (
	"errors"
	"fmt"

	"github.com/andr-235/vk_api/internal/transport"
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

func newVKError(src *transport.APIError) *VKError {
	if src == nil {
		return nil
	}

	params := make([]RequestParam, 0, len(src.RequestParams))
	for _, p := range src.RequestParams {
		params = append(params, RequestParam{
			Key:   p.Key,
			Value: p.Value,
		})
	}

	return &VKError{
		Code:             src.Code,
		Message:          src.Message,
		RequestParams:    params,
		CaptchaSID:       src.CaptchaSID,
		CaptchaImg:       src.CaptchaImg,
		RedirectURI:      src.RedirectURI,
		ConfirmationText: src.ConfirmationText,
	}
}
