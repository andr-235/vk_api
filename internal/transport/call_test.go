package transport

import (
	"context"
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"
)

type roundTripFunc func(req *http.Request) (*http.Response, error)

func (f roundTripFunc) Do(req *http.Request) (*http.Response, error) {
	return f(req)
}

func newHTTPResponse(status int, body string) *http.Response {
	return &http.Response{
		StatusCode: status,
		Header:     make(http.Header),
		Body:       io.NopCloser(strings.NewReader(body)),
	}
}

type testParams struct {
	OwnerID int `url:"owner_id,omitempty"`
	Count   int `url:"count,omitempty"`
}

type testResponse struct {
	Count int `json:"count"`
}

func TestCall_Success(t *testing.T) {
	tests := []struct {
		name      string
		cfg       Config
		params    any
		method    string
		checkReq  func(t *testing.T, req *http.Request)
		wantCount int
	}{
		{
			name: "token in params",
			cfg: Config{
				BaseURL:     "https://api.vk.com/method/",
				Version:     "5.199",
				Lang:        "ru",
				TestMode:    true,
				Token:       "test-token",
				TokenSource: TokenInParams,
			},
			method: "users.get",
			params: testParams{
				OwnerID: 1,
				Count:   10,
			},
			checkReq: func(t *testing.T, req *http.Request) {
				t.Helper()

				if req.Method != http.MethodPost {
					t.Fatalf("expected method %q, got %q", http.MethodPost, req.Method)
				}

				if got := req.URL.Path; got != "/method/users.get" {
					t.Fatalf("expected path %q, got %q", "/method/users.get", got)
				}

				q := req.URL.Query()

				if got := q.Get("owner_id"); got != "1" {
					t.Fatalf("expected owner_id=1, got %q", got)
				}
				if got := q.Get("count"); got != "10" {
					t.Fatalf("expected count=10, got %q", got)
				}
				if got := q.Get("v"); got != "5.199" {
					t.Fatalf("expected v=5.199, got %q", got)
				}
				if got := q.Get("lang"); got != "ru" {
					t.Fatalf("expected lang=ru, got %q", got)
				}
				if got := q.Get("test_mode"); got != "1" {
					t.Fatalf("expected test_mode=1, got %q", got)
				}
				if got := q.Get("access_token"); got != "test-token" {
					t.Fatalf("expected access_token=test-token, got %q", got)
				}

				if got := req.Header.Get("Accept"); got != "application/json" {
					t.Fatalf("expected Accept header application/json, got %q", got)
				}

				if got := req.Header.Get("Authorization"); got != "" {
					t.Fatalf("expected empty Authorization header, got %q", got)
				}
			},
			wantCount: 42,
		},
		{
			name: "token in header",
			cfg: Config{
				BaseURL:     "https://api.vk.com/method/",
				Version:     "5.199",
				Token:       "header-token",
				TokenSource: TokenInHeader,
			},
			method: "users.get",
			params: testParams{
				OwnerID: 1,
			},
			checkReq: func(t *testing.T, req *http.Request) {
				t.Helper()

				if got := req.URL.Path; got != "/method/users.get" {
					t.Fatalf("expected path %q, got %q", "/method/users.get", got)
				}

				q := req.URL.Query()

				if got := q.Get("access_token"); got != "" {
					t.Fatalf("expected no access_token in query, got %q", got)
				}

				if got := q.Get("v"); got != "5.199" {
					t.Fatalf("expected v=5.199, got %q", got)
				}

				if got := req.Header.Get("Authorization"); got != "Bearer header-token" {
					t.Fatalf("expected Authorization header %q, got %q", "Bearer header-token", got)
				}

				if got := req.Header.Get("Accept"); got != "application/json" {
					t.Fatalf("expected Accept header application/json, got %q", got)
				}
			},
			wantCount: 42,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := tt.cfg
			cfg.HTTPClient = roundTripFunc(func(req *http.Request) (*http.Response, error) {
				tt.checkReq(t, req)
				return newHTTPResponse(200, `{"response":{"count":42}}`), nil
			})

			var out testResponse
			err := Call(context.Background(), cfg, tt.method, tt.params, &out)
			if err != nil {
				t.Fatalf("expected nil error, got %v", err)
			}

			if out.Count != tt.wantCount {
				t.Fatalf("expected count=%d, got %d", tt.wantCount, out.Count)
			}
		})
	}
}

func TestCall_ConfigValidation(t *testing.T) {
	tests := []struct {
		name    string
		cfg     Config
		method  string
		wantErr string
	}{
		{
			name: "missing version",
			cfg: Config{
				BaseURL: "https://api.vk.com/method/",
				HTTPClient: roundTripFunc(func(req *http.Request) (*http.Response, error) {
					return nil, nil
				}),
			},
			method:  "users.get",
			wantErr: "vk: api version is required",
		},
		{
			name: "missing base url",
			cfg: Config{
				Version: "5.199",
				HTTPClient: roundTripFunc(func(req *http.Request) (*http.Response, error) {
					return nil, nil
				}),
			},
			method:  "users.get",
			wantErr: "vk: base url is required",
		},
		{
			name: "missing http client",
			cfg: Config{
				BaseURL: "https://api.vk.com/method/",
				Version: "5.199",
			},
			method:  "users.get",
			wantErr: "vk: http client is required",
		},
		{
			name: "missing method",
			cfg: Config{
				BaseURL: "https://api.vk.com/method/",
				Version: "5.199",
				HTTPClient: roundTripFunc(func(req *http.Request) (*http.Response, error) {
					return nil, nil
				}),
			},
			method:  "",
			wantErr: "vk: method is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Call(context.Background(), tt.cfg, tt.method, nil, nil)
			if err == nil {
				t.Fatal("expected error, got nil")
			}
			if err.Error() != tt.wantErr {
				t.Fatalf("expected error %q, got %q", tt.wantErr, err.Error())
			}
		})
	}
}

