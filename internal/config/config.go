package config

import (
	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
)

type Config struct {
	DatabaseConfig DatabaseConfig `mapstructure:"database"`
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

func NewConfig(validate *validator.Validate) (*Config, error) {
	v := viper.New()

	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath(".")

	cfg := new(Config)

	if err := v.ReadInConfig(); err != nil {
		return nil, err
	}

	if err := v.Unmarshal(cfg); err != nil {
		return nil, err
	}

	if err := validate.Struct(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}
