package util

import (
	"crypto/rand"
)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
const defaultLength = 6

func GenerateShortCode() string {
	b := make([]byte, defaultLength)
	rand.Read(b)

	for i := range b {
		b[i] = charset[int(b[i])%len(charset)]
	}
	return string(b)
}
