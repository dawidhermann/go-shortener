package shortener

import (
	"fmt"
	"math/rand"
	"time"
)

const shortenStringLength int8 = 7

func ShortenUrl() string {
	return randomString(shortenStringLength)
}

func randomString(length int8) string {
	rand.Seed(time.Now().UnixNano())
	b := make([]byte, length)
	rand.Read(b)
	return fmt.Sprintf("%x", b)[:length]
}
