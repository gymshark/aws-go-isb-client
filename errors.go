package isbclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// APIRequestError wraps errors related to making API requests.
type APIRequestError struct {
	Op  string
	URL string
	Err error
}

func (e *APIRequestError) Error() string {
	return fmt.Sprintf("api request error [%s %s]: %v", e.Op, e.URL, e.Err)
}

func (e *APIRequestError) Unwrap() error {
	return e.Err
}

// APIResponseError is a generic error for unexpected HTTP responses.
type APIResponseError struct {
	StatusCode int
	Body       string
	Message    string
}

func (e *APIResponseError) Error() string {
	if e.Message != "" {
		return fmt.Sprintf("api response error: %s (status %d)", e.Message, e.StatusCode)
	}
	return fmt.Sprintf("api response error: status %d, body: %s", e.StatusCode, e.Body)
}

// FailErrorDetail represents the structure of errors in a fail response.
type FailErrorDetail struct {
	Message string `json:"message"`
}

// FailResponseError represents a 'fail' response from the API (status: fail, data.errors).
type FailResponseError struct {
	Status     string // always "fail"
	Errors     []FailErrorDetail
	StatusCode int
}

func (e *FailResponseError) Error() string {
	return fmt.Sprintf("fail response: %s (status %d) errors: %v", e.Status, e.StatusCode, e.Errors)
}

// ErrorErrorDetail represents the structure of data in an error response.
type ErrorErrorDetail struct {
	Status  string                 `json:"status"`
	Message string                 `json:"message"`
	Code    int                    `json:"code,omitempty"`
	Data    map[string]interface{} `json:"data,omitempty"`
}

// ServerError represents a 500+ server error.
type ServerError struct {
	APIResponseError
	Code int
	Data map[string]interface{}
}

func (e *ServerError) Error() string {
	return fmt.Sprintf("server error: %s (status %d, code %d)", e.Message, e.StatusCode, e.Code)
}

// JSONDecodingError wraps errors related to JSON decoding.
type JSONDecodingError struct {
	Err error
}

func (e *JSONDecodingError) Error() string {
	return fmt.Sprintf("json decoding error: %v", e.Err)
}

func (e *JSONDecodingError) Unwrap() error {
	return e.Err
}

// BadRequestError represents a 400 Bad Request error.
type BadRequestError struct {
	APIResponseError
	Errors      []FailErrorDetail // from spec.yaml: data.errors
	RequestBody string
}

func (e *BadRequestError) Error() string {
	if len(e.Errors) > 0 {
		return fmt.Sprintf("bad request: %s (status %d) errors: %v body: %s", e.Message, e.StatusCode, e.Errors, e.RequestBody)
	}
	return fmt.Sprintf("bad request: %s (status %d) body: %s", e.Message, e.StatusCode, e.RequestBody)
}

// UnauthorizedError represents a 401 or 403 Unauthorized error.
type UnauthorizedError struct {
	APIResponseError
}

func (e *UnauthorizedError) Error() string {
	return fmt.Sprintf("unauthorized: %s (status %d)", e.Message, e.StatusCode)
}

// NotFoundError represents a 404 Not Found error.
type NotFoundError struct {
	APIResponseError
	Errors []FailErrorDetail // from spec.yaml: data.errors
}

func (e *NotFoundError) Error() string {
	if len(e.Errors) > 0 {
		return fmt.Sprintf("not found: %s (status %d) errors: %v", e.Message, e.StatusCode, e.Errors)
	}
	return fmt.Sprintf("not found: %s (status %d)", e.Message, e.StatusCode)
}

// ConflictError represents a 409 Conflict error.
type ConflictError struct {
	APIResponseError
	Errors []FailErrorDetail // from spec.yaml: data.errors
}

func (e *ConflictError) Error() string {
	if len(e.Errors) > 0 {
		return fmt.Sprintf("conflict: %s (status %d) errors: %v", e.Message, e.StatusCode, e.Errors)
	}
	return fmt.Sprintf("conflict: %s (status %d)", e.Message, e.StatusCode)
}

