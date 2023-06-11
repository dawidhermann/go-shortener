package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/dawidhermann/shortener-api/api"
	"github.com/dawidhermann/shortener-api/appbase/config"
	"github.com/dawidhermann/shortener-api/internal/database"
)

func main() {
	appConfig := config.GetAppConfiguration()
	db, err := database.Connect(appConfig.DatabaseConfig)
	if err != nil {
		log.Fatalf(err.Error())
	}
	apiMux := api.APIMux(db)
	api := http.Server{
		Addr:    fmt.Sprintf(":%d", appConfig.ApiPort),
		Handler: apiMux,
	}
	log.Printf("Listening on port %s", api.Addr)
	if err := api.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatal(err)
	}
}
