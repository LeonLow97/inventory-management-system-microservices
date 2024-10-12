package config

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Server struct {
		Name string `mapstructure:"name"`
		Port int    `mapstructure:"port"`
	} `mapstructure:"server"`
	InventoryService InventoryService `mapstructure:"inventory_service"`
	KafkaConfig      KafkaConfig      `mapstructure:"kafka"`
	PostgresConfig   PostgresConfig   `mapstructure:"postgres"`
}

type InventoryService struct {
	URL string `mapstructure:"url"`
}

type KafkaConfig struct {
	BrokerAddress string `mapstructure:"broker_address"`
}

type PostgresConfig struct {
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	DB       string `mapstructure:"db"`
}

const (
	ModeDevelopment = "development"
	ModeDocker      = "docker"
	ModeStaging     = "staging" // local kubernetes
)

func LoadConfig() (*Config, error) {
	vpr := viper.New()
	vpr.SetEnvKeyReplacer(strings.NewReplacer(".", "_")) // Replace '.' with '_' for environment variables
	vpr.AutomaticEnv()

	mode := vpr.GetString("MODE")
	if len(mode) == 0 {
		mode = "development" // Default mode
	}
	vpr.SetConfigName(mode)
	vpr.AddConfigPath("/app/config")

	if err := vpr.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Println("config file not found", err)
			return nil, errors.New("config file not found")
		}
		log.Println("error reading config file", err)
		return nil, err
	}

	var config Config
	if err := vpr.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %v", err)
	}

	return &config, nil
}
