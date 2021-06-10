package env

import (
	"fmt"

	"github.com/spf13/viper"
)

// handle global config vars
func GetEnvByKey(key string) string {
	viper.AddConfigPath(".")
	viper.SetConfigName("app")
	viper.SetConfigType("env")
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Errorf("Failed to read config file: %s \n", err)
		return ""
	}
	value := viper.Get(key).(string)

	return value
}
