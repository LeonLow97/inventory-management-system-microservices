package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	AWS AWS `mapstructure:"aws"`
}

type AWS struct {
	Region           string         `mapstructure:"region"`
	Endpoint         string         `mapstructure:"endpoint"`
	Bucket           string         `mapstructure:"bucket"`
	S3ForcePathStyle bool           `mapstructure:"s3_force_path_style"`
	Credentials      AWSCredentials `mapstructure:"credentials"`
}

type AWSCredentials struct {
	Token  string
	ID     string `mapstructure:"id"`
	Secret string `mapstructure:"secret"`
}

func LoadConfig() (Config, error) {
	var cfg Config

	viper.SetConfigName("development")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("/app/config")

	if err := viper.ReadInConfig(); err != nil {
		return cfg, err
	}

	if err := viper.Unmarshal(&cfg); err != nil {
		return cfg, err
	}

	return cfg, nil
}
