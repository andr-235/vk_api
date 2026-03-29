package transport

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/andr-235/vk_api/pkg/config"
)

// DefaultHTTPClient возвращает HTTP-клиент с оптимальными настройками.
func DefaultHTTPClient() *http.Client {
	return &http.Client{
		Timeout: 30 * time.Second,
		Transport: &http.Transport{
			MaxIdleConns:        100,
			MaxIdleConnsPerHost: 10,
			IdleConnTimeout:     90 * time.Second,
		},
	}
}

// Doer определяет интерфейс для выполнения HTTP-запросов.
type Doer interface {
	Do(req *http.Request) (*http.Response, error)
}

// Transport выполняет HTTP-запросы к VK API.
type Transport struct {
	baseURL     string
	version     string
	lang        string
	testMode    bool
	token       string
	tokenSource config.TokenSource
	httpClient  Doer
}

// New создаёт новый Transport из конфигурации.
func New(cfg config.Config, httpClient Doer) *Transport {
	return &Transport{
		baseURL:     cfg.BaseURL,
		version:     cfg.Version,
		lang:        cfg.Lang,
		testMode:    cfg.TestMode,
		token:       cfg.Token,
		tokenSource: cfg.TokenSource,
		httpClient:  httpClient,
	}
}

// Call вызывает метод VK API и декодирует ответ в out.
func (t *Transport) Call(ctx context.Context, method string, params url.Values, out any) error {
	respEnv, err := t.doRequest(ctx, method, params)
	if err != nil {
		return err
	}

	if out == nil || len(respEnv.Response) == 0 {
		return nil
	}

	if err := json.Unmarshal(respEnv.Response, out); err != nil {
		return fmt.Errorf("transport: decode response payload: %w", err)
	}

	return nil
}

// CallRaw вызывает метод VK API с обработчиком сырого ответа.
func (t *Transport) CallRaw(ctx context.Context, method string, params url.Values, handler func(json.RawMessage) error) error {
	respEnv, err := t.doRequest(ctx, method, params)
	if err != nil {
		return err
	}

	if handler == nil || len(respEnv.Response) == 0 {
		return nil
	}

	return handler(respEnv.Response)
}

// doRequest выполняет один запрос к VK API.
func (t *Transport) doRequest(ctx context.Context, method string, params url.Values) (*responseEnvelope, error) {
	// Строим URL
	reqURL := t.baseURL + method + "?" + t.encodeParams(params)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, reqURL, nil)
	if err != nil {
		return nil, fmt.Errorf("transport: build request: %w", err)
	}

	req.Header.Set("Accept", "application/json")

	if t.token != "" && t.tokenSource == config.TokenInHeader {
		req.Header.Set("Authorization", "Bearer "+t.token)
	}

	// Выполняем запрос
	resp, err := t.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("transport: do request: %w", err)
	}
	defer resp.Body.Close()

	// Читаем тело
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("transport: read response: %w", err)
	}

	// Проверяем HTTP статус
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, &HTTPError{
			StatusCode: resp.StatusCode,
			Body:       string(body),
		}
	}

	// Парсим ответ VK API
	var respEnv responseEnvelope
	if err := json.Unmarshal(body, &respEnv); err != nil {
		return nil, fmt.Errorf("transport: decode envelope: %w", err)
	}

	// Обрабатываем ошибку VK API
	if respEnv.Error != nil {
		return nil, MapError(respEnv.Error)
	}

	return &respEnv, nil
}

// encodeParams кодирует параметры запроса.
func (t *Transport) encodeParams(params url.Values) string {
	values := make(url.Values, len(params)+4)

	// Копируем параметры
	for k, v := range params {
		values[k] = v
	}

	// Добавляем обязательные параметры VK API
	values.Set("v", t.version)
	if t.lang != "" {
		values.Set("lang", t.lang)
	}
	if t.testMode {
		values.Set("test_mode", "1")
	}
	if t.token != "" && t.tokenSource == config.TokenInParams {
		values.Set("access_token", t.token)
	}

	return values.Encode()
}

// responseEnvelope представляет ответ VK API.
type responseEnvelope struct {
	Response json.RawMessage `json:"response"`
	Error    *vkErrorEnvelope `json:"error"`
}

// vkErrorEnvelope представляет ошибку VK API.
type vkErrorEnvelope struct {
	Code             int            `json:"error_code"`
	Message          string         `json:"error_msg"`
	RequestParams    []RequestParam `json:"request_params,omitempty"`
	CaptchaSID       string         `json:"captcha_sid,omitempty"`
	CaptchaImg       string         `json:"captcha_img,omitempty"`
	RedirectURI      string         `json:"redirect_uri,omitempty"`
	ConfirmationText string         `json:"confirmation_text,omitempty"`
}

// RequestParam представляет параметр запроса.
type RequestParam struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// HTTPError представляет HTTP-ошибку.
type HTTPError struct {
	StatusCode int
	Body       string
}

// Error реализует интерфейс error.
func (e *HTTPError) Error() string {
	return fmt.Sprintf("http error %d: %s", e.StatusCode, e.Body)
}
