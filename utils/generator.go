package utils

import "math/rand"

// charset is the set of characters allowed in a short code.
const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"

// GenerateShortCode returns a random 6-character alphanumeric string,
// e.g. "Ab12Cd". It does not check for uniqueness itself — the caller
// (handlers.ShortenURL) is responsible for retrying on collision.
func GenerateShortCode() string {
	code := make([]byte, 6)
	for i := range code {
		code[i] = charset[rand.Intn(len(charset))]
	}
	return string(code)
}
