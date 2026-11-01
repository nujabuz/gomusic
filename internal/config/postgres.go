package config

import (
	"fmt"
	"log"
	"log/slog"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type ConfigDatabase struct {
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	Host     string `mapstructure:"host"`
	SSL      string `mapstructure:"ssl"`
	DBName   string `mapstructure:"dbname"`
	Port     int    `mapstructure:"port"`
}

func NewDatabase(cfg *ConfigDatabase) *gorm.DB {
	dsn := fmt.Sprintf("user=%s password=%s host=%s port=%d dbname=%s sslmode=%s", cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DBName, cfg.SSL)
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second,   // Slow SQL threshold
			LogLevel:                  logger.Silent, // Log level
			IgnoreRecordNotFoundError: true,          // Ignore ErrRecordNotFound error for logger
			ParameterizedQueries:      true,          // Don't include params in the SQL log
			Colorful:                  false,         // Disable color
		},
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger:         newLogger,
		TranslateError: true,
	})
	if err != nil {
		slog.Error("gorm", "error", err)
		os.Exit(1)
	}

	return db
}
