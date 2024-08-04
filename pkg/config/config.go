package config

import (
	"log"

	"github.com/spf13/viper"
)

var (
	config *viper.Viper
)

const (
	USERID        = "userId"
	USERTYPE      = "userType"
	USERNAME      = "username"
	AUTHORIZATION = "Authorization"
	ADMIN         = "ADMIN"
	CUSTOMER      = "CUSTOMER"
)

// Init is an exported method that takes the environment starts the viper
// (external lib) and returns the configuration struct.
func Load(env string, configPaths ...string) {
	var err error
	config = viper.New()
	config.SetConfigType("yaml")
	config.SetConfigName(env)
	config.AddConfigPath("config/")
	config.AddConfigPath("../../")
	config.AddConfigPath(".")
	if len(configPaths) != 0 {
		for _, path := range configPaths {
			config.AddConfigPath(path)
		}
	}
	err = config.ReadInConfig()
	if err != nil {
		log.Fatal("error on parsing configuration file", err)
		return
	}
}

func GetConfig() *viper.Viper {
	return config
}
