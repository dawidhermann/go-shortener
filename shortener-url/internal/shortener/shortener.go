// Generates shortened urls
package shortener

import (
	"math/rand"
	"strings"
	"time"
)

const (
	characterSet = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
)

func shortenUrl() string {
	return randomString()
}

func randomString() string {
	r := rand.New(rand.NewSource(time.Hour.Microseconds()))
	return base62Encode(r.Uint64())
}

func base62Encode(number uint64) string {
	length := len(characterSet)
	var encodedBuilder strings.Builder
	encodedBuilder.Grow(10)
	for ; number > 0; number = number / uint64(length) {
		encodedBuilder.WriteByte(characterSet[(number % uint64(length))])
	}

	return encodedBuilder.String()
}
