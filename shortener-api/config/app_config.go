package config

import (
	"log"
	"os"
	"strconv"
)

type DbConfig struct {
	DbUser     string
	DbPassword string
	DbAddr     string
	DbName     string
}

type RpcConfig struct {
	ServerHost string
	ServerPort string
}

type ServerConfig struct {
	ServerPort string
}

type AuthConfig struct {
	JwtExpTime   int
	JwtSecretKey string
}

type AppConfig struct {
	DbAppConfig     DbConfig
	RpcAppConfig    RpcConfig
	ServerAppConfig ServerConfig
	AuthAppConfig   AuthConfig
}

func newDbConfig() DbConfig {
	return DbConfig{
		DbAddr:     os.Getenv("POSTGRES_ADDR"),
		DbName:     os.Getenv("POSTGRES_DB"),
		DbPassword: os.Getenv("POSTGRES_PASSWORD"),
		DbUser:     os.Getenv("POSTGRES_USER"),
	}
}

func newRpcConfig() RpcConfig {
	return RpcConfig{
		ServerHost: os.Getenv("GRPC_SERVER_HOST"),
		ServerPort: os.Getenv("GRPC_SERVER_PORT"),
	}
}

func newServerConfig() ServerConfig {
	return ServerConfig{
		ServerPort: os.Getenv("SHORTENER_API_PORT"),
	}
}

func newAuthConfig() AuthConfig {
	tokenExpirationEnv := os.Getenv("JWT_AUTH_TIME_SEC")
	tokenExp, err := strconv.Atoi(tokenExpirationEnv)
	if err != nil {
		log.Fatalf("Cannot get token expiration time: %s", err)
	}
	return AuthConfig{
		JwtExpTime:   tokenExp,
		JwtSecretKey: os.Getenv("JWT_SECRET_KEY"),
	}
}

func NewAppConfig() AppConfig {
	return AppConfig{
		DbAppConfig:     newDbConfig(),
		AuthAppConfig:   newAuthConfig(),
		RpcAppConfig:    newRpcConfig(),
		ServerAppConfig: newServerConfig(),
	}
}
