package shortener

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

const (
	characterSet = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
)

func ShortenUrl() string {
	return randomString()
}

func randomString() string {
	rand.Seed(time.Now().UnixNano())
	//b := make([]byte, length)
	//rand.Read(b)
	//return fmt.Sprintf("%x", b)[:length]
	return base62Encode(rand.Uint64())
}

func base62Encode(number uint64) string {
	fmt.Println(number)
	length := len(characterSet)
	var encodedBuilder strings.Builder
	encodedBuilder.Grow(10)
	for ; number > 0; number = number / uint64(length) {
		encodedBuilder.WriteByte(characterSet[(number % uint64(length))])
	}

	return encodedBuilder.String()
}
