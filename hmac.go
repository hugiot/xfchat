package xfchat

import (
	"crypto/hmac"
	"crypto/sha256"
)

func HmacSHA256(key string, data string) []byte {
	mac := hmac.New(sha256.New, []byte(key))
	_, _ = mac.Write([]byte(data))
	return mac.Sum(nil)
}
