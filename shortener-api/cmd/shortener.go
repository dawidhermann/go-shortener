// Web app entry point
package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/dawidhermann/shortener-api/api"
	"github.com/dawidhermann/shortener-api/internal/auth"
	"github.com/dawidhermann/shortener-api/internal/config"
	"github.com/dawidhermann/shortener-api/internal/database"
	"github.com/dawidhermann/shortener-api/internal/rpc"
)

func main() {
	appConfig := config.GetAppConfiguration()
	db, err := database.Connect(appConfig.DatabaseConfig)
	if err != nil {
		log.Fatalf(err.Error())
	}
	auth := auth.New(appConfig.AuthConfig.SecretKey, appConfig.AuthConfig.JwtAuthTimeSec)
	rpcConn := rpc.Connect(appConfig.GrpcServerConfig)
	apiMux := api.APIMux(api.AppConfig{
		Auth:    auth,
		Db:      db,
		RpcConn: rpcConn,
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
