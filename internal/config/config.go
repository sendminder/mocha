package config

import (
	"log/slog"

	"github.com/spf13/viper"
)

func Init() {
	loadConfig()
}

func loadConfig() {
	slog.Info("Config Loaded")
	viper.SetConfigFile("config.yaml")
	if err := viper.ReadInConfig(); err != nil {
		slog.Error("Failed to read config file", "err", err)
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