func TestCall_ErrorCases(t *testing.T) {
	tests := []struct {
		name      string
		resp      *http.Response
		doErr     error
		out       any
		wantErr   string
		errPrefix bool
		checkErr  func(t *testing.T, err error)
	}{
		{
			name:    "do request error",
			doErr:   errors.New("network down"),
			wantErr: "vk: do request: network down",
		},
		{
			name:    "unexpected http status",
			resp:    newHTTPResponse(500, `internal server error`),
			wantErr: "vk: unexpected http status 500: internal server error",
		},
		{
			name:      "decode envelope error",
			resp:      newHTTPResponse(200, `{"response":`),
			wantErr:   "vk: decode envelope: ",
			errPrefix: true,
		},
		{
			name: "api error",
			resp: newHTTPResponse(200, `{
				"error": {
					"error_code": 5,
					"error_msg": "User authorization failed"
				}
			}`),
			checkErr: func(t *testing.T, err error) {
				t.Helper()

				var apiErr *APIError
				if !errors.As(err, &apiErr) {
					t.Fatalf("expected APIError, got %T", err)
				}

				if apiErr.Code != 5 {
					t.Fatalf("expected code=5, got %d", apiErr.Code)
				}

				if apiErr.Message != "User authorization failed" {
					t.Fatalf("expected message %q, got %q", "User authorization failed", apiErr.Message)
				}
			},
		},
		{
			name:      "decode response payload error",
			resp:      newHTTPResponse(200, `{"response":"not-an-object"}`),
			out:       &testResponse{},
			wantErr:   "vk: decode response payload: ",
			errPrefix: true,
		},
		{
			name: "nil out ignores response payload",
			resp: newHTTPResponse(200, `{"response":{"count":123}}`),
			out:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := Config{
				BaseURL: "https://api.vk.com/method/",
				Version: "5.199",
				HTTPClient: roundTripFunc(func(req *http.Request) (*http.Response, error) {
					if tt.doErr != nil {
						return nil, tt.doErr
					}
					return tt.resp, nil
				}),
			}

			err := Call(context.Background(), cfg, "users.get", nil, tt.out)

			if tt.wantErr == "" && tt.checkErr == nil {
				if err != nil {
					t.Fatalf("expected nil error, got %v", err)
				}
				return
			}

			if err == nil {
				t.Fatal("expected error, got nil")
			}

			if tt.errPrefix {
				if !strings.HasPrefix(err.Error(), tt.wantErr) {
					t.Fatalf("expected error prefix %q, got %q", tt.wantErr, err.Error())
				}
			} else if tt.wantErr != "" {
				if err.Error() != tt.wantErr {
					t.Fatalf("expected error %q, got %q", tt.wantErr, err.Error())
				}
			}

			if tt.checkErr != nil {
				tt.checkErr(t, err)
			}
		})
	}
}

func TestAPIError_Error(t *testing.T) {
	tests := []struct {
		name string
		err  *APIError
		want string
	}{
		{
			name: "nil",
			err:  nil,
			want: "<nil>",
		},
		{
			name: "without message",
			err: &APIError{
				Code: 123,
			},
			want: "vk api error 123",
		},
		{
			name: "with message",
			err: &APIError{
				Code:    5,
				Message: "User authorization failed",
			},
			want: "vk api error 5: User authorization failed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.err.Error(); got != tt.want {
				t.Fatalf("expected %q, got %q", tt.want, got)
			}
		})
	}
}

func TestQueryFromParams(t *testing.T) {
	tests := []struct {
		name      string
		params    any
		checkFunc func(t *testing.T, values map[string]string)
	}{
		{
			name: "encodes params",
			params: testParams{
				OwnerID: 1,
				Count:   10,
			},
			checkFunc: func(t *testing.T, values map[string]string) {
				t.Helper()

				if got := values["owner_id"]; got != "1" {
					t.Fatalf("expected owner_id=1, got %q", got)
				}
				if got := values["count"]; got != "10" {
					t.Fatalf("expected count=10, got %q", got)
				}
			},
		},
		{
			name:   "nil params",
			params: nil,
			checkFunc: func(t *testing.T, values map[string]string) {
				t.Helper()

				if len(values) != 0 {
					t.Fatalf("expected empty values, got %#v", values)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := QueryFromParams(tt.params)
			if err != nil {
				t.Fatalf("expected nil error, got %v", err)
			}

			values := make(map[string]string, len(got))
			for k := range got {
				values[k] = got.Get(k)
			}

			tt.checkFunc(t, values)
		})
	}
}
