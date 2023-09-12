// App configuration definition
package config

import "github.com/spf13/viper"

type StoreConfig struct {
	Address  string
	Password string
}

type GrpcConfig struct {
	Port string
}

type AppConfig struct {
	Store StoreConfig
	Grpc  GrpcConfig
}

func init() {
	viper.AutomaticEnv()
	viper.SetEnvPrefix("SHORTENER_URL")
	setDefaults()
}

// Create new application's config
func GetAppConfiguration() AppConfig {
	grpcConfig := GrpcConfig{
		Port: viper.GetString("GRPC_PORT"),
	}
	storeConfig := StoreConfig{
		Address:  viper.GetString("STORE_ADDRESS"),
		Password: viper.GetString("STORE_PASSWORD"),
	}
	return AppConfig{
		Grpc:  grpcConfig,
		Store: storeConfig,
	}
}

func setDefaults() {
	viper.SetDefault("GRPC_PORT", 8090)
}
