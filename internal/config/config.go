package config

import (
	"log/slog"
	"os"

	"github.com/spf13/viper"
)

type Config struct {
	Database *ConfigDatabase `mapstructure:"database"`
	App      *ConfigApp      `mapstructure:"app"`
}

func NewConfig() *Config {
	v := viper.New()

	v.SetConfigName("config")
	v.SetConfigType("yaml")

	v.AddConfigPath(".")
	if err := v.ReadInConfig(); err != nil {
		slog.Error("viper", "error", err)
		os.Exit(1)
	}

	cfg := new(Config)
	if err := v.Unmarshal(cfg); err != nil {
		slog.Error("viper", "error", err)
		os.Exit(1)
	}

	return cfg
}
