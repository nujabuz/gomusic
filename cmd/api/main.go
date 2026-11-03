// Package main Entry point of the application. Responsible for loading configuration, initializing database connections,
// instantiating repositories & use cases, wiring dependencies together (Dependency Injection),
// and starting the HTTP server.
package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/nujabuz/gomusic/internal/config"
	"github.com/nujabuz/gomusic/internal/delivery/rest"
	"github.com/nujabuz/gomusic/internal/repository/postgres"
	"github.com/nujabuz/gomusic/internal/security"
	"github.com/nujabuz/gomusic/internal/usecase"
)

func main() {
	cfg := config.NewConfig()
	db, dbConn := config.NewDatabase(cfg.Database)

	hasher := security.NewPasswordHasher()
	tokenGenerator := security.NewTokenGenerator(cfg.JWT.Issuer, cfg.JWT.ExpiratedInSecond, cfg.JWT.Key)

	userRepository := postgres.NewUserRepository(db)

	userUseCase := usecase.NewUserUseCase(userRepository)
	authUseCase := usecase.NewAuthUseCase(hasher, tokenGenerator, userRepository, userUseCase)

	userHandler := rest.NewUserHandler(userUseCase, tokenGenerator)
	authhandler := rest.NewAuthHandler(authUseCase)

	app := config.NewApp(
		userHandler,
		authhandler,
	)
	isShuttingDown := new(atomic.Bool)

	app.Hooks().OnPostShutdown(func(error) error {
		slog.Info("shutting down connection database...")
		err := dbConn.Close()
		if err != nil {
			slog.Error("database shutdown failed with error:", "error", err)
		} else {
			slog.Info("database connection pool closed cleanly")
		}
		return nil
	})

	// Start
	go func() {
		if err := app.Listen(":3000", fiber.ListenConfig{
			DisableStartupMessage: true,
		}); err != nil {
			slog.Error("listen error", "error", err)
		} else {
			slog.Info("successfully listened on", "addr", cfg.App.Addr)
		}
	}()

	// wait for signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	sig := <-quit
	slog.Info("get a signal", "signal", sig)
	slog.Info("shutting down...l")

	// Begin graceful shutdown
	isShuttingDown.Store(true)
	time.Sleep(5 * time.Second) // Let LB drain

	ctx, cancel := context.WithTimeout(context.Background(), 25*time.Second)
	defer cancel()

	done := make(chan error, 1)
	go func() { done <- app.Shutdown() }()

	select {
	case <-done:
		slog.Info("Shutdown complete")
	case <-ctx.Done():
		slog.Info("Shutdown timed out")
	}

	os.Exit(0)
}
