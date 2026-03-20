package cryptolib

import "crypto/rand"

func GenerateToken(byteSize int, encoder func([]byte) string) (string, error) {
	b := make([]byte, byteSize)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	return encoder(b), nil
}
