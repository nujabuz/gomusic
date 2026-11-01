package main

import (
	"log/slog"
	"os"

	"github.com/nujabuz/gomusic/internal/config"
	"github.com/nujabuz/gomusic/internal/delivery/rest"
	"github.com/nujabuz/gomusic/internal/repository/postgres"
	"github.com/nujabuz/gomusic/internal/security"
	"github.com/nujabuz/gomusic/internal/usecase"
)

func main() {
	cfg := config.NewConfig()
	db := config.NewDatabase(cfg.Database)

	hasher := security.NewPasswordHasher()

	userRepository := postgres.NewUserRepository(db)
	userUseCase := usecase.NewUserUseCase(userRepository, hasher)
	userHandler := rest.NewUserHandler(userUseCase)

	app := config.NewApp(
		userHandler,
	)

	if err := app.Listen(cfg.App.Addr); err != nil {
		slog.Error("fiber", "error", err)
		os.Exit(1)
	}
}
