package webapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// APIError represents an error response from the EOS Web API.
type APIError struct {
	HTTPStatus int
	ErrorCode  string
	Message    string
}

// Error returns a human-readable representation of the API error.
func (e *APIError) Error() string {
	if e.ErrorCode != "" {
		return fmt.Sprintf("webapi: HTTP %d: %s: %s", e.HTTPStatus, e.ErrorCode, e.Message)
	}
	return fmt.Sprintf("webapi: HTTP %d: %s", e.HTTPStatus, e.Message)
}

// Is supports errors.Is matching by HTTP status code. Sentinels like
// ErrRateLimited carry only a status code; any APIError with the same
// status matches.
func (e *APIError) Is(target error) bool {
	t, ok := target.(*APIError)
	if !ok {
		return false
	}
	if t.ErrorCode == "" {
		return e.HTTPStatus == t.HTTPStatus
	}
	return e.HTTPStatus == t.HTTPStatus && e.ErrorCode == t.ErrorCode
}

// Sentinel errors for common HTTP status codes. Use errors.Is to match.
var (
	// ErrUnauthorized indicates a 401 Unauthorized response.
	ErrUnauthorized = &APIError{HTTPStatus: 401}
	// ErrForbidden indicates a 403 Forbidden response.
	ErrForbidden = &APIError{HTTPStatus: 403}
	// ErrNotFound indicates a 404 Not Found response.
	ErrNotFound = &APIError{HTTPStatus: 404}
	// ErrRateLimited indicates a 429 Too Many Requests response.
	ErrRateLimited = &APIError{HTTPStatus: 429}
	// ErrServerError indicates a 500 Internal Server Error response.
	ErrServerError = &APIError{HTTPStatus: 500}
)

type errorBody struct {
	ErrorCode    string `json:"errorCode"`
	ErrorMessage string `json:"errorMessage"`
}

func parseErrorResponse(resp *http.Response) *APIError {
	defer resp.Body.Close() //nolint:errcheck

	data, err := io.ReadAll(resp.Body)
	if err != nil || len(data) == 0 {
		return &APIError{
			HTTPStatus: resp.StatusCode,
			Message:    http.StatusText(resp.StatusCode),
		}
	}

	var body errorBody
	if json.Unmarshal(data, &body) != nil || body.ErrorMessage == "" {
		return &APIError{
			HTTPStatus: resp.StatusCode,
			Message:    http.StatusText(resp.StatusCode),
		}
	}

	return &APIError{
		HTTPStatus: resp.StatusCode,
		ErrorCode:  body.ErrorCode,
		Message:    body.ErrorMessage,
	}
}
