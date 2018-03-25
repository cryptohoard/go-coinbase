package coinbase

import (
	"fmt"
	"net/http"
	"strings"

	"google.golang.org/api/iterator"
)

const (
	buysURL = "/buys"
)

func (a Account) Buys(p ...PaginationParams) *BuyIterator {
	it := &BuyIterator{
		c:         a.c,
		params:    p,
		accountID: a.ID,
		firstCall: true,
	}
	it.pageInfo, it.nextFunc = iterator.NewPageInfo(
		it.fetch,
		func() int { return len(it.items) },
		func() interface{} { b := it.items; it.items = nil; return b },
	)
	return it
}

type BuyIterator struct {
	pageInfo  *iterator.PageInfo
	nextFunc  func() error
	c         *Client
	params    []PaginationParams
	accountID string
	firstCall bool
	r         buysResp
	items     []Buy
}

func (it *BuyIterator) PageInfo() *iterator.PageInfo {
	return it.pageInfo
}

func (it *BuyIterator) fetch(
	pageSize int,
	pageToken string,
) (string, error) {
	url := accountsURL + "/" + it.accountID + buysURL

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

	r := buysResp{}

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

func (it *BuyIterator) Next() (Buy, error) {
	if err := it.nextFunc(); err != nil {
		return Buy{}, err
	}
	item := it.items[0]
	it.items = it.items[1:]
	return item, nil
}

type Buy struct {
	ID            string `json:"id"`
	Status        string `json:"status"`
	PaymentMethod struct {
		ID           string `json:"id"`
		Resource     string `json:"resource"`
		ResourcePath string `json:"resource_path"`
	} `json:"payment_method"`
	Transaction struct {
		ID           string `json:"id"`
		Resource     string `json:"resource"`
		ResourcePath string `json:"resource_path"`
	} `json:"transaction"`
	Amount struct {
		Amount   string `json:"amount"`
		Currency string `json:"currency"`
	} `json:"amount"`
	Total struct {
		Amount   string `json:"amount"`
		Currency string `json:"currency"`
	} `json:"total"`
	Subtotal struct {
		Amount   string `json:"amount"`
		Currency string `json:"currency"`
	} `json:"subtotal"`
	CreatedAt    string `json:"created_at"`
	UpdatedAt    string `json:"updated_at"`
	Resource     string `json:"resource"`
	ResourcePath string `json:"resource_path"`
	Committed    bool   `json:"committed"`
	Instant      bool   `json:"instant"`
	Fee          struct {
		Amount   string `json:"amount"`
		Currency string `json:"currency"`
	} `json:"fee"`
	PayoutAt string `json:"payout_at"`
}

type buysResp struct {
	Data       []Buy      `json:"data"`
	Pagination pagination `json:"pagination"`
}

func (a Account) Buy(buyID string) (Buy, error) {
	url := accountsURL + "/" + a.ID + buysURL + "/" + buyID

	r := buyResp{}

	_, err := a.c.Request(
		http.MethodGet,
		url,
		nil,
		&r,
	)
	if err != nil {
		return Buy{}, err
	}

	return r.Data, nil
}

type buyResp struct {
	Data Buy `json:"data"`
}
