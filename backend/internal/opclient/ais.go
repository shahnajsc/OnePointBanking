package opclient

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type AISClient struct {
	HTTP            *http.Client
	MTLSBase        string
	ClientID        string
	ClientSecret    string
	APIKey          string
	FAPIFinancialID string
}

type tokenResp struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	Scope       string `json:"scope"`
	ExpiresIn   any    `json:"expires_in"`
	Status      string `json:"status"`
}

func (c *AISClient) ClientCredentialsToken(ctx context.Context) (string, error) {
	form := url.Values{}
	form.Set("grant_type", "client_credentials")
	form.Set("scope", "accounts") // OP requirement: one scope per request
	form.Set("client_id", c.ClientID)
	form.Set("client_secret", c.ClientSecret)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.MTLSBase+"/oauth/token", strings.NewReader(form.Encode()))
	if err != nil {
		return "", fmt.Errorf("create token request: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := c.HTTP.Do(req)
	if err != nil {
		return "", fmt.Errorf("token request failed: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode/100 != 2 {
		return "", fmt.Errorf("token non-2xx: %s body=%s", resp.Status, string(body))
	}

	var tr tokenResp
	if err := json.Unmarshal(body, &tr); err != nil {
		return "", fmt.Errorf("parse token response: %w body=%s", err, string(body))
	}
	if tr.AccessToken == "" {
		return "", fmt.Errorf("missing access_token in response body=%s", string(body))
	}
	return tr.AccessToken, nil
}

type createAuthReq struct {
	Expires string `json:"expires"`
}

type createAuthResp struct {
	AuthorizationID string `json:"authorizationId"`
	Status          string `json:"status"`
	Expires         string `json:"expires"`
}

func (c *AISClient) CreateAuthorization(ctx context.Context, bearerToken string) (string, error) {
	// Set expires (now+1h) for sandbox
	expires := time.Now().Add(1 * time.Hour).Format(time.RFC3339)
	payload := fmt.Sprintf(`{"expires":"%s"}`, expires)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.MTLSBase+"/accounts-psd2/v1/authorizations", strings.NewReader(payload))
	if err != nil {
		return "", fmt.Errorf("create authorizations request: %w", err)
	}

	req.Header.Set("x-api-key", c.APIKey)
	req.Header.Set("Authorization", "Bearer "+bearerToken)
	req.Header.Set("x-fapi-financial-id", c.FAPIFinancialID)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.HTTP.Do(req)
	if err != nil {
		return "", fmt.Errorf("authorizations request failed: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode/100 != 2 {
		return "", fmt.Errorf("authorizations non-2xx: %s body=%s", resp.Status, string(body))
	}

	var ar createAuthResp
	if err := json.Unmarshal(body, &ar); err != nil {
		return "", fmt.Errorf("parse authorizations response: %w body=%s", err, string(body))
	}
	if ar.AuthorizationID == "" {
		return "", fmt.Errorf("missing authorizationId body=%s", string(body))
	}
	return ar.AuthorizationID, nil
}
