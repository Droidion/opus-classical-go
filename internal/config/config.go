package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	Port string `mapstructure:"Port"`
}

var config *Config

// Load reads in config file and ENV variables if set.
func Load() error {
	var err error
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()

	if err = viper.ReadInConfig(); err != nil {
		return err
	}

	config = &Config{}
	err = viper.Unmarshal(config)
	return err
}

// Get returns the loaded config
func Get() *Config {
	if config == nil {
		panic("Config not loaded")
	}
	return config
}
