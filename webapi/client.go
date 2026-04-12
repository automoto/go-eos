package webapi

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/time/rate"
)

const defaultBaseURL = "https://api.epicgames.dev"

// Client is a pure-Go client for the EOS Web API.
// Safe for concurrent use from multiple goroutines.
type Client struct {
	httpClient   *http.Client
	baseURL      string
	deploymentID string
	limiter      *rate.Limiter
	logger       *slog.Logger
	retryPolicy  RetryPolicy
	tokenSource  oauth2.TokenSource
	authBuilder  func(c *Client) // deferred auth setup; runs after all options are applied
}

// RetryPolicy configures exponential backoff with jitter for transient errors.
type RetryPolicy struct {
	MaxRetries  int
	BaseDelay   time.Duration
	MaxDelay    time.Duration
	JitterRatio float64
}

// Option configures a Client.
type Option func(*Client)

// WithHTTPClient sets a custom *http.Client for the underlying transport.
func WithHTTPClient(c *http.Client) Option {
	return func(cl *Client) { cl.httpClient = c }
}

// WithBaseURL overrides the default EOS Web API base URL.
func WithBaseURL(url string) Option {
	return func(c *Client) { c.baseURL = url }
}

// WithLogger sets a structured logger for request/response logging.
func WithLogger(l *slog.Logger) Option {
	return func(c *Client) { c.logger = l }
}

// WithRateLimit sets the proactive rate limiter parameters.
func WithRateLimit(rps float64, burst int) Option {
	return func(c *Client) { c.limiter = rate.NewLimiter(rate.Limit(rps), burst) }
}

// WithRetryPolicy sets the retry policy for transient errors.
func WithRetryPolicy(p RetryPolicy) Option {
	return func(c *Client) { c.retryPolicy = p }
}

// New creates a new Web API client. deploymentID is required.
// At least one auth option (WithClientCredentials or WithExchangeCode)
// must be provided.
func New(deploymentID string, opts ...Option) (*Client, error) {
	if deploymentID == "" {
		return nil, fmt.Errorf("webapi: deploymentID is required")
	}

	c := &Client{
		httpClient:   http.DefaultClient,
		baseURL:      defaultBaseURL,
		deploymentID: deploymentID,
		limiter:      rate.NewLimiter(10, 20),
		logger:       slog.Default(),
		retryPolicy: RetryPolicy{
			MaxRetries:  3,
			BaseDelay:   500 * time.Millisecond,
			MaxDelay:    30 * time.Second,
			JitterRatio: 0.5,
		},
	}

	for _, opt := range opts {
		opt(c)
	}

	if c.authBuilder != nil {
		c.authBuilder(c)
	}

	if c.tokenSource == nil {
		return nil, fmt.Errorf("webapi: auth option required (use WithClientCredentials or WithExchangeCode)")
	}

	return c, nil
}

// do executes an authenticated JSON API request using the client's token source.
func (c *Client) do(ctx context.Context, method, path string, body, result any) error {
	return c.doRequest(ctx, method, path, "", body, result)
}

// doWithBearerToken is like do but uses an explicit bearer token instead of
// the client's token source. Used by VerifyToken to verify a third-party token.
func (c *Client) doWithBearerToken(ctx context.Context, method, path, token string, body, result any) error {
	return c.doRequest(ctx, method, path, token, body, result)
}

// doRequest is the shared implementation for do and doWithBearerToken.
// If tokenOverride is non-empty, it is used as the Bearer token; otherwise
// the client's token source provides it.
func (c *Client) doRequest(ctx context.Context, method, path, tokenOverride string, body, result any) error {
	if err := c.limiter.Wait(ctx); err != nil {
		return err
	}

	var bodyBytes []byte
	if body != nil {
		var err error
		bodyBytes, err = json.Marshal(body)
		if err != nil {
			return fmt.Errorf("webapi: marshal request body: %w", err)
		}
	}

	url := c.baseURL + path
	resp, err := c.doWithRetry(ctx, method, url, bodyBytes, tokenOverride)
	if err != nil {
		return err
	}

	// parseErrorResponse closes the body itself.
	if resp.StatusCode >= 400 {
		return parseErrorResponse(resp)
	}
	defer resp.Body.Close() //nolint:errcheck

	if result == nil {
		return nil
	}

	return json.NewDecoder(resp.Body).Decode(result)
}

func (c *Client) buildRequest(ctx context.Context, method, url string, body []byte, tokenOverride string) (*http.Request, error) {
	var bodyReader io.Reader
	if body != nil {
		bodyReader = bytes.NewReader(body)
	}

	req, err := http.NewRequestWithContext(ctx, method, url, bodyReader)
	if err != nil {
		return nil, err
	}

	if tokenOverride != "" {
		req.Header.Set("Authorization", "Bearer "+tokenOverride)
	} else {
		tok, err := c.tokenSource.Token()
		if err != nil {
			return nil, fmt.Errorf("webapi: token: %w", err)
		}
		req.Header.Set("Authorization", "Bearer "+tok.AccessToken)
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	req.Header.Set("Accept", "application/json")

	return req, nil
}

// staticTokenSource returns a token source that always returns the
// same access token. Used in tests.
func staticTokenSource(token string) oauth2.TokenSource {
	return oauth2.StaticTokenSource(&oauth2.Token{
		AccessToken: token,
		TokenType:   "bearer",
		Expiry:      time.Now().Add(1 * time.Hour),
	})
}
