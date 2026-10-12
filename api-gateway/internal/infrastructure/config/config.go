package config

import (
	"fmt"
	"strings"
	"sync"

	"github.com/spf13/viper"
)

const (
	ModeDevelopment = "development"
	ModeDocker      = "docker"
	ModeKubernetes  = "kubernetes"
	ModeProduction  = "production"
)

var modeConfigFiles = map[string]string{
	ModeDevelopment: ModeDevelopment,
	ModeDocker:      ModeDocker,
	ModeKubernetes:  "kind",
	ModeProduction:  "production-aws",
}

func LoadConfig() (*Config, error) {
	vpr := viper.New()
	vpr.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	vpr.AutomaticEnv()

	envKeys := []string{
		"server.base_url",
		"server.port",
		"authify_service.base_url",
		"authify_service.port",
		"auth_jwt_token.name",
		"auth_jwt_token.secret",
		"auth_jwt_token.max_age",
		"auth_jwt_token.domain",
		"auth_jwt_token.secure",
		"auth_jwt_token.http_only",
		"auth_jwt_token.path",
		"redis.port",
		"redis.database_index",
		"redis.address",
		"redis.password",
	}

	for _, key := range envKeys {
		if err := vpr.BindEnv(key); err != nil {
			return nil, fmt.Errorf("failed binding env var %s: %w", key, err)
		}
	}

	mode := vpr.GetString("MODE")
	if mode == "" {
		mode = ModeDevelopment
	}

	configName, ok := modeConfigFiles[mode]
	if !ok {
		return nil, fmt.Errorf("unknown MODE %q", mode)
	}

	// Load mode config
	vpr.SetConfigName(configName)
	vpr.SetConfigType("yaml")
	vpr.AddConfigPath("./config/environments")
	if err := vpr.MergeInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read mode config: %w", err)
	}

	// Load ratelimit config
	ratelimitVpr := viper.New()
	ratelimitVpr.SetConfigName("auth")
	ratelimitVpr.SetConfigType("yaml")
	ratelimitVpr.AddConfigPath("./config/rate-limit")
	if err := ratelimitVpr.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read ratelimit config: %w", err)
	}

	// Merge ratelimit into main config
	if err := vpr.MergeConfigMap(ratelimitVpr.AllSettings()); err != nil {
		return nil, fmt.Errorf("failed to merge ratelimit config: %w", err)
	}

	var c Config
	c.setDefaultValues(vpr)
	if err := vpr.Unmarshal(&c); err != nil {
		return nil, fmt.Errorf("failed unmarshaling config: %w", err)
	}
	c.Mode = mode
	return &c, nil
}

var (
	configOnce   sync.Once
	loadedConfig *Config
	configErr    error
)

// GetConfig returns the global config instance (loaded once).
// Safe for concurrent use. The caller handles any load error.
func GetConfig() (*Config, error) {
	configOnce.Do(func() {
		loadedConfig, configErr = LoadConfig()
	})
	return loadedConfig, configErr
}

func (c *Config) setDefaultValues(vpr *viper.Viper) {
	vpr.SetDefault("server.base_url", "localhost")
	vpr.SetDefault("server.port", 80)
}
