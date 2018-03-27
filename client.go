package coinbase

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

const (
	coinbaseAPIVer     = "/v2"
	coinbaseAPIVerDate = "2017-12-16"
)

type client struct {
	BaseURL      string
	secret       string
	key          string
	httpClient   *http.Client
	oauth        bool
	clientID     string
	clientSecret string
	redirectURL  string
	oauthToken   string
}

func (c *client) Request(
	method string,
	url string,
	params interface{},
	result interface{},
) (res *http.Response, err error) {
	var data []byte
	body := bytes.NewReader(make([]byte, 0))

	if params != nil {
		data, err = json.Marshal(params)
		if err != nil {
			return res, err
		}

		body = bytes.NewReader(data)
	}

	fullURL := fmt.Sprintf("%s%s%s", c.BaseURL, coinbaseAPIVer, url)
	req, err := http.NewRequest(method, fullURL, body)
	if err != nil {
		return res, err
	}

	timestamp := strconv.FormatInt(time.Now().Unix(), 10)

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("User-Agent", "griphook 1.0")

	var h map[string]string
	if c.oauth {
		h, err = c.oauthHeaders()
	} else {
		h, err = c.apiKeyHeaders(
			method,
			coinbaseAPIVer+url,
			timestamp,
			string(data),
		)
	}
	for k, v := range h {
		req.Header.Add(k, v)
	}

	res, err = c.httpClient.Do(req)
	if err != nil {
		return res, err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		defer res.Body.Close()
		reqErr := Error{}
		decoder := json.NewDecoder(res.Body)
		if err := decoder.Decode(&reqErr); err != nil {
			return res, err
		}

		return res, error(reqErr)
	}

	if result != nil {
		decoder := json.NewDecoder(res.Body)
		if err = decoder.Decode(&result); err != nil {
			return res, err
		}
	}

	return res, nil
}
