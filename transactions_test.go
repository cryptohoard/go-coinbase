package coinbase_test

import (
	"testing"

	coinbase "github.com/cryptohoard/go-coinbase"
	"github.com/stretchr/testify/assert"
	"google.golang.org/api/iterator"
)

func TestTransactionsNext(t *testing.T) {
	assert := assert.New(t)
	c := coinbase.NewClient(secret, key)
	assert.NotNil(c)

	account, err := c.Account("ff6d5390-7054-509e-ac84-cf315860b4d2")
	assert.Empty(err)

	it := account.Transactions()
loop:
	for {
		transaction, err := it.Next()
		switch err {
		case nil:
			assert.NotEmpty(transaction.ID)
		case iterator.Done:
			break loop
		default:
			t.Fatal(err)
		}
	}
}

func TestTransactionsPager(t *testing.T) {
	assert := assert.New(t)
	c := coinbase.NewClient(secret, key)
	assert.NotNil(c)

	account, err := c.Account("ff6d5390-7054-509e-ac84-cf315860b4d2")
	assert.Empty(err)

	it := account.Transactions()
	p := iterator.NewPager(it, 25, "")
	for {
		var transactions []coinbase.Transaction
		nextPageToken, err := p.NextPage(&transactions)
		if err != nil {
			t.Fatal(err)
		}
		for _, t := range transactions {
			assert.NotEmpty(t.ID)
		}
		if nextPageToken == "" {
			break
		}
	}
}

func TestTransactionsNextWithPagination(t *testing.T) {
	assert := assert.New(t)
	c := coinbase.NewClient(secret, key)
	assert.NotNil(c)

	account, err := c.Account("ff6d5390-7054-509e-ac84-cf315860b4d2")
	assert.Empty(err)

	it := account.Transactions(
		coinbase.PaginationParams{
			Limit: 5,
		},
	)
loop:
	for {
		transaction, err := it.Next()
		switch err {
		case nil:
			assert.NotEmpty(transaction.ID)
		case iterator.Done:
			break loop
		default:
			t.Fatal(err)
		}
	}
}

func TestTransaction(t *testing.T) {
	assert := assert.New(t)
	c := coinbase.NewClient(secret, key)
	assert.NotNil(c)

	account, err := c.Account("ff6d5390-7054-509e-ac84-cf315860b4d2")
	assert.Empty(err)

	transaction, tErr := account.Transaction(
		"0344816f-052c-5d59-bdab-6e45ad765ed9",
	)
	assert.Empty(tErr)
	assert.NotEmpty(transaction.ID)
}
