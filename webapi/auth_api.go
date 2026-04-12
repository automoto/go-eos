package webapi

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

// TokenInfo contains information about a verified access token.
type TokenInfo struct {
	Active        bool   `json:"active"`
	Scope         string `json:"scope"`
	TokenType     string `json:"token_type"`
	ExpiresIn     int    `json:"expires_in"`
	ExpiresAt     string `json:"expires_at"`
	AccountID     string `json:"account_id"`
	ClientID      string `json:"client_id"`
	ProductID     string `json:"product_id"`
	ApplicationID string `json:"application_id"`
	DeploymentID  string `json:"deployment_id"`
}

// VerifyToken verifies an access token and returns its metadata.
// The token parameter is the token to verify — it is sent as the
// Bearer token for this specific request, overriding the client's
// own credentials.
func (c *Client) VerifyToken(ctx context.Context, token string) (*TokenInfo, error) {
	if token == "" {
		return nil, fmt.Errorf("webapi: token is required")
	}

	var info TokenInfo
	if err := c.doWithBearerToken(ctx, http.MethodGet, "/auth/v1/oauth/verify", token, nil, &info); err != nil {
		return nil, err
	}
	return &info, nil
}

// AccountInfo contains public account information.
type AccountInfo struct {
	AccountID         string `json:"accountId"`
	DisplayName       string `json:"displayName"`
	PreferredLanguage string `json:"preferredLanguage"`
}

// GetAccounts retrieves account information for the given account IDs.
// Maximum 100 account IDs per request.
func (c *Client) GetAccounts(ctx context.Context, accountIDs []string) ([]AccountInfo, error) {
	if len(accountIDs) == 0 {
		return nil, fmt.Errorf("webapi: accountIDs is required")
	}
	if len(accountIDs) > 100 {
		return nil, fmt.Errorf("webapi: maximum 100 accountIDs per request")
	}

	params := url.Values{}
	for _, id := range accountIDs {
		params.Add("accountId", id)
	}
	path := "/auth/v1/accounts?" + params.Encode()

	var accounts []AccountInfo
	if err := c.doGet(ctx, path, &accounts); err != nil {
		return nil, err
	}
	return accounts, nil
}

// escapePathSegment ensures a path segment is safe for URL inclusion.
func escapePathSegment(s string) string {
	return strings.ReplaceAll(url.PathEscape(s), "+", "%2B")
}
