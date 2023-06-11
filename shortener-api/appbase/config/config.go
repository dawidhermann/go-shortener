package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type DbConfig struct {
	Port     int
	Host     string
	User     string
	Password string
	DbName   string
}

type GrpcConfig struct {
	GrpcServerPort int
	GrpcServerHost string
}

type AuthConfig struct {
	JwtSecretKey   string
	JwtAuthTimeSec int
}

type AppConfig struct {
	ApiPort          int
	DatabaseConfig   DbConfig
	GrpcServerConfig GrpcConfig
}

func init() {
	// loadEnvFile()
	viper.AutomaticEnv()
	viper.SetEnvPrefix("SHORTENER_API")
	setDefaults()
}

func GetAppConfiguration() AppConfig {
	dbConfig := DbConfig{
		Port:     viper.GetInt("DB_PORT"),
		Host:     viper.GetString("DB_HOST"),
		User:     viper.GetString("DB_USER"),
		Password: viper.GetString("DB_PASSWORD"),
		DbName:   viper.GetString("DB_NAME"),
	}
	grpcConfig := GrpcConfig{
		GrpcServerPort: viper.GetInt("GRPC_SERVER_PORT"),
		GrpcServerHost: viper.GetString("GRPC_SERVER_HOST"),
	}
	return AppConfig{
		ApiPort:          viper.GetInt("API_PORT"),
		DatabaseConfig:   dbConfig,
		GrpcServerConfig: grpcConfig,
	}
}

func loadEnvFile() {
	viper.SetConfigName(".env")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}
}

func setDefaults() {
	viper.SetDefault("API_PORT", 8090)
}
