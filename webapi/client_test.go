package webapi

import (
	"log/slog"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_new_should_require_deployment_id(t *testing.T) {
	_, err := New("")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "deploymentID")
}

func Test_new_should_require_auth_option(t *testing.T) {
	_, err := New("deploy-123")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "auth")
}

func Test_new_should_succeed_with_valid_config(t *testing.T) {
	c, err := New("deploy-123", withFakeAuth())

	assert.NoError(t, err)
	assert.NotNil(t, c)
	assert.Equal(t, "deploy-123", c.deploymentID)
	assert.Equal(t, defaultBaseURL, c.baseURL)
}

func Test_with_base_url_should_override_default(t *testing.T) {
	c, err := New("deploy-123", withFakeAuth(), WithBaseURL("https://custom.example.com"))

	assert.NoError(t, err)
	assert.Equal(t, "https://custom.example.com", c.baseURL)
}

func Test_with_http_client_should_be_used(t *testing.T) {
	custom := &http.Client{Timeout: 42 * time.Second}
	c, err := New("deploy-123", withFakeAuth(), WithHTTPClient(custom))

	assert.NoError(t, err)
	assert.Equal(t, custom, c.httpClient)
}

func Test_with_logger_should_be_set(t *testing.T) {
	logger := slog.Default()
	c, err := New("deploy-123", withFakeAuth(), WithLogger(logger))

	assert.NoError(t, err)
	assert.Equal(t, logger, c.logger)
}

func Test_with_rate_limit_should_configure_limiter(t *testing.T) {
	c, err := New("deploy-123", withFakeAuth(), WithRateLimit(50, 100))

	assert.NoError(t, err)
	assert.InDelta(t, 50.0, float64(c.limiter.Limit()), 0.01)
	assert.Equal(t, 100, c.limiter.Burst())
}

func Test_with_retry_policy_should_override_default(t *testing.T) {
	p := RetryPolicy{
		MaxRetries:  5,
		BaseDelay:   1 * time.Second,
		MaxDelay:    60 * time.Second,
		JitterRatio: 0.3,
	}
	c, err := New("deploy-123", withFakeAuth(), WithRetryPolicy(p))

	assert.NoError(t, err)
	assert.Equal(t, 5, c.retryPolicy.MaxRetries)
	assert.Equal(t, 1*time.Second, c.retryPolicy.BaseDelay)
}

func Test_default_retry_policy_should_have_sane_defaults(t *testing.T) {
	c, err := New("deploy-123", withFakeAuth())

	assert.NoError(t, err)
	assert.Equal(t, 3, c.retryPolicy.MaxRetries)
	assert.Equal(t, 500*time.Millisecond, c.retryPolicy.BaseDelay)
}

// withFakeAuth is a test helper that sets a non-nil tokenSource to
// satisfy the auth-required validation.
func withFakeAuth() Option {
	return func(c *Client) {
		c.tokenSource = staticTokenSource("test-token")
	}
}
