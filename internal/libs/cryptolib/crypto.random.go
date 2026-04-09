package cryptolib

import (
	"crypto/rand"
	"encoding/base64"
)

func GenerateToken(byteSize int, encoder func([]byte) string) (string, error) {
	b := make([]byte, byteSize)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	return encoder(b), nil
}

func RandomString(n int) string {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		return ""
	}
	return base64.URLEncoding.EncodeToString(b)
}
