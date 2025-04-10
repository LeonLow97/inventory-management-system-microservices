package config

import (
	"fmt"
	"log"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Mode                  string                 `mapstructure:"mode"`
	Server                ServerConfig           `mapstructure:"server"`
	AuthJWTToken          AuthJWTTokenConfig     `mapstructure:"auth_jwt_token"`
	AuthService           AuthServiceConfig      `mapstructure:"auth_service"`
	InventoryService      InventoryServiceConfig `mapstructure:"inventory_service"`
	OrderService          OrderServiceConfig     `mapstructure:"order_service"`
	HashicorpConsulConfig HashicorpConsulConfig  `mapstructure:"hashicorp_consul"`
	RedisServer           RedisServerConfig      `mapstructure:"redis_server"`
	RateLimiting          RateLimitingConfig     `mapstructure:"rate_limiting"`
	AdminWhitelistedIPs   []string               `mapstructure:"admin_whitelisted_ips"`
}

type ServerConfig struct {
	URL  string `mapstructure:"url"`
	Port int    `mapstructure:"port"`
}

type AuthJWTTokenConfig struct {
	Name     string `mapstructure:"name"`
	Secret   string `mapstructure:"secret"`
	MaxAge   int    `mapstructure:"max_age"`
	Domain   string `mapstructure:"domain"`
	Secure   bool   `mapstructure:"secure"`
	HTTPOnly bool   `mapstructure:"http_only"`
	Path     string `mapstructure:"path"`
}

type AuthServiceConfig struct {
	Name string `mapstructure:"name"`
}

type InventoryServiceConfig struct {
	Name string `mapstructure:"name"`
}

type OrderServiceConfig struct {
	Name string `mapstructure:"name"`
}

type HashicorpConsulConfig struct {
	Port    int    `mapstructure:"port"`
	Address string `mapstructure:"address"`
}

type RedisServerConfig struct {
	Port          int    `mapstructure:"port"`
	Address       string `mapstructure:"address"`
	Password      string `mapstructure:"password"`
	DatabaseIndex int    `mapstructure:"database_index"`
}

type RateLimitingConfig struct {
	BucketLockExpiration int `mapstructure:"bucket_lock_expiration"`
	DistributedLocks     struct {
		Write  string `mapstructure:"write"`
		Read   string `mapstructure:"read"`
		Global string `mapstructure:"global"`
	} `mapstructure:"distributed_locks"`
}

const (
	ModeDocker     = "docker"
	ModeKubernetes = "kubernetes" // local kubernetes
)

// LoadConfig reads configuration based on the current mode.
// It loads environment variables, reads the config file, and Unmarshal it into a Config struct
func LoadConfig() (*Config, error) {
	// Create a new Viper instance to manage configuration settings
	vpr := viper.New()

	// Replace '.' with '_' in environment variable names to match common conventions
	// E.g., "DATABASE.HOST" becomes "DATABASE_HOST"
	vpr.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// Automatically bind environment variables to Viper keys
	// This allows configuration values to be overridden by environment variables
	vpr.AutomaticEnv()

	// Retrieve the "MODE" environment variable to determine the configuration mode
	// If not set, default to Development mode
	mode := vpr.GetString("MODE")
	if mode == "" {
		mode = ModeDocker
	}

	// Set the configuration file name based on the mode (e.g., "development.yaml")
	vpr.SetConfigName(mode)

	// Add the directory "/app/config" as a location to search for the configuration file
	vpr.AddConfigPath("/app/config")

	if err := vpr.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return nil, fmt.Errorf("config file not found with error: %v", err)
		}
		log.Println("Failed to read config file using Viper with error:", err)
		return nil, err
	}

	// Unmarshal the configuration file into the Config struct
	var config Config
	if err := vpr.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config with error: %v", err)
	}

	fmt.Println("Config -->", config)

	return &config, nil
}
