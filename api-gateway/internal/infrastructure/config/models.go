package config

import "time"

type Config struct {
	Mode           string               `mapstructure:"mode"`
	Server         ServerConfig         `mapstructure:"server"`
	RedisConfig    RedisConfig          `mapstructure:"redis"`
	AuthifyService AuthifyServiceConfig `mapstructure:"authify_service"`
	AuthJWTToken   AuthJWTTokenConfig   `mapstructure:"auth_jwt_token"`

	// Rate Limiting
	AuthRateLimit AuthRateLimitConfig `mapstructure:"auth"`
}
type ServerConfig struct {
	BaseUrl string `mapstructure:"base_url"`
	Port    int    `mapstructure:"port"`
}
type RedisConfig struct {
	Port          int    `mapstructure:"port"`
	Address       string `mapstructure:"address"`
	Password      string `mapstructure:"password"`
	DatabaseIndex int    `mapstructure:"database_index"`
}
type AuthifyServiceConfig struct {
	BaseUrl string `mapstructure:"base_url"`
	Port    int    `mapstructure:"port"`
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

type RateLimitConfig struct {
	Window      time.Duration `mapstructure:"window"`
	MaxRequests int64         `mapstructure:"max_requests"`
}
type LoginRateLimitConfig struct {
	Key       string          `mapstructure:"key"`
	RateLimit RateLimitConfig `mapstructure:"rate_limit"`
}
type AuthRateLimitConfig struct {
	Login struct {
		Burst LoginRateLimitConfig `mapstructure:"burst"`
		Fail  LoginRateLimitConfig `mapstructure:"fail"`
	} `mapstructure:"login"`
}
