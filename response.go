package vk

import "encoding/json"

type responseEnvelope struct {
	Response json.RawMessage `json:"response"`
	Error    *vkErrorEnvelope `json:"error"`
}

// vkErrorEnvelope представляет ошибку VK API.
type vkErrorEnvelope struct {
	Code          int            `json:"error_code"`
	Message       string         `json:"error_msg"`
	RequestParams []RequestParam `json:"request_params,omitempty"`

	CaptchaSID  string `json:"captcha_sid,omitempty"`
	CaptchaImg  string `json:"captcha_img,omitempty"`
	RedirectURI string `json:"redirect_uri,omitempty"`

	ConfirmationText string `json:"confirmation_text,omitempty"`
}

// apiError — устаревшее имя, используйте vkErrorEnvelope.
// Deprecated: используйте vkErrorEnvelope.
type apiError = vkErrorEnvelope
