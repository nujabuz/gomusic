// Package main Entry point of the application. Responsible for loading configuration, initializing database connections,
// instantiating repositories & use cases, wiring dependencies together (Dependency Injection),
// and starting the HTTP server.
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
	tokenGenerator := security.NewTokenGenerator(cfg.JWT.Issuer, cfg.JWT.ExpiratedInSecond, cfg.JWT.Key)

	userRepository := postgres.NewUserRepository(db)

	userUseCase := usecase.NewUserUseCase(userRepository)
	authUseCase := usecase.NewAuthUseCase(hasher, tokenGenerator, userRepository, userUseCase)

	userHandler := rest.NewUserHandler(userUseCase)
	authhandler := rest.NewAuthHandler(authUseCase)

	app := config.NewApp(
		userHandler,
		authhandler,
	)

	if err := app.Listen(cfg.App.Addr); err != nil {
		slog.Error("fiber", "error", err)
		os.Exit(1)
	}
}
