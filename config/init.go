package config

import (
	"github.com/spf13/viper"
)

func Init() error {
	viper.AddConfigPath("./")
	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	return viper.ReadInConfig()
}
