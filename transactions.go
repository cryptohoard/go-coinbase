package coinbase

import (
	"fmt"
	"net/http"
	"strings"

	"google.golang.org/api/iterator"
)

const (
	transactionsURL = "/transactions"
)

type SortOrder int

const (
	OrderNone SortOrder = iota
	OrderAscending
	OrderDescending
)

type PaginationParams struct {
	Limit int
	Order SortOrder
}

func (a Account) Transactions(p ...PaginationParams) *TransactionIterator {
	it := &TransactionIterator{
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

type TransactionIterator struct {
	pageInfo  *iterator.PageInfo
	nextFunc  func() error
	c         *Client
	params    []PaginationParams
	accountID string
	firstCall bool
	r         transactionsResp
	items     []Transaction
}

func (it *TransactionIterator) PageInfo() *iterator.PageInfo {
	return it.pageInfo
}

func (it *TransactionIterator) fetch(
	pageSize int,
	pageToken string,
) (string, error) {
	url := accountsURL + "/" + it.accountID + transactionsURL

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

	r := transactionsResp{}

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

func (it *TransactionIterator) Next() (Transaction, error) {
	if err := it.nextFunc(); err != nil {
		return Transaction{}, err
	}
	item := it.items[0]
	it.items = it.items[1:]
	return item, nil
}

type Transaction struct {
	ID     string `json:"id"`
	Type   string `json:"type"`
	Status string `json:"status"`
	Amount struct {
		Amount   string `json:"amount"`
		Currency string `json:"currency"`
	} `json:"amount"`
	NativeAmount struct {
		Amount   string `json:"amount"`
		Currency string `json:"currency"`
	} `json:"native_amount"`
	Description  string `json:"description"`
	CreatedAt    string `json:"created_at"`
	UpdatedAt    string `json:"updated_at"`
	Resource     string `json:"resource"`
	ResourcePath string `json:"resource_path"`
	Buy          struct {
		ID           string `json:"id"`
		Resource     string `json:"resource"`
		ResourcePath string `json:"resource_path"`
	} `json:"buy,omitempty"`
	Details struct {
		Title    string `json:"title"`
		Subtitle string `json:"subtitle"`
	} `json:"details"`
	To struct {
		Resource string `json:"resource"`
		Email    string `json:"email"`
	} `json:"to,omitempty"`
	InstantExchange bool `json:"instant_exchange,omitempty"`
	Sell            struct {
		ID           string `json:"id"`
		Resource     string `json:"resource"`
		ResourcePath string `json:"resource_path"`
	} `json:"sell,omitempty"`
	Network struct {
		Status string `json:"status"`
		Name   string `json:"name"`
	} `json:"network,omitempty"`
}

type pagination struct {
	EndingBefore  string `json:"ending_before"`
	StartingAfter string `json:"starting_after"`
	Limit         int    `json:"limit"`
	Order         string `json:"order"`
	PreviousURI   string `json:"previous_uri"`
	NextURI       string `json:"next_uri"`
}

type transactionsResp struct {
	Data       []Transaction `json:"data"`
	Pagination pagination    `json:"pagination"`
}

func (a Account) Transaction(transactionID string) (Transaction, error) {
	url := accountsURL + "/" + a.ID + transactionsURL + "/" + transactionID

	r := transactionResp{}

	_, err := a.c.Request(
		http.MethodGet,
		url,
		nil,
		&r,
	)
	if err != nil {
		return Transaction{}, err
	}

	return r.Data, nil
}

type transactionResp struct {
	Data Transaction `json:"data"`
}
