package main

import (
	"github.com/spf13/viper"
)

type Config struct {
	Port string `mapstructure:"Part"`
}

func loadConfig() (config Config, err error) {
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()
	err = viper.ReadInConfig()
	if err != nil {
		return
	}
	err = viper.Unmarshal(&config)
	return
}
