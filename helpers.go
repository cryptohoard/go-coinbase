package coinbase

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
)

func generateSig(message, secret string) (string, error) {
	signature := hmac.New(sha256.New, []byte(secret))
	_, err := signature.Write([]byte(message))
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(signature.Sum(nil)), nil
}
