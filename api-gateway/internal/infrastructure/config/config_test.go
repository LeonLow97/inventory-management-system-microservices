package config

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

// writeTempConfig creates a config file in config/environments using the runtime mode's file name.
func writeTempConfig(t *testing.T, baseDir, mode, yaml string) {
	t.Helper()
	cfgDir := filepath.Join(baseDir, "config", "environments")
	if err := os.MkdirAll(cfgDir, 0755); err != nil {
		t.Fatal(err)
	}
	file := filepath.Join(cfgDir, modeConfigFiles[mode]+".yaml")
	if err := os.WriteFile(file, []byte(yaml), 0644); err != nil {
		t.Fatal(err)
	}
}

func writeTempRateLimitConfig(t *testing.T, baseDir string) {
	t.Helper()
	cfgDir := filepath.Join(baseDir, "config", "rate-limit")
	if err := os.MkdirAll(cfgDir, 0755); err != nil {
		t.Fatal(err)
	}
	file := filepath.Join(cfgDir, "auth.yaml")
	const rateLimitYAML = `auth:
  login:
    burst:
      key: "login:burst:%s:%s"
      rate_limit:
        window: 1s
        max_requests: 1
    fail:
      key: "login:fail:%s:%s"
      rate_limit:
        window: 60s
        max_requests: 10
`
	if err := os.WriteFile(file, []byte(rateLimitYAML), 0644); err != nil {
		t.Fatal(err)
	}
}

// withWorkingDir changes to a temporary working directory for the duration of the test.
// It automatically restores the original directory afterward.
func withWorkingDir(t *testing.T, f func(tempDir string)) {
	orig, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := os.Chdir(orig); err != nil {
			t.Errorf("restore working directory: %v", err)
		}
	})

	tempDir := t.TempDir()
	if err := os.Chdir(tempDir); err != nil {
		t.Fatal(err)
	}
	f(tempDir)
}

// withEnv sets an environment variable for the duration of the test.
// If value is empty, it unsets the variable.
func withEnv(t *testing.T, key, value string) {
	t.Helper()
	if value == "" {
		original, existed := os.LookupEnv(key)
		if err := os.Unsetenv(key); err != nil {
			t.Fatal(err)
		}
		t.Cleanup(func() {
			if existed {
				if err := os.Setenv(key, original); err != nil {
					t.Errorf("restore environment variable %s: %v", key, err)
				}
			} else if err := os.Unsetenv(key); err != nil {
				t.Errorf("unset environment variable %s: %v", key, err)
			}
		})
		return
	}
	t.Setenv(key, value)
}

func Test_LoadConfig(t *testing.T) {
	const validYAML = `
mode: development
server:
  base_url: "127.0.0.1"
  port: 8080
`

	t.Run("development mode", func(t *testing.T) {
		t.Run("config file exists", func(t *testing.T) {
			withWorkingDir(t, func(tempDir string) {
				writeTempConfig(t, tempDir, ModeDevelopment, validYAML)
				writeTempRateLimitConfig(t, tempDir)
				withEnv(t, "MODE", ModeDevelopment)
				cfg, err := LoadConfig()
				assert.NoError(t, err)
				assert.NotNil(t, cfg)
				assert.Equal(t, ModeDevelopment, cfg.Mode)
				assert.Equal(t, time.Second, cfg.AuthRateLimit.Login.Burst.RateLimit.Window)
				assert.Equal(t, int64(1), cfg.AuthRateLimit.Login.Burst.RateLimit.MaxRequests)
				assert.Equal(t, time.Minute, cfg.AuthRateLimit.Login.Fail.RateLimit.Window)
				assert.Equal(t, int64(10), cfg.AuthRateLimit.Login.Fail.RateLimit.MaxRequests)
			})
		})

		t.Run("config file does not exist", func(t *testing.T) {
			withWorkingDir(t, func(_ string) {
				withEnv(t, "MODE", ModeDevelopment)
				cfg, err := LoadConfig()
				assert.Error(t, err)
				assert.Nil(t, cfg)
			})
		})

		t.Run("Mode env var not set, defaults to development", func(t *testing.T) {
			withWorkingDir(t, func(tempDir string) {
				writeTempConfig(t, tempDir, ModeDevelopment, validYAML)
				writeTempRateLimitConfig(t, tempDir)
				withEnv(t, "MODE", "") // unset
				cfg, err := LoadConfig()
				assert.NoError(t, err)
				assert.NotNil(t, cfg)
				assert.Equal(t, ModeDevelopment, cfg.Mode)
			})
		})

		t.Run("unknown Mode env var", func(t *testing.T) {
			withWorkingDir(t, func(_ string) {
				withEnv(t, "MODE", "unknown")
				cfg, err := LoadConfig()
				assert.Error(t, err)
				assert.Nil(t, cfg)
			})
		})
	})
}

func Test_LoadConfigRecognizesEveryMode(t *testing.T) {
	testCases := []struct {
		name string
		mode string
	}{
		{name: "development", mode: ModeDevelopment},
		{name: "docker", mode: ModeDocker},
		{name: "kind Kubernetes", mode: ModeKubernetes},
		{name: "AWS production", mode: ModeProduction},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			withWorkingDir(t, func(tempDir string) {
				writeTempConfig(t, tempDir, testCase.mode, "server:\n  port: 9000\n")
				writeTempRateLimitConfig(t, tempDir)
				withEnv(t, "MODE", testCase.mode)

				cfg, err := LoadConfig()
				if err != nil {
					t.Fatalf("LoadConfig() error = %v", err)
				}
				if cfg.Mode != testCase.mode || cfg.Server.Port != 9000 {
					t.Errorf("loaded config mode/port = (%q, %d), want (%q, 9000)", cfg.Mode, cfg.Server.Port, testCase.mode)
				}
			})
		})
	}
}

func Test_LoadConfig_EnvOverridesNestedValues(t *testing.T) {
	const productionYAML = `
mode: production
server:
  base_url: "0.0.0.0"
  port: 80
authify_service:
  base_url: "0.0.0.0"
  port: 50051
auth_jwt_token:
  name: authify-token
  max_age: 3600
  secure: false
  http_only: true
  path: /
  secret: CHANGE_ME
  domain: CHANGE_ME
redis:
  port: 6379
  database_index: 0
  address: CHANGE_ME
  password: CHANGE_ME
`

	withWorkingDir(t, func(tempDir string) {
		writeTempConfig(t, tempDir, ModeProduction, productionYAML)
		writeTempRateLimitConfig(t, tempDir)
		withEnv(t, "MODE", ModeProduction)
		withEnv(t, "SERVER_PORT", "8080")
		withEnv(t, "AUTHIFY_SERVICE_BASE_URL", "10.0.3.10")
		withEnv(t, "REDIS_ADDRESS", "authify-redis.internal")
		withEnv(t, "AUTH_JWT_TOKEN_DOMAIN", "localhost")

		cfg, err := LoadConfig()
		assert.NoError(t, err)
		assert.NotNil(t, cfg)
		assert.Equal(t, 8080, cfg.Server.Port)
		assert.Equal(t, "10.0.3.10", cfg.AuthifyService.BaseUrl)
		assert.Equal(t, "authify-redis.internal", cfg.RedisConfig.Address)
		assert.Equal(t, "localhost", cfg.AuthJWTToken.Domain)
	})
}

func Test_setDefaultValues(t *testing.T) {
	var cfg Config
	vpr := viper.New()

	cfg.setDefaultValues(vpr)
	assert.Equal(t, "localhost", vpr.GetString("server.base_url"))
	assert.Equal(t, 80, vpr.GetInt("server.port"))
}
