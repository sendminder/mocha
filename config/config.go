package config

import (
	"log"

	"github.com/spf13/viper"
)

func init() {
	loadConfig()
}

func loadConfig() {
	log.Println("Config Loaded")
	viper.SetConfigFile("./config/config.yaml")
	if err := viper.ReadInConfig(); err != nil {
		log.Println("Failed to read config file:", err)
		return
	}
}

func GetString(key string) string {
	return viper.GetString(key)
}

func GetInt(key string) int {
	return viper.GetInt(key)
}

func GetBool(key string) bool {
	return viper.GetBool(key)
}
