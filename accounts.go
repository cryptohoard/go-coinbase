package coinbase

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"google.golang.org/api/iterator"
)

const (
	accountsURL = "/accounts"
)

func (c *Client) Accounts(p ...PaginationParams) *AccountIterator {
	it := &AccountIterator{
		c:         c,
		params:    p,
		firstCall: true,
	}
	it.pageInfo, it.nextFunc = iterator.NewPageInfo(
		it.fetch,
		func() int { return len(it.items) },
		func() interface{} { b := it.items; it.items = nil; return b },
	)
	return it
}

type AccountIterator struct {
	pageInfo  *iterator.PageInfo
	nextFunc  func() error
	c         *Client
	params    []PaginationParams
	firstCall bool
	r         accountsResp
	items     []Account
}

func (it *AccountIterator) PageInfo() *iterator.PageInfo {
	return it.pageInfo
}

func (it *AccountIterator) fetch(
	pageSize int,
	pageToken string,
) (string, error) {
	url := accountsURL

	if !it.firstCall {
		if it.r.Pagination.NextURI == "" {
			return "", iterator.Done
		}

		url = strings.Replace(it.r.Pagination.NextURI, "/v2", "", -1)
	} else if len(it.params) > 0 {
		paramsStr := ""

		if it.params[0].Limit > 0 {
			paramsStr += fmt.Sprintf("&limit=%d", it.params[0].Limit)
		}

		if it.params[0].Order == OrderAscending {
			paramsStr += "&order=asc"
		}

		if paramsStr != "" {
			url += "?" + paramsStr
		}
	}

	it.firstCall = false

	r := accountsResp{}

	_, err := it.c.Request(
		http.MethodGet,
		url,
		nil,
		&r,
	)
	if err != nil {
		return "", err
	}

	it.r = r
	items := r.Data[:]
	it.items = append(it.items, items...)

	return it.r.Pagination.NextURI, nil
}

func (it *AccountIterator) Next() (Account, error) {
	if err := it.nextFunc(); err != nil {
		return Account{}, err
	}
	item := it.items[0]
	it.items = it.items[1:]
	item.c = it.c
	return item, nil
}

type Account struct {
	c       *Client
	Balance struct {
		Amount   string `json:"amount"`
		Currency string `json:"currency"`
	} `json:"balance"`
	CreatedAt time.Time `json:"created_at"`
	Currency  struct {
		AddressRegex string `json:"address_regex"`
		Code         string `json:"code"`
		Color        string `json:"color"`
		Exponent     int    `json:"exponent"`
		Name         string `json:"name"`
		Type         string `json:"type"`
	} `json:"currency"`
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	Primary      bool      `json:"primary"`
	Resource     string    `json:"resource"`
	ResourcePath string    `json:"resource_path"`
	Type         string    `json:"type"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type accountsResp struct {
	Data       []Account  `json:"data"`
	Pagination pagination `json:"pagination"`
}

type accountResp struct {
	Data Account `json:"data"`
}

func (c *Client) Account(ID string) (Account, error) {
	r := accountResp{}
	_, err := c.Request(http.MethodGet, accountsURL+"/"+ID, nil, &r)
	if err != nil {
		return Account{}, err
	}
	r.Data.c = c
	return r.Data, nil
}
