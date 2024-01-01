package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	Development Development `mapstructure:"development"`
}

type Development struct {
	URL string `mapstructure:"url"`
}

func LoadConfig(environment string) (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	viper.AddConfigPath("./config")

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	return &config, nil
}
