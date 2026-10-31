package config

import (
	"database/sql"
	"fmt"
	"log/slog"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewDatabase(cfg *DatabaseConfig) (*gorm.DB, *sql.DB) {
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		cfg.Username,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Name,
		cfg.SSL,
	)

	db, err := gorm.Open(postgres.Open(dsn))
	if err != nil {
		slog.Error("failed when initializing db", "error", err)
		os.Exit(1)
	}

	connection, err := db.DB()
	if err != nil {
		slog.Error("failed connection to database", "error", err)
		os.Exit(1)
	}

	connection.SetMaxOpenConns(cfg.MaxConn)
	connection.SetMaxIdleConns(cfg.MaxIdleConn)
	connection.SetConnMaxLifetime(time.Duration(cfg.MaxLifetime) * time.Minute)

	return db, connection
}
