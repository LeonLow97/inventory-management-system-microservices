package config

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/spf13/viper"
)

type Config struct {
	Mode   string `mapstructure:"mode"`
	Server struct {
		URL  string `mapstructure:"url"`
		Port int    `mapstructure:"port"`
	} `mapstructure:"server"`
	JWT struct {
		Secret string `mapstructure:"secret"`
	} `mapstructure:"jwt"`
	AuthService struct {
		URL string `mapstructure:"url"`
	} `mapstructure:"auth_service"`
	InventoryService struct {
		URL string `mapstructure:"url"`
	} `mapstructure:"inventory_service"`
	OrderService struct {
		URL string `mapstructure:"url"`
	} `mapstructure:"order_service"`
}

func LoadConfig() (*Config, error) {
	mode := os.Getenv("MODE")
	if len(mode) == 0 {
		mode = "development"
	}

	viper.SetConfigName(mode)
	viper.SetConfigType("yaml")
	viper.AddConfigPath("/app/config")

	// inject environment variables to override matching keys in configuration files (.yaml)
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Println("config file not found", err)
			return nil, errors.New("config file not found")
		} else {
			log.Println("error reading config file", err)
			return nil, err
		}
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %v", err)
	}

	return &config, nil
}
