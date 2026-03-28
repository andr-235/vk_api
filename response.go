package vk

import "encoding/json"

type responseEnvelope struct {
	Response json.RawMessage `json:"response"`
	Error    *apiError       `json:"error"`
}

type apiError struct {
	Code          int            `json:"error_code"`
	Message       string         `json:"error_msg"`
	RequestParams []RequestParam `json:"request_params,omitempty"`

	CaptchaSID  string `json:"captcha_sid,omitempty"`
	CaptchaImg  string `json:"captcha_img,omitempty"`
	RedirectURI string `json:"redirect_uri,omitempty"`

	ConfirmationText string `json:"confirmation_text,omitempty"`
}
