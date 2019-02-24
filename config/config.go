package config

import (
	"goqueue/helper"

	"github.com/spf13/viper"
)

func Init() {
	viper.SetConfigName("config")
	viper.SetConfigType("yml")
	viper.AddConfigPath(".")
	helper.FailOnError(viper.ReadInConfig(), "Failed to read config file")
}
