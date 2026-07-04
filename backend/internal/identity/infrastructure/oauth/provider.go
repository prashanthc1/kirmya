// Package oauth implements OpenID-Connect-style OAuth providers (Google,
// LinkedIn) satisfying application.OAuthProvider. Construct only the providers
// whose client credentials are configured.
package oauth

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// Provider is a generic authorization-code + userinfo OIDC client.
type Provider struct {
	name          string
	clientID      string
	clientSecret  string
	redirectURL   string
	scope         string
	authEndpoint  string
	tokenEndpoint string
	userinfoURL   string
	http          *http.Client
}

func (p *Provider) Name() string { return p.name }

// AuthURL builds the provider authorization URL.
func (p *Provider) AuthURL(state string) string {
	q := url.Values{}
	q.Set("response_type", "code")
	q.Set("client_id", p.clientID)
	q.Set("redirect_uri", p.redirectURL)
	q.Set("scope", p.scope)
	q.Set("state", state)
	return p.authEndpoint + "?" + q.Encode()
}

// Exchange swaps the authorization code for tokens, then fetches the OIDC
// userinfo, returning (providerUID, email, fullName).
func (p *Provider) Exchange(ctx context.Context, code string) (string, string, string, error) {
	form := url.Values{}
	form.Set("grant_type", "authorization_code")
	form.Set("code", code)
	form.Set("redirect_uri", p.redirectURL)
	form.Set("client_id", p.clientID)
	form.Set("client_secret", p.clientSecret)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, p.tokenEndpoint, strings.NewReader(form.Encode()))
	if err != nil {
		return "", "", "", err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")

	resp, err := p.http.Do(req)
	if err != nil {
		return "", "", "", err
	}
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode != http.StatusOK {
		return "", "", "", errors.New("oauth token exchange failed: " + resp.Status)
	}
	var tok struct {
		AccessToken string `json:"access_token"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&tok); err != nil {
		return "", "", "", err
	}
	if tok.AccessToken == "" {
		return "", "", "", errors.New("oauth token exchange returned no access token")
	}
	return p.userinfo(ctx, tok.AccessToken)
}

func (p *Provider) userinfo(ctx context.Context, accessToken string) (string, string, string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, p.userinfoURL, nil)
	if err != nil {
		return "", "", "", err
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)
	resp, err := p.http.Do(req)
	if err != nil {
		return "", "", "", err
	}
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode != http.StatusOK {
		return "", "", "", errors.New("oauth userinfo failed: " + resp.Status)
	}
	var info struct {
		Sub   string `json:"sub"`
		Email string `json:"email"`
		Name  string `json:"name"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&info); err != nil {
		return "", "", "", err
	}
	if info.Sub == "" || info.Email == "" {
		return "", "", "", errors.New("oauth userinfo missing sub/email")
	}
	return info.Sub, info.Email, info.Name, nil
}

func newClient() *http.Client { return &http.Client{Timeout: 10 * time.Second} }

// NewGoogle builds a Google OIDC provider.
func NewGoogle(clientID, clientSecret, redirectURL string) *Provider {
	return &Provider{
		name:          "google",
		clientID:      clientID,
		clientSecret:  clientSecret,
		redirectURL:   redirectURL,
		scope:         "openid email profile",
		authEndpoint:  "https://accounts.google.com/o/oauth2/v2/auth",
		tokenEndpoint: "https://oauth2.googleapis.com/token",
		userinfoURL:   "https://openidconnect.googleapis.com/v1/userinfo",
		http:          newClient(),
	}
}

// NewLinkedIn builds a LinkedIn OIDC provider.
func NewLinkedIn(clientID, clientSecret, redirectURL string) *Provider {
	return &Provider{
		name:          "linkedin",
		clientID:      clientID,
		clientSecret:  clientSecret,
		redirectURL:   redirectURL,
		scope:         "openid email profile",
		authEndpoint:  "https://www.linkedin.com/oauth/v2/authorization",
		tokenEndpoint: "https://www.linkedin.com/oauth/v2/accessToken",
		userinfoURL:   "https://api.linkedin.com/v2/userinfo",
		http:          newClient(),
	}
}
