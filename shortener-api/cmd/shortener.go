package main

import (
	"log"
	"net/http"
	"os"

	"github.com/dawidhermann/shortener-api/api"
	"github.com/dawidhermann/shortener-api/internal/database"
)

func main() {
	db, err := database.Connect(newDbConfig())
	if err != nil {
		log.Fatalf("DB ERROR")
	}
	log.Println("Starting app")
	apiMux := api.APIMux(db)
	api := http.Server{
		Addr:    ":8080",
		Handler: apiMux,
	}
	if err := api.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatal(err)
	}
}

func newDbConfig() database.DbConfig {
	return database.DbConfig{
		Host:     os.Getenv("POSTGRES_ADDR"),
		Name:     os.Getenv("POSTGRES_DB"),
		Password: os.Getenv("POSTGRES_PASSWORD"),
		User:     os.Getenv("POSTGRES_USER"),
		Schema:   "public",
	}
}
