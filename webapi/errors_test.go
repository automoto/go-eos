package webapi

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_api_error_should_format_with_error_code(t *testing.T) {
	err := &APIError{
		HTTPStatus: 429,
		ErrorCode:  "errors.com.epicgames.common.throttled",
		Message:    "Rate limit exceeded",
	}

	assert.Equal(t, "webapi: HTTP 429: errors.com.epicgames.common.throttled: Rate limit exceeded", err.Error())
}

func Test_api_error_should_format_without_error_code(t *testing.T) {
	err := &APIError{
		HTTPStatus: 500,
		Message:    "Internal Server Error",
	}

	assert.Equal(t, "webapi: HTTP 500: Internal Server Error", err.Error())
}

func Test_api_error_should_match_sentinel_by_status_code(t *testing.T) {
	err := &APIError{
		HTTPStatus: 429,
		ErrorCode:  "errors.com.epicgames.common.throttled",
		Message:    "Rate limit exceeded",
	}

	assert.True(t, errors.Is(err, ErrRateLimited))
	assert.False(t, errors.Is(err, ErrUnauthorized))
}

func Test_api_error_should_match_all_sentinels(t *testing.T) {
	tests := []struct {
		name     string
		status   int
		sentinel *APIError
	}{
		{"unauthorized", 401, ErrUnauthorized},
		{"forbidden", 403, ErrForbidden},
		{"not found", 404, ErrNotFound},
		{"rate limited", 429, ErrRateLimited},
		{"server error", 500, ErrServerError},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := &APIError{HTTPStatus: tt.status, Message: "test"}
			assert.True(t, errors.Is(err, tt.sentinel))
		})
	}
}

func Test_api_error_should_support_errors_as(t *testing.T) {
	err := &APIError{HTTPStatus: 401, ErrorCode: "invalid_token", Message: "expired"}
	wrapped := fmt.Errorf("auth failed: %w", err)

	var apiErr *APIError
	assert.True(t, errors.As(wrapped, &apiErr))
	assert.Equal(t, 401, apiErr.HTTPStatus)
	assert.Equal(t, "invalid_token", apiErr.ErrorCode)
	assert.Equal(t, "expired", apiErr.Message)
}

func Test_parse_error_response_should_decode_json_body(t *testing.T) {
	body := `{"errorCode":"errors.com.epicgames.common.throttled","errorMessage":"Rate limit exceeded"}`
	resp := &http.Response{
		StatusCode: 429,
		Body:       io.NopCloser(strings.NewReader(body)),
	}

	err := parseErrorResponse(resp)

	assert.Equal(t, 429, err.HTTPStatus)
	assert.Equal(t, "errors.com.epicgames.common.throttled", err.ErrorCode)
	assert.Equal(t, "Rate limit exceeded", err.Message)
}

func Test_parse_error_response_should_fallback_on_invalid_json(t *testing.T) {
	resp := &http.Response{
		StatusCode: 503,
		Body:       io.NopCloser(strings.NewReader("not json")),
	}

	err := parseErrorResponse(resp)

	assert.Equal(t, 503, err.HTTPStatus)
	assert.Equal(t, "", err.ErrorCode)
	assert.Equal(t, "Service Unavailable", err.Message)
}

func Test_parse_error_response_should_handle_empty_body(t *testing.T) {
	resp := &http.Response{
		StatusCode: 500,
		Body:       io.NopCloser(strings.NewReader("")),
	}

	err := parseErrorResponse(resp)

	assert.Equal(t, 500, err.HTTPStatus)
	assert.Equal(t, "Internal Server Error", err.Message)
}
