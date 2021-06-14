package config

import (
	"github.com/spf13/viper"
	LOGS "go_svelte_lighthouse/logs"
)

// handle global config vars
func GetEnvByKey(key string) string {
	viper.AddConfigPath(".")
	viper.SetConfigName("app")
	viper.SetConfigType("env")
	err := viper.ReadInConfig()
	if err != nil {
		LOGS.ErrorLogger.Fatalf("Failed to read config file: %s \n", err)
		return ""
	}
	value := viper.Get(key).(string)

	return value
}
