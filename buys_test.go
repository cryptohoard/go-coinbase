package coinbase_test

import (
	"testing"

	coinbase "github.com/cryptohoard/go-coinbase"
	"github.com/stretchr/testify/assert"
	"google.golang.org/api/iterator"
)

func TestBuysNext(t *testing.T) {
	assert := assert.New(t)
	c := coinbase.NewAPIKeyClient(secret, key)
	assert.NotNil(c)

	account, err := c.Account("ff6d5390-7054-509e-ac84-cf315860b4d2")
	assert.Empty(err)

	it := account.Buys()
loop:
	for {
		buy, err := it.Next()
		switch err {
		case nil:
			assert.NotEmpty(buy.ID)
		case iterator.Done:
			break loop
		default:
			t.Fatal(err)
		}
	}
}

func TestBuysPager(t *testing.T) {
	assert := assert.New(t)
	c := coinbase.NewAPIKeyClient(secret, key)
	assert.NotNil(c)

	account, err := c.Account("ff6d5390-7054-509e-ac84-cf315860b4d2")
	assert.Empty(err)

	it := account.Buys()
	p := iterator.NewPager(it, 25, "")
	for {
		var buys []coinbase.Buy
		nextPageToken, err := p.NextPage(&buys)
		if err != nil {
			t.Fatal(err)
		}
		for _, b := range buys {
			assert.NotEmpty(b.ID)
		}
		if nextPageToken == "" {
			break
		}
	}
}

func TestBuysNextWithPagination(t *testing.T) {
	assert := assert.New(t)
	c := coinbase.NewAPIKeyClient(secret, key)
	assert.NotNil(c)

	account, err := c.Account("ff6d5390-7054-509e-ac84-cf315860b4d2")
	assert.Empty(err)

	it := account.Buys(
		coinbase.PaginationParams{
			Limit: 5,
		},
	)
loop:
	for {
		buy, err := it.Next()
		switch err {
		case nil:
			assert.NotEmpty(buy.ID)
		case iterator.Done:
			break loop
		default:
			t.Fatal(err)
		}
	}
}

func TestBuy(t *testing.T) {
	assert := assert.New(t)
	c := coinbase.NewAPIKeyClient(secret, key)
	assert.NotNil(c)

	account, err := c.Account("ff6d5390-7054-509e-ac84-cf315860b4d2")
	assert.Empty(err)

	buy, tErr := account.Buy(
		"7e2f7d84-bad9-58ef-841b-7e07f0a26d26",
	)
	assert.Empty(tErr)
	assert.NotEmpty(buy.ID)
}
