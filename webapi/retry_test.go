package webapi

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_retry_should_succeed_after_transient_errors(t *testing.T) {
	var calls atomic.Int32
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		n := calls.Add(1)
		if r.URL.Path == "/auth/v1/oauth/token" {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]any{
				"access_token": "tok", "token_type": "bearer", "expires_in": 3600,
			})
			return
		}
		if n <= 2 {
			w.WriteHeader(500)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"ok":true}`))
	}))
	defer ts.Close()

	c, err := New("deploy-1",
		WithClientCredentials("id", "secret"),
		WithBaseURL(ts.URL),
		WithRateLimit(1000, 1000),
		WithRetryPolicy(RetryPolicy{
			MaxRetries: 3, BaseDelay: 1 * time.Millisecond,
			MaxDelay: 10 * time.Millisecond, JitterRatio: 0,
		}),
	)
	assert.NoError(t, err)

	var result map[string]bool
	doErr := c.doGet(context.Background(), "/test", &result)

	assert.NoError(t, doErr)
	assert.True(t, result["ok"])
}

func Test_retry_should_return_error_when_max_retries_exceeded(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/auth/v1/oauth/token" {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]any{
				"access_token": "tok", "token_type": "bearer", "expires_in": 3600,
			})
			return
		}
		w.WriteHeader(503)
	}))
	defer ts.Close()

	c, err := New("deploy-1",
		WithClientCredentials("id", "secret"),
		WithBaseURL(ts.URL),
		WithRateLimit(1000, 1000),
		WithRetryPolicy(RetryPolicy{
			MaxRetries: 2, BaseDelay: 1 * time.Millisecond,
			MaxDelay: 10 * time.Millisecond, JitterRatio: 0,
		}),
	)
	assert.NoError(t, err)

	doErr := c.doGet(context.Background(), "/fail", nil)
	assert.Error(t, doErr)
	assert.Contains(t, doErr.Error(), "max retries exceeded")
}

func Test_retry_should_handle_429_with_retry_after(t *testing.T) {
	var calls atomic.Int32
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/auth/v1/oauth/token" {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]any{
				"access_token": "tok", "token_type": "bearer", "expires_in": 3600,
			})
			return
		}
		n := calls.Add(1)
		if n == 1 {
			w.Header().Set("Retry-After", "0")
			w.WriteHeader(429)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"ok":true}`))
	}))
	defer ts.Close()

	c, err := New("deploy-1",
		WithClientCredentials("id", "secret"),
		WithBaseURL(ts.URL),
		WithRateLimit(1000, 1000),
		WithRetryPolicy(RetryPolicy{
			MaxRetries: 3, BaseDelay: 1 * time.Millisecond,
			MaxDelay: 10 * time.Millisecond, JitterRatio: 0,
		}),
	)
	assert.NoError(t, err)

	var result map[string]bool
	doErr := c.doGet(context.Background(), "/throttled", &result)

	assert.NoError(t, doErr)
	assert.True(t, result["ok"])
	assert.Equal(t, int32(2), calls.Load())
}

func Test_retry_should_cancel_on_context(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/auth/v1/oauth/token" {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]any{
				"access_token": "tok", "token_type": "bearer", "expires_in": 3600,
			})
			return
		}
		w.WriteHeader(500)
	}))
	defer ts.Close()

	c, err := New("deploy-1",
		WithClientCredentials("id", "secret"),
		WithBaseURL(ts.URL),
		WithRateLimit(1000, 1000),
		WithRetryPolicy(RetryPolicy{
			MaxRetries: 10, BaseDelay: 10 * time.Second,
			MaxDelay: 30 * time.Second, JitterRatio: 0,
		}),
	)
	assert.NoError(t, err)

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	doErr := c.doGet(ctx, "/slow", nil)
	assert.Error(t, doErr)
	assert.ErrorIs(t, doErr, context.DeadlineExceeded)
}

func Test_backoff_delay_should_grow_exponentially(t *testing.T) {
	c := &Client{retryPolicy: RetryPolicy{
		BaseDelay: 100 * time.Millisecond, MaxDelay: 10 * time.Second, JitterRatio: 0,
	}}

	d1 := c.backoffDelay(1)
	d2 := c.backoffDelay(2)
	d3 := c.backoffDelay(3)

	assert.Equal(t, 100*time.Millisecond, d1)
	assert.Equal(t, 200*time.Millisecond, d2)
	assert.Equal(t, 400*time.Millisecond, d3)
}

func Test_backoff_delay_should_cap_at_max(t *testing.T) {
	c := &Client{retryPolicy: RetryPolicy{
		BaseDelay: 1 * time.Second, MaxDelay: 5 * time.Second, JitterRatio: 0,
	}}

	d := c.backoffDelay(10)
	assert.Equal(t, 5*time.Second, d)
}

func Test_parse_retry_after_should_handle_seconds(t *testing.T) {
	assert.Equal(t, 5*time.Second, parseRetryAfter("5"))
}

func Test_parse_retry_after_should_default_on_empty(t *testing.T) {
	assert.Equal(t, 1*time.Second, parseRetryAfter(""))
}

func Test_is_transient_error_should_match_5xx(t *testing.T) {
	tests := []struct {
		status    int
		transient bool
	}{
		{200, false},
		{400, false},
		{401, false},
		{429, false},
		{500, true},
		{502, true},
		{503, true},
		{504, true},
	}
	for _, tt := range tests {
		assert.Equal(t, tt.transient, isTransientError(tt.status), "status %d", tt.status)
	}
}
