package webapi

import (
	"context"
	"fmt"
	"math"
	"math/rand/v2"
	"net/http"
	"strconv"
	"time"
)

func (c *Client) doWithRetry(ctx context.Context, method, url string, body []byte, tokenOverride string) (*http.Response, error) {
	var lastErr error
	for attempt := range c.retryPolicy.MaxRetries + 1 {
		if attempt > 0 {
			delay := c.backoffDelay(attempt)
			select {
			case <-time.After(delay):
			case <-ctx.Done():
				return nil, ctx.Err()
			}
		}

		req, err := c.buildRequest(ctx, method, url, body, tokenOverride)
		if err != nil {
			return nil, err
		}

		start := time.Now()
		resp, err := c.httpClient.Do(req)
		dur := time.Since(start)

		if err != nil {
			c.logger.DebugContext(ctx, "request failed",
				"method", method, "url", url, "attempt", attempt, "error", err)
			lastErr = err
			continue
		}

		c.logger.DebugContext(ctx, "request complete",
			"method", method, "url", url, "status", resp.StatusCode,
			"duration", dur, "attempt", attempt)

		if resp.StatusCode == http.StatusTooManyRequests {
			delay := parseRetryAfter(resp.Header.Get("Retry-After"))
			resp.Body.Close() //nolint:errcheck
			lastErr = ErrRateLimited
			select {
			case <-time.After(delay):
			case <-ctx.Done():
				return nil, ctx.Err()
			}
			continue
		}

		if isTransientError(resp.StatusCode) {
			resp.Body.Close() //nolint:errcheck
			lastErr = &APIError{HTTPStatus: resp.StatusCode, Message: http.StatusText(resp.StatusCode)}
			continue
		}

		return resp, nil
	}
	return nil, fmt.Errorf("webapi: max retries exceeded: %w", lastErr)
}

func isTransientError(status int) bool {
	return status == 500 || status == 502 || status == 503 || status == 504
}

func (c *Client) backoffDelay(attempt int) time.Duration {
	delay := float64(c.retryPolicy.BaseDelay) * math.Pow(2, float64(attempt-1))
	if delay > float64(c.retryPolicy.MaxDelay) {
		delay = float64(c.retryPolicy.MaxDelay)
	}
	jitter := delay * c.retryPolicy.JitterRatio * (2*rand.Float64() - 1)
	d := time.Duration(delay + jitter)
	if d < 0 {
		return 0
	}
	return d
}

func parseRetryAfter(val string) time.Duration {
	if val == "" {
		return 1 * time.Second
	}
	if secs, err := strconv.Atoi(val); err == nil {
		return time.Duration(secs) * time.Second
	}
	if t, err := time.Parse(time.RFC1123, val); err == nil {
		d := time.Until(t)
		if d > 0 {
			return d
		}
	}
	return 1 * time.Second
}

// doGet is a convenience wrapper for GET requests without a body.
func (c *Client) doGet(ctx context.Context, path string, result any) error {
	return c.do(ctx, http.MethodGet, path, nil, result)
}
