package coinbase

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type OAuthClient struct {
	*client
}

func NewOAuthClient(clientID, clientSecret, redirectURL string) *OAuthClient {
	return &OAuthClient{
		&client{
			BaseURL:      "https://api.coinbase.com",
			clientID:     clientID,
			clientSecret: clientSecret,
			redirectURL:  redirectURL,
			oauth:        true,
			httpClient: &http.Client{
				Timeout: 15 * time.Second,
			},
		},
	}
}

func (o *OAuthClient) CreateAuthorizeUrl(scope []string, state string) string {
	path := "https://www.coinbase.com/oauth/authorize"

	v := url.Values{}

	v.Set("response_type", "code")

	if o.clientID != "" {
		v.Add("client_id", o.clientID)
	}

	if o.redirectURL != "" {
		v.Add("redirect_uri", o.redirectURL)
	}

	if state != "" {
		v.Add("state", state)
	}

	v.Add("account", "all")

	if len(scope) > 0 {
		scopeStr := strings.Join(scope, ",")
		v.Add("scope", scopeStr)
	}

	return path + "?" + v.Encode()
}

func (o *OAuthClient) SetToken(token string) {
	o.oauthToken = token
}

type OAuthToken struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	Scope        string `json:"scope"`
}

func (o *OAuthClient) Tokens(code string) (*OAuthToken, error) {
	fullURL := fmt.Sprintf("%s/oauth/token", o.BaseURL)
	req, err := http.NewRequest(http.MethodPost, fullURL, nil)
	if err != nil {
		return nil, err
	}

	q := req.URL.Query()
	q.Add("grant_type", "authorization_code")
	q.Add("code", code)
	q.Add("client_id", o.clientID)
	q.Add("client_secret", o.clientSecret)
	q.Add("redirect_uri", o.redirectURL)
	req.URL.RawQuery = q.Encode()

	res, err := o.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		defer res.Body.Close()
		reqErr := Error{}
		decoder := json.NewDecoder(res.Body)
		if err := decoder.Decode(&reqErr); err != nil {
			return nil, err
		}

		return nil, error(reqErr)
	}

	tok := &OAuthToken{}
	decoder := json.NewDecoder(res.Body)
	if err := decoder.Decode(tok); err != nil {
		return nil, err
	}

	return tok, nil
}

func (o *OAuthClient) RefreshTokens(refreshToken string) (*OAuthToken, error) {
	fullURL := fmt.Sprintf("%s/oauth/token", o.BaseURL)
	req, err := http.NewRequest(http.MethodPost, fullURL, nil)
	if err != nil {
		return nil, err
	}

	q := req.URL.Query()
	q.Add("grant_type", "refresh_token")
	q.Add("client_id", o.clientID)
	q.Add("client_secret", o.clientSecret)
	q.Add("refresh_token", refreshToken)
	req.URL.RawQuery = q.Encode()

	res, err := o.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		defer res.Body.Close()
		reqErr := Error{}
		decoder := json.NewDecoder(res.Body)
		if err := decoder.Decode(&reqErr); err != nil {
			return nil, err
		}

		return nil, error(reqErr)
	}

	tok := &OAuthToken{}
	decoder := json.NewDecoder(res.Body)
	if err := decoder.Decode(tok); err != nil {
		return nil, err
	}

	return tok, nil
}

func (c *client) oauthHeaders() (map[string]string, error) {
	return map[string]string{
		"Authorization": "Bearer " + c.oauthToken,
	}, nil
}
