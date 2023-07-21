package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/dawidhermann/shortener-api/api"
	"github.com/dawidhermann/shortener-api/appbase/config"
	"github.com/dawidhermann/shortener-api/internal/auth"
	"github.com/dawidhermann/shortener-api/internal/database"
)

func main() {
	appConfig := config.GetAppConfiguration()
	db, err := database.Connect(appConfig.DatabaseConfig)
	if err != nil {
		log.Fatalf(err.Error())
	}
	keyManager, err := auth.NewKeyReader(appConfig.AuthConfig.PrivateKeyPath, appConfig.AuthConfig.PublicKeyPath)
	if err != nil {
		log.Fatal(err.Error())
	}
	auth := auth.New(keyManager, appConfig.AuthConfig.JwtAuthTimeSec)
	apiMux := api.APIMux(api.AppConfig{
		Auth: auth,
		Db:   db,
	})
	api := http.Server{
		Addr:    fmt.Sprintf(":%d", appConfig.ApiPort),
		Handler: apiMux,
	}
	log.Printf("Listening on port %s", api.Addr)
	if err := api.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatal(err)
	}
}
