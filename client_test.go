package coinbase_test

import "os"

var (
	key    = ""
	secret = ""
)

func init() {
	key = os.Getenv("GO_COINBASE_TEST_KEY")
	secret = os.Getenv("GO_COINBASE_TEST_SECRET")
}
