package vk

import "encoding/json"

type responseEnvelope struct {
	Response json.RawMessage `json:"response"`
	Error    *apiError       `json:"error"`
}

type apiError struct {
	Code          int            `json:"error_code"`
	Messages      string         `json:"error_msg"`
	RequestParams []RequestParam `json:"request_params"`
	CaptchaSID    string         `json:"captha_sid"`
	CaptchaImg    string         `json:"captha_img"`
	RedirectURI   string         `json:"redirect_url"`
}
