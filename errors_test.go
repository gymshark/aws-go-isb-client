package isbclient

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"testing"
)

type mockReadCloser struct {
	*bytes.Buffer
}

func (m *mockReadCloser) Close() error { return nil }

func newMockResponse(statusCode int, body interface{}) *http.Response {
	var b bytes.Buffer
	if body != nil {
		_ = json.NewEncoder(&b).Encode(body)
	}
	return &http.Response{
		StatusCode: statusCode,
		Body:       &mockReadCloser{Buffer: &b},
	}
}

func TestDecodeAPIError_Conflict(t *testing.T) {
	failBody := FailResponseBody{
		Status: "fail",
		Data: struct {
			Errors []FailErrorDetail `json:"errors"`
		}{Errors: []FailErrorDetail{{Message: "conflict error"}}},
	}
	resp := newMockResponse(409, failBody)
	resp.Request = &http.Request{URL: &url.URL{Path: "/leases/123"}}
	err := DecodeAPIError(nil, resp)
	if _, ok := err.(*LeaseConflictError); !ok {
		t.Errorf("expected LeaseConflictError, got %T", err)
	}
}

func TestDecodeAPIError_ServerError(t *testing.T) {
	errBody := ErrorResponseBody{
		Status:  "error",
		Message: "server error",
		Code:    500,
		Data:    map[string]interface{}{"foo": "bar"},
	}
	resp := newMockResponse(500, errBody)
	resp.Request = &http.Request{URL: &url.URL{Path: "/leases/123"}} // Set Request for DecodeAPIError
	err := DecodeAPIError(nil, resp)
	if serr, ok := err.(*ServerError); !ok || serr.Code != 500 {
		t.Errorf("expected ServerError with code 500, got %T %+v", err, err)
	}
}

func TestDecodeAPIError_FallbackFail(t *testing.T) {
	failBody := FailResponseBody{
		Status: "fail",
		Data: struct {
			Errors []FailErrorDetail `json:"errors"`
		}{Errors: []FailErrorDetail{{Message: "fail error"}}},
	}
	resp := newMockResponse(418, failBody)
	// Remove resource-specific path for fallback test
	resp.Request = &http.Request{URL: &url.URL{Path: "/other"}}
	err := DecodeAPIError(nil, resp)
	if ferr, ok := err.(*FailResponseError); !ok || ferr.Status != "fail" {
		t.Errorf("expected FailResponseError, got %T %+v", err, err)
	}
}

func TestDecodeAPIError_FallbackError(t *testing.T) {
	errBody := ErrorResponseBody{
		Status:  "error",
		Message: "some error",
		Code:    400,
		Data:    nil,
	}
	resp := newMockResponse(418, errBody)
	resp.Request = &http.Request{URL: &url.URL{Path: "/other"}}
	err := DecodeAPIError(nil, resp)
	if serr, ok := err.(*ServerError); !ok || serr.Code != 400 || serr.Message != "some error" || serr.StatusCode != 418 {
		t.Errorf("expected ServerError with code 400, got %T %+v", err, err)
	}
}

func TestDecodeAPIError_Generic(t *testing.T) {
	resp := &http.Response{
		StatusCode: 418,
		Body:       ioutil.NopCloser(strings.NewReader("teapot error")),
		Request:    &http.Request{URL: &url.URL{Path: "/other"}},
	}
	err := DecodeAPIError(nil, resp)
	if apiErr, ok := err.(*APIResponseError); !ok || !strings.Contains(apiErr.Body, "teapot error") {
		t.Errorf("expected APIResponseError with body, got %T %+v", err, err)
	}
}
