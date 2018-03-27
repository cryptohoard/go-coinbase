package coinbase_test

import (
	"testing"

	coinbase "github.com/cryptohoard/go-coinbase"
	"github.com/stretchr/testify/assert"
	"google.golang.org/api/iterator"
)

func TestAccountsNext(t *testing.T) {
	assert := assert.New(t)

	c := coinbase.NewAPIKeyClient(secret, key)
	assert.NotNil(c)

	it := c.Accounts()
loop:
	for {
		account, err := it.Next()
		switch err {
		case nil:
			assert.NotEmpty(account.ID)
		case iterator.Done:
			break loop
		default:
			t.Fatal(err)
		}
	}
}

func TestAccountsPager(t *testing.T) {
	assert := assert.New(t)
	c := coinbase.NewAPIKeyClient(secret, key)
	assert.NotNil(c)

	it := c.Accounts()
	p := iterator.NewPager(it, 25, "")
	for {
		var accounts []coinbase.Account
		nextPageToken, err := p.NextPage(&accounts)
		if err != nil {
			t.Fatal(err)
		}
		for _, a := range accounts {
			assert.NotEmpty(a.ID)
		}
		if nextPageToken == "" {
			break
		}
	}
}

func TestAccount(t *testing.T) {
	assert := assert.New(t)

	c := coinbase.NewAPIKeyClient(secret, key)
	assert.NotNil(c)

	account, err := c.Account("d27743b6-fb0a-593e-82ab-fec85f0746e2")
	assert.Nil(err)
	assert.NotEmpty(account)
	assert.NotEmpty(account.ID)
}
