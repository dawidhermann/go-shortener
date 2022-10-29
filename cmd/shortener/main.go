package main

import (
	"github.com/dawidhermann/go-shortener/internal/api"
)

func main() {
	//url := shortener.ShortenUrl()
	//fmt.Print(url)
	api.StartServer()
}
