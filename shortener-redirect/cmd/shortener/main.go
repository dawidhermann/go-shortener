package main

import (
	"fmt"
	"net/http"

	"github.com/dawidhermann/shortener-redirect/api/redirectionctrl"
	"github.com/dawidhermann/shortener-redirect/internal/config"
	"github.com/dawidhermann/shortener-redirect/internal/db"
)

func main() {
	cfg := config.GetAppConfiguration()
	fmt.Println(cfg)
	store := db.New(cfg.Store)
	redirectionctrl := redirectionctrl.New(store)
	mux := http.NewServeMux()
	mux.HandleFunc("/", redirectionctrl.Redirect)
	port := fmt.Sprintf(":%s", cfg.Api.Port)
	http.ListenAndServe(port, mux)
}
