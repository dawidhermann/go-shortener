// App configuration definition
package config

import "github.com/spf13/viper"

type StoreConfig struct {
	Address  string
	Password string
}

type ApiConfig struct {
	Port string
}

type AppConfig struct {
	Store StoreConfig
	Api   ApiConfig
}

func init() {
	viper.AutomaticEnv()
	viper.SetEnvPrefix("SHORTENER_REDIRECT")
	setDefaults()
}

// Create new application's config
func GetAppConfiguration() AppConfig {
	apiConfig := ApiConfig{
		Port: viper.GetString("API_PORT"),
	}
	storeConfig := StoreConfig{
		Address:  viper.GetString("STORE_ADDRESS"),
		Password: viper.GetString("STORE_PASSWORD"),
	}
	return AppConfig{
		Api:   apiConfig,
		Store: storeConfig,
	}
}

func setDefaults() {
	viper.SetDefault("API_PORT", 8090)
}
