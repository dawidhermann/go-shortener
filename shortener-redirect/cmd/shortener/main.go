package main

import (
	"net/http"

	"github.com/dawidhermann/shortener-redirect/api/redirectionctrl"
	"github.com/dawidhermann/shortener-redirect/internal/config"
	"github.com/dawidhermann/shortener-redirect/internal/db"
)

func main() {
	cfg := config.StoreConfig{
		Address:  "redis://localhost:6379",
		Password: "s3CreTP4sS",
	}
	store := db.New(cfg)
	redirectionctrl := redirectionctrl.New(store)
	mux := http.NewServeMux()
	mux.HandleFunc("/", redirectionctrl.Redirect)
	http.ListenAndServe(":8092", mux)
}
