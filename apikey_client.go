package coinbase

import (
	"fmt"
	"net/http"
	"time"
)

type APIKeyClient struct {
	*client
}

func NewAPIKeyClient(secret, key string) *APIKeyClient {
	return &APIKeyClient{
		&client{
			BaseURL: "https://api.coinbase.com",
			secret:  secret,
			key:     key,
			httpClient: &http.Client{
				Timeout: 15 * time.Second,
			},
		},
	}
}

// Headers generates a map that can be used as headers to authenticate a request
func (c *client) apiKeyHeaders(
	method string,
	url string,
	timestamp string,
	data string,
) (map[string]string, error) {
	h := make(map[string]string)
	h["CB-VERSION"] = coinbaseAPIVerDate
	h["CB-ACCESS-KEY"] = c.key
	h["CB-ACCESS-TIMESTAMP"] = timestamp

	message := fmt.Sprintf(
		"%s%s%s%s",
		timestamp,
		method,
		url,
		data,
	)

	sig, err := generateSig(message, c.secret)
	if err != nil {
		return nil, err
	}
	h["CB-ACCESS-SIGN"] = sig

	return h, nil
}
