package webapi

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"golang.org/x/oauth2"
)

// WithClientCredentials configures OAuth2 client_credentials authentication.
// This is the standard flow for backend services.
func WithClientCredentials(clientID, clientSecret string) Option {
	return func(c *Client) {
		c.authBuilder = func(c *Client) {
			src := &clientCredentialsTokenSource{
				clientID:     clientID,
				clientSecret: clientSecret,
				deploymentID: c.deploymentID,
				tokenURL:     c.baseURL + "/auth/v1/oauth/token",
				httpClient:   c.httpClient,
			}
			c.tokenSource = oauth2.ReuseTokenSource(nil, src)
		}
	}
}

// WithExchangeCode configures OAuth2 exchange_code authentication.
// Used by game clients launched via the Epic Games Launcher.
func WithExchangeCode(clientID, clientSecret, code string) Option {
	return func(c *Client) {
		c.authBuilder = func(c *Client) {
			src := &exchangeCodeTokenSource{
				clientID:     clientID,
				clientSecret: clientSecret,
				code:         code,
				deploymentID: c.deploymentID,
				tokenURL:     c.baseURL + "/auth/v1/oauth/token",
				httpClient:   c.httpClient,
			}
			c.tokenSource = oauth2.ReuseTokenSource(nil, src)
		}
	}
}

type clientCredentialsTokenSource struct {
	clientID     string
	clientSecret string
	deploymentID string
	tokenURL     string
	httpClient   *http.Client
}

func (s *clientCredentialsTokenSource) Token() (*oauth2.Token, error) {
	data := url.Values{
		"grant_type":    {"client_credentials"},
		"deployment_id": {s.deploymentID},
	}

	req, err := http.NewRequest("POST", s.tokenURL, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, fmt.Errorf("webapi: client credentials request: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.SetBasicAuth(s.clientID, s.clientSecret)

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("webapi: client credentials: %w", err)
	}
	defer resp.Body.Close() //nolint:errcheck

	if resp.StatusCode != http.StatusOK {
		return nil, parseErrorResponse(resp)
	}

	return decodeTokenResponse(resp)
}

type exchangeCodeTokenSource struct {
	clientID     string
	clientSecret string
	code         string
	deploymentID string
	tokenURL     string
	httpClient   *http.Client
}

func (s *exchangeCodeTokenSource) Token() (*oauth2.Token, error) {
	data := url.Values{
		"grant_type":    {"exchange_code"},
		"exchange_code": {s.code},
		"deployment_id": {s.deploymentID},
	}

	req, err := http.NewRequest("POST", s.tokenURL, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, fmt.Errorf("webapi: exchange code request: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.SetBasicAuth(s.clientID, s.clientSecret)

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("webapi: exchange code: %w", err)
	}
	defer resp.Body.Close() //nolint:errcheck

	if resp.StatusCode != http.StatusOK {
		return nil, parseErrorResponse(resp)
	}

	return decodeTokenResponse(resp)
}

func decodeTokenResponse(resp *http.Response) (*oauth2.Token, error) {
	var body struct {
		AccessToken  string `json:"access_token"`
		TokenType    string `json:"token_type"`
		ExpiresIn    int    `json:"expires_in"`
		RefreshToken string `json:"refresh_token"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return nil, fmt.Errorf("webapi: decode token response: %w", err)
	}

	return &oauth2.Token{
		AccessToken:  body.AccessToken,
		TokenType:    body.TokenType,
		RefreshToken: body.RefreshToken,
		Expiry:       time.Now().Add(time.Duration(body.ExpiresIn) * time.Second),
	}, nil
}
