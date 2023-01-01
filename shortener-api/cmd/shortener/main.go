package main

import (
	"github.com/dawidhermann/shortener-api/api"
	"github.com/dawidhermann/shortener-api/config"
	"github.com/dawidhermann/shortener-api/internal/auth"
	"github.com/dawidhermann/shortener-api/internal/db"
	"github.com/dawidhermann/shortener-api/internal/rpc"
	"log"
)

func main() {
	appConfig := config.NewAppConfig()
	sqlConnection := db.Connect(appConfig.DbAppConfig)
	authManager, err := auth.NewAuthenticationManager(appConfig.AuthAppConfig)
	if err != nil {
		log.Fatalf("cannot initialize auth manager. Reason: %s", err)
	}
	connRpc := rpc.Connect(appConfig.RpcAppConfig)
	api.StartServer(appConfig.ServerAppConfig, connRpc, sqlConnection, authManager)
}
