package coinbase_test

import (
	"testing"

	coinbase "github.com/cryptohoard/go-coinbase"
	"github.com/stretchr/testify/assert"
)

func TestUser(t *testing.T) {
	assert := assert.New(t)

	c := coinbase.NewAPIKeyClient(secret, key)
	assert.NotNil(c)

	user, err := c.User()
	assert.Nil(err)
	assert.NotEmpty(user.Email)
	assert.NotEmpty(user.Name)
}
