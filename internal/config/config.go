package config

import (
	"log/slog"
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
)

type Config struct {
	DatabaseConfig DatabaseConfig `mapstructure:"database"`
	AppConfig      AppConfig      `mapstructure:"app"`
}

type DatabaseConfig struct {
	Username    string `mapstructure:"username" validate:"required"`
	Password    string `mapstructure:"password" validate:"required"`
	Host        string `mapstructure:"host" validate:"required"`
	Name        string `mapstructure:"name" validate:"required"`
	SSL         string `mapstructure:"ssl" validate:"required"`
	Port        int    `mapstructure:"port" validate:"required,min=1,max=65535"`
	MaxConn     int    `mapstructure:"max" validate:"required,min=1"`
	MaxIdleConn int    `mapstructure:"idle" validate:"required,min=1"`
	MaxLifetime int    `mapstructure:"lifetime" validate:"required,min=1"`
}

type AppConfig struct {
	Addr          string `mapstructure:"addr" validate:"required"`
	Name          string `mapstructure:"name" validate:"required"`
	Prefork       bool   `mapstructure:"prefork"`
	StrictRouting bool   `mapstructure:"strictRouting"`
	CaseSensitive bool   `mapstructure:"caseSensitive"`
	BodyLimit     int    `mapstructure:"bodyLimit" validate:"required,min=1024"`
	ReadTimeout   int    `mapstructure:"readTimeout" validate:"required,min=1"`
	WriteTimeout  int    `mapstructure:"writeTimeout" validate:"required,min=1"`
	IdleTimeout   int    `mapstructure:"idleTimeout" validate:"required,min=1"`
}

func NewConfig(validate *validator.Validate) *Config {
	v := viper.New()

	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath(".")

	cfg := new(Config)

	if err := v.ReadInConfig(); err != nil {
		slog.Error("failed to read config.yaml file", "error", err)
		os.Exit(1)
	}

	if err := v.Unmarshal(cfg); err != nil {
		slog.Error("failed in unmarshal struct config", "error", err)
		os.Exit(1)
	}

	if err := validate.Struct(cfg); err != nil {
		slog.Error("failed struct config validation", "error", err)
		os.Exit(1)
	}

	return cfg
}
