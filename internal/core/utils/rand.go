package utils

import (
	"crypto/rand"
)

const (
	charset = "abcdefghijklmnopqrstuvwxyz0123456789"
)

// randomString generates a random string of the given length.
func RandomString(length int) (string, error) {
	b := make([]byte, length)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	for i := 0; i < length; i++ {
		b[i] = charset[int(b[i])%len(charset)]
	}
	return string(b), nil
}
