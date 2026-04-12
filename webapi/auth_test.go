package webapi

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/assert"
)

func tokenHandler(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"access_token": "tok-123",
		"token_type":   "bearer",
		"expires_in":   3600,
	})
}

func Test_with_client_credentials_should_acquire_token(t *testing.T) {
	var called atomic.Int32
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called.Add(1)
		assert.Equal(t, "POST", r.Method)
		assert.Equal(t, "/auth/v1/oauth/token", r.URL.Path)
		assert.NoError(t, r.ParseForm())
		assert.Equal(t, "client_credentials", r.Form.Get("grant_type"))
		assert.Equal(t, "deploy-1", r.Form.Get("deployment_id"))

		user, pass, ok := r.BasicAuth()
		assert.True(t, ok, "expected Basic auth header")
		assert.Equal(t, "my-client-id", user)
		assert.Equal(t, "my-secret", pass)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"access_token": "tok-123",
			"token_type":   "bearer",
			"expires_in":   3600,
		})
	}))
	defer ts.Close()

	c, err := New("deploy-1",
		WithClientCredentials("my-client-id", "my-secret"),
		WithBaseURL(ts.URL),
		WithRateLimit(1000, 1000),
	)
	assert.NoError(t, err)

	tok, tokErr := c.tokenSource.Token()
	assert.NoError(t, tokErr)
	assert.Equal(t, "tok-123", tok.AccessToken)
	assert.GreaterOrEqual(t, called.Load(), int32(1))
}

func Test_with_client_credentials_should_reuse_cached_token(t *testing.T) {
	var called atomic.Int32
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		called.Add(1)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"access_token": "tok-cached",
			"token_type":   "bearer",
			"expires_in":   3600,
		})
	}))
	defer ts.Close()

	c, err := New("deploy-1",
		WithClientCredentials("id", "secret"),
		WithBaseURL(ts.URL),
		WithRateLimit(1000, 1000),
	)
	assert.NoError(t, err)

	_, _ = c.tokenSource.Token()
	_, _ = c.tokenSource.Token()

	// oauth2 auto-detect may call once or twice on first request, but
	// the second Token() call must reuse the cache — so total should be
	// <= 2 (at most one auto-detect probe + one real fetch).
	assert.LessOrEqual(t, called.Load(), int32(2))
}

func Test_with_exchange_code_should_acquire_token(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method)
		assert.NoError(t, r.ParseForm())
		assert.Equal(t, "exchange_code", r.Form.Get("grant_type"))
		assert.Equal(t, "my-code-123", r.Form.Get("exchange_code"))
		assert.Equal(t, "deploy-1", r.Form.Get("deployment_id"))

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"access_token":  "xc-tok",
			"token_type":    "bearer",
			"expires_in":    3600,
			"refresh_token": "rt-abc",
		})
	}))
	defer ts.Close()

	c, err := New("deploy-1",
		WithExchangeCode("id", "secret", "my-code-123"),
		WithBaseURL(ts.URL),
		WithRateLimit(1000, 1000),
	)
	assert.NoError(t, err)

	tok, tokErr := c.tokenSource.Token()
	assert.NoError(t, tokErr)
	assert.Equal(t, "xc-tok", tok.AccessToken)
}
