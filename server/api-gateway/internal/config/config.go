package config

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Mode   string `mapstructure:"mode"`
	Server struct {
		URL  string `mapstructure:"url"`
		Port int    `mapstructure:"port"`
	} `mapstructure:"server"`
	AuthJWTToken struct {
		Name     string `mapstructure:"name"`
		Secret   string `mapstructure:"secret"`
		MaxAge   int    `mapstructure:"max_age"`
		Domain   string `mapstructure:"domain"`
		Secure   bool   `mapstructure:"secure"`
		HTTPOnly bool   `mapstructure:"http_only"`
		Path     string `mapstructure:"path"`
	} `mapstructure:"auth_jwt_token"`
	AuthService struct {
		Name string `mapstructure:"name"`
	} `mapstructure:"auth_service"`
	InventoryService struct {
		Name string `mapstructure:"name"`
	} `mapstructure:"inventory_service"`
	OrderService struct {
		Name string `mapstructure:"name"`
	} `mapstructure:"order_service"`
	HashicorpConsulConfig struct {
		Port    int    `mapstructure:"port"`
		Address string `mapstructure:"address"`
	} `mapstructure:"hashicorp_consul"`
	RedisServer struct {
		Port          int    `mapstructure:"port"`
		Address       string `mapstructure:"address"`
		Password      string `mapstructure:"password"`
		DatabaseIndex int    `mapstructure:"database_index"`
	} `mapstructure:"redis_server"`
	RateLimiting struct {
		BucketLockExpiration int `mapstructure:"bucket_lock_expiration"`
		DistributedLocks     struct {
			Write  string `mapstructure:"write"`
			Read   string `mapstructure:"read"`
			Global string `mapstructure:"global"`
		} `mapstructure:"distributed_locks"`
	} `mapstructure:"rate_limiting"`
}

const (
	ModeDocker  = "docker"
	ModeStaging = "staging" // local kubernetes
)

func LoadConfig() (*Config, error) {
	vpr := viper.New()
	vpr.AutomaticEnv()
	vpr.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.SetConfigType("yaml")

	mode := vpr.GetString("MODE")
	if len(mode) == 0 {
		mode = ModeDocker
	}
	viper.SetConfigName(mode)

	viper.AddConfigPath("/app/config")

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
