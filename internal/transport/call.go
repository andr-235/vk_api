package transport

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	internalencode "github.com/andr-235/vk_api/internal/encode"
)

type Doer interface {
	Do(req *http.Request) (*http.Response, error)
}

type TokenSource int

const (
	TokenInParams TokenSource = iota
	TokenInHeader
)

type Config struct {
	BaseURL     string
	Version     string
	Lang        string
	TestMode    bool
	Token       string
	TokenSource TokenSource
	HTTPClient  Doer
}

// EncodeValues кодирует параметры в url.Values.
func EncodeValues(v any) (url.Values, error) {
	return internalencode.Values(v)
}

type ResponseEnvelope struct {
	Response json.RawMessage `json:"response"`
	Error    *_APIError      `json:"error"`
}

type _APIError struct {
	Code          int    `json:"error_code"`
	Message       string `json:"error_msg"`
	RequestParams []struct {
		Key   string `json:"key"`
		Value string `json:"value"`
	} `json:"request_params,omitempty"`

	CaptchaSID  string `json:"captcha_sid,omitempty"`
	CaptchaImg  string `json:"captcha_img,omitempty"`
	RedirectURI string `json:"redirect_uri,omitempty"`

	ConfirmationText string `json:"confirmation_text,omitempty"`
}

// TransportError — ошибка транспорта (не экспортируется наружу).
type TransportError struct {
	Code         int
	Message      string
	CaptchaSID   string
	CaptchaImg   string
	RedirectURI  string
	Confirmation string
}

func (e *TransportError) Error() string {
	if e == nil {
		return "<nil>"
	}
	if e.Message == "" {
		return fmt.Sprintf("vk api error %d", e.Code)
	}
	return fmt.Sprintf("vk api error %d: %s", e.Code, e.Message)
}
