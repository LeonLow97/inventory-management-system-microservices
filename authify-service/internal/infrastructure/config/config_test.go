package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	testCases := []struct {
		name           string
		mode           string
		configName     string
		configContents string
		env            map[string]string
		wantMode       string
		wantServerName string
		wantServerPort int
		wantError      string
	}{
		{
			name:           "defaults to development and applies config defaults",
			configName:     ModeDevelopment,
			configContents: "server:\n  name: local-authify-service\n",
			wantMode:       ModeDevelopment,
			wantServerName: "local-authify-service",
			wantServerPort: 50051,
		},
		{
			name:           "environment overrides config values",
			mode:           ModeDocker,
			configName:     ModeDocker,
			configContents: "server:\n  name: docker-authify-service\n  port: 50051\n",
			env:            map[string]string{"SERVER_PORT": "6000"},
			wantMode:       ModeDocker,
			wantServerName: "docker-authify-service",
			wantServerPort: 6000,
		},
		{
			name:      "rejects unknown mode",
			mode:      "staging",
			wantError: `unknown MODE "staging"`,
		},
		{
			name:      "reports a missing mode file",
			mode:      ModeProduction,
			wantError: `no config file found for mode "production"`,
		},
		{
			name:           "reports malformed yaml",
			mode:           ModeDevelopment,
			configName:     ModeDevelopment,
			configContents: "server: [\n",
			wantError:      "failed reading config",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			configDir := t.TempDir()
			if testCase.configName != "" {
				writeConfigFile(t, configDir, testCase.configName, testCase.configContents)
			}
			t.Chdir(configDir)
			t.Setenv("MODE", testCase.mode)
			t.Setenv("SERVER_NAME", "")
			t.Setenv("SERVER_PORT", "")
			for key, value := range testCase.env {
				t.Setenv(key, value)
			}

			got, err := LoadConfig()
			if testCase.wantError != "" {
				if err == nil || !strings.Contains(err.Error(), testCase.wantError) {
					t.Fatalf("LoadConfig() error = %v, want message containing %q", err, testCase.wantError)
				}
				if got != nil {
					t.Fatalf("LoadConfig() config = %#v, want nil on error", got)
				}
				return
			}
			if err != nil {
				t.Fatalf("LoadConfig() error = %v", err)
			}
			if got.Mode != testCase.wantMode || got.Server.Name != testCase.wantServerName || got.Server.Port != testCase.wantServerPort {
				t.Errorf("LoadConfig() = (mode=%q, server=(%q, %d)), want (mode=%q, server=(%q, %d))",
					got.Mode, got.Server.Name, got.Server.Port,
					testCase.wantMode, testCase.wantServerName, testCase.wantServerPort)
			}
		})
	}
}

func TestGetConfig(t *testing.T) {
	configDir := t.TempDir()
	writeConfigFile(t, configDir, ModeDevelopment, "server:\n  name: singleton-authify-service\n")
	t.Chdir(configDir)
	t.Setenv("MODE", ModeDevelopment)
	t.Setenv("SERVER_NAME", "")
	t.Setenv("SERVER_PORT", "")

	got, err := GetConfig()
	if err != nil {
		t.Fatalf("GetConfig() error = %v", err)
	}
	if got == nil || got.Server.Name != "singleton-authify-service" {
		t.Fatalf("GetConfig() = %#v, want configured singleton", got)
	}
}

func writeConfigFile(t *testing.T, directory, name, contents string) {
	t.Helper()

	configDirectory := filepath.Join(directory, "config", "environments")
	if err := os.MkdirAll(configDirectory, 0o755); err != nil {
		t.Fatalf("create config directory: %v", err)
	}
	if err := os.WriteFile(filepath.Join(configDirectory, name+".yaml"), []byte(contents), 0o600); err != nil {
		t.Fatalf("write config file: %v", err)
	}
}
