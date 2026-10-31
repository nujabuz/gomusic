package main

import (
	"log/slog"

	"github.com/go-playground/validator/v10"
	"github.com/nujabuz/music-cloud/internal/config"
	"github.com/nujabuz/music-cloud/internal/delivery/rest"
	"github.com/nujabuz/music-cloud/internal/repository/postgres"
	"github.com/nujabuz/music-cloud/internal/security"
	"github.com/nujabuz/music-cloud/internal/usecase"
)

func main() {
	vld := validator.New()
	cfg := config.NewConfig(vld)
	hsr := security.NewPasswordHasher()

	db, connection := config.NewDatabase(&cfg.DatabaseConfig)
	defer func() {
		slog.Info("closing database connection pool")
		if err := connection.Close(); err != nil {
			slog.Error("error closing database", "error", err)
		}
	}()

	userRepository := postgres.NewUserRepository(db)
	userUseCase := usecase.NewUserUseCase(userRepository, hsr)
	userHandler := rest.NewUserHandler(userUseCase)

	app := config.NewApp(cfg.AppConfig, vld, userHandler)
	app.SetupRoute()
	if err := app.Start(); err != nil {
		slog.Error("server shutdown with error", "error", err)
	} else {
		slog.Info("server exited cleanly")
	}

}
