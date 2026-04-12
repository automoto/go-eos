package webapi

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func newTestClient(t *testing.T, handlers map[string]http.HandlerFunc) *Client {
	t.Helper()
	mux := http.NewServeMux()
	mux.HandleFunc("POST /auth/v1/oauth/token", tokenHandler)
	for pattern, handler := range handlers {
		mux.HandleFunc(pattern, handler)
	}
	ts := httptest.NewServer(mux)
	t.Cleanup(ts.Close)

	c, err := New("test-deployment",
		WithClientCredentials("test-id", "test-secret"),
		WithBaseURL(ts.URL),
		WithRateLimit(1000, 1000),
		WithRetryPolicy(RetryPolicy{MaxRetries: 0}),
	)
	assert.NoError(t, err)
	return c
}

func Test_verify_token_should_return_token_info(t *testing.T) {
	c := newTestClient(t, map[string]http.HandlerFunc{
		"GET /auth/v1/oauth/verify": func(w http.ResponseWriter, r *http.Request) {
			assert.Contains(t, r.Header.Get("Authorization"), "Bearer")
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]any{
				"active":        true,
				"token_type":    "bearer",
				"expires_in":    3000,
				"account_id":    "acc-123",
				"client_id":     "cli-456",
				"deployment_id": "dep-789",
			})
		},
	})

	info, err := c.VerifyToken(context.Background(), "some-token")

	assert.NoError(t, err)
	assert.True(t, info.Active)
	assert.Equal(t, "acc-123", info.AccountID)
	assert.Equal(t, "dep-789", info.DeploymentID)
}

func Test_verify_token_should_return_error_on_401(t *testing.T) {
	c := newTestClient(t, map[string]http.HandlerFunc{
		"GET /auth/v1/oauth/verify": func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(401)
			json.NewEncoder(w).Encode(map[string]any{
				"errorCode":    "invalid_token",
				"errorMessage": "Token is expired",
			})
		},
	})

	_, err := c.VerifyToken(context.Background(), "expired-tok")

	assert.Error(t, err)
	assert.ErrorIs(t, err, ErrUnauthorized)
}

func Test_verify_token_should_require_non_empty_token(t *testing.T) {
	c := newTestClient(t, nil)

	_, err := c.VerifyToken(context.Background(), "")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "token")
}

func Test_get_accounts_should_return_account_info(t *testing.T) {
	c := newTestClient(t, map[string]http.HandlerFunc{
		"GET /auth/v1/accounts": func(w http.ResponseWriter, r *http.Request) {
			ids := r.URL.Query()["accountId"]
			assert.Equal(t, []string{"a1", "a2"}, ids)

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode([]map[string]any{
				{"accountId": "a1", "displayName": "Alice"},
				{"accountId": "a2", "displayName": "Bob"},
			})
		},
	})

	accounts, err := c.GetAccounts(context.Background(), []string{"a1", "a2"})

	assert.NoError(t, err)
	assert.Len(t, accounts, 2)
	assert.Equal(t, "Alice", accounts[0].DisplayName)
	assert.Equal(t, "Bob", accounts[1].DisplayName)
}

func Test_get_accounts_should_reject_empty_ids(t *testing.T) {
	c := newTestClient(t, nil)

	_, err := c.GetAccounts(context.Background(), nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "accountIDs")
}

func Test_get_accounts_should_reject_over_100_ids(t *testing.T) {
	c := newTestClient(t, nil)

	ids := make([]string, 101)
	for i := range ids {
		ids[i] = "id"
	}
	_, err := c.GetAccounts(context.Background(), ids)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "100")
}