// LeaseNotFoundError represents a 404 Not Found error for a lease resource.
type LeaseNotFoundError struct {
	APIResponseError
	Errors []FailErrorDetail
}

func (e *LeaseNotFoundError) Error() string {
	if len(e.Errors) > 0 {
		return fmt.Sprintf("lease not found: %s (status %d) errors: %v", e.Message, e.StatusCode, e.Errors)
	}
	return fmt.Sprintf("lease not found: %s (status %d)", e.Message, e.StatusCode)
}

// LeaseTemplateNotFoundError represents a 404 Not Found error for a lease template resource.
type LeaseTemplateNotFoundError struct {
	APIResponseError
	Errors []FailErrorDetail
}

func (e *LeaseTemplateNotFoundError) Error() string {
	if len(e.Errors) > 0 {
		return fmt.Sprintf("lease template not found: %s (status %d) errors: %v", e.Message, e.StatusCode, e.Errors)
	}
	return fmt.Sprintf("lease template not found: %s (status %d)", e.Message, e.StatusCode)
}

// AccountNotFoundError represents a 404 Not Found error for an account resource.
type AccountNotFoundError struct {
	APIResponseError
	Errors []FailErrorDetail
}

func (e *AccountNotFoundError) Error() string {
	if len(e.Errors) > 0 {
		return fmt.Sprintf("account not found: %s (status %d) errors: %v", e.Message, e.StatusCode, e.Errors)
	}
	return fmt.Sprintf("account not found: %s (status %d)", e.Message, e.StatusCode)
}

// LeaseConflictError represents a 409 Conflict error for a lease resource.
type LeaseConflictError struct {
	APIResponseError
	Errors []FailErrorDetail
}

func (e *LeaseConflictError) Error() string {
	if len(e.Errors) > 0 {
		return fmt.Sprintf("lease conflict: %s (status %d) errors: %v", e.Message, e.StatusCode, e.Errors)
	}
	return fmt.Sprintf("lease conflict: %s (status %d)", e.Message, e.StatusCode)
}

// LeaseTemplateConflictError represents a 409 Conflict error for a lease template resource.
type LeaseTemplateConflictError struct {
	APIResponseError
	Errors []FailErrorDetail
}

func (e *LeaseTemplateConflictError) Error() string {
	if len(e.Errors) > 0 {
		return fmt.Sprintf("lease template conflict: %s (status %d) errors: %v", e.Message, e.StatusCode, e.Errors)
	}
	return fmt.Sprintf("lease template conflict: %s (status %d)", e.Message, e.StatusCode)
}

// AccountConflictError represents a 409 Conflict error for an account resource.
type AccountConflictError struct {
	APIResponseError
	Errors []FailErrorDetail
}

func (e *AccountConflictError) Error() string {
	if len(e.Errors) > 0 {
		return fmt.Sprintf("account conflict: %s (status %d) errors: %v", e.Message, e.StatusCode, e.Errors)
	}
	return fmt.Sprintf("account conflict: %s (status %d)", e.Message, e.StatusCode)
}

// DecodeAPIError decodes the API error response and returns the appropriate error type.
func DecodeAPIError(reqBody []byte, resp *http.Response) error {
	defer resp.Body.Close()

	// Read the body into a buffer so we can decode multiple times
	b := new(bytes.Buffer)
	_, _ = b.ReadFrom(resp.Body)
	bodyBytes := b.Bytes()

	var (
		failBody struct {
			Status string `json:"status"`
			Data   struct {
				Errors []FailErrorDetail `json:"errors"`
			} `json:"data"`
		}
		errorBody struct {
			Status  string                 `json:"status"`
			Message string                 `json:"message"`
			Code    int                    `json:"code,omitempty"`
			Data    map[string]interface{} `json:"data,omitempty"`
		}
	)

	urlPath := resp.Request.URL.Path

	switch resp.StatusCode {
	case 400:
		if err := json.NewDecoder(bytes.NewReader(bodyBytes)).Decode(&failBody); err == nil && failBody.Status == "fail" {
			return &BadRequestError{
				APIResponseError: APIResponseError{StatusCode: 400, Message: "bad request"},
				Errors:           failBody.Data.Errors,
				RequestBody:      string(reqBody),
			}
		}
	case 401, 403:
		if err := json.NewDecoder(bytes.NewReader(bodyBytes)).Decode(&failBody); err == nil && failBody.Status == "fail" {
			return &UnauthorizedError{
				APIResponseError: APIResponseError{StatusCode: resp.StatusCode, Message: "unauthorized"},
			}
		}
	case 404:
		if err := json.NewDecoder(bytes.NewReader(bodyBytes)).Decode(&failBody); err == nil && failBody.Status == "fail" {
			// Resource-specific not found errors
			if strings.HasPrefix(urlPath, "/leases/") {
				return &LeaseNotFoundError{
					APIResponseError: APIResponseError{StatusCode: 404, Message: "lease not found"},
					Errors:           failBody.Data.Errors,
				}
			} else if strings.HasPrefix(urlPath, "/leaseTemplates/") {
				return &LeaseTemplateNotFoundError{
					APIResponseError: APIResponseError{StatusCode: 404, Message: "lease template not found"},
					Errors:           failBody.Data.Errors,
				}
			} else if strings.HasPrefix(urlPath, "/accounts/") {
				return &AccountNotFoundError{
					APIResponseError: APIResponseError{StatusCode: 404, Message: "account not found"},
					Errors:           failBody.Data.Errors,
				}
			}
			// fallback
			return &NotFoundError{
				APIResponseError: APIResponseError{StatusCode: 404, Message: "not found"},
				Errors:           failBody.Data.Errors,
			}
		}
	case 409:
		if err := json.NewDecoder(bytes.NewReader(bodyBytes)).Decode(&failBody); err == nil && failBody.Status == "fail" {
			// Resource-specific conflict errors
			if strings.HasPrefix(urlPath, "/leases/") {
				return &LeaseConflictError{
					APIResponseError: APIResponseError{StatusCode: 409, Message: "lease conflict"},
					Errors:           failBody.Data.Errors,
				}
			} else if strings.HasPrefix(urlPath, "/leaseTemplates/") {
				return &LeaseTemplateConflictError{
					APIResponseError: APIResponseError{StatusCode: 409, Message: "lease template conflict"},
					Errors:           failBody.Data.Errors,
				}
			} else if strings.HasPrefix(urlPath, "/accounts/") {
				return &AccountConflictError{
					APIResponseError: APIResponseError{StatusCode: 409, Message: "account conflict"},
					Errors:           failBody.Data.Errors,
				}
			}
			// fallback
			return &ConflictError{
				APIResponseError: APIResponseError{StatusCode: 409, Message: "conflict"},
				Errors:           failBody.Data.Errors,
			}
		}
	case 500:
		if err := json.NewDecoder(bytes.NewReader(bodyBytes)).Decode(&errorBody); err == nil && errorBody.Status == "error" {
			return &ServerError{
				APIResponseError: APIResponseError{StatusCode: 500, Message: errorBody.Message},
				Code:             errorBody.Code,
				Data:             errorBody.Data,
			}
		}
	}
	// fallback: try to decode as fail
	if err := json.NewDecoder(bytes.NewReader(bodyBytes)).Decode(&failBody); err == nil && failBody.Status == "fail" {
		return &FailResponseError{
			Status:     failBody.Status,
			Errors:     failBody.Data.Errors,
			StatusCode: resp.StatusCode,
		}
	}
	// fallback: try to decode as error
	if err := json.NewDecoder(bytes.NewReader(bodyBytes)).Decode(&errorBody); err == nil && errorBody.Status == "error" {
		return &ServerError{
			APIResponseError: APIResponseError{StatusCode: resp.StatusCode, Message: errorBody.Message},
			Code:             errorBody.Code,
			Data:             errorBody.Data,
		}
	}
	// fallback: generic error
	return &APIResponseError{StatusCode: resp.StatusCode, Body: string(bodyBytes)}
}
