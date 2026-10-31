package config

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/limiter"
	"github.com/nujabuz/music-cloud/internal/delivery/rest"
	"github.com/nujabuz/music-cloud/internal/delivery/rest/dto"
)

type structValidator struct {
	validate *validator.Validate
}

func (v *structValidator) Validate(out any) error {
	return v.validate.Struct(out)
}

type App struct {
	Fiber          *fiber.App
	Config         AppConfig
	UserHandler    rest.UserHandler
	AuthMiddleware fiber.Handler
}

func NewApp(cfg AppConfig, vld *validator.Validate, userHandler rest.UserHandler, authMiddleware fiber.Handler) *App {
	app := fiber.New(fiber.Config{
		AppName:         cfg.Name,
		StrictRouting:   cfg.StrictRouting,
		CaseSensitive:   cfg.CaseSensitive,
		BodyLimit:       cfg.BodyLimit,
		IdleTimeout:     time.Duration(cfg.IdleTimeout) * time.Second,
		ReadTimeout:     time.Duration(cfg.ReadTimeout) * time.Second,
		WriteTimeout:    time.Duration(cfg.WriteTimeout) * time.Second,
		StructValidator: &structValidator{validate: vld},
	})
	return &App{
		Config:         cfg,
		Fiber:          app,
		UserHandler:    userHandler,
		AuthMiddleware: authMiddleware,
	}
}

func (a *App) Start() error {
	// 1. Channel to receive OS signals (Ctrl+C, kill, docker stop, etc.)
	shutdownChan := make(chan os.Signal, 1)
	signal.Notify(shutdownChan, os.Interrupt, syscall.SIGTERM)

	// 2. Start server in a background goroutine
	serverErr := make(chan error, 1)
	go func() {
		slog.Info("server starting", "addr", a.Config.Addr)
		if err := a.Fiber.Listen(a.Config.Addr, fiber.ListenConfig{
			EnablePrefork: a.Config.Prefork,
		}); err != nil {
			serverErr <- err
		}
	}()

	// 3. Block until either the server fails to start or we receive a shutdown signal
	select {
	case err := <-serverErr:
		return err
	case sig := <-shutdownChan:
		slog.Info("shutdown signal received", "signal", sig.String())
	}

	// 4. Give in-flight requests up to 10 seconds to finish
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	slog.Info("shutting down HTTP server gracefully...")
	return a.Fiber.ShutdownWithContext(ctx)
}

func (a *App) SetupRoute() {
	api := a.Fiber.Group("/api/v1")

	a.setupAutheRoute(api)
	a.setupAuthoRoute(api)
}

func (a *App) setupAutheRoute(api fiber.Router) {
	auth := api.Group("/auth")
	auth.Use(limiter.New(limiter.Config{
		Max:               3,
		Expiration:        time.Minute * 1,
		LimiterMiddleware: limiter.SlidingWindow{},
		LimitReached: func(ctx fiber.Ctx) error {
			return dto.ResponseFailed(ctx, fiber.StatusTooManyRequests, "Limit reached (3 requests). Please wait 1 minute.", "Too Many Requests")
		},
	}))

	auth.Post("/register", a.UserHandler.Register)
	auth.Post("/login", a.UserHandler.Login)
}

func (a *App) setupAuthoRoute(api fiber.Router) {
	music := api.Group("/music")

	music.Use(a.AuthMiddleware)

	music.Get("/ping", func(ctx fiber.Ctx) error {
		username, ok := fiber.ValueFromContext[string](ctx, "username")
		if ok {
			return dto.ResponseSuccess(ctx, fiber.StatusOK, "pong!", map[string]string{"username": username})
		}
		return dto.ResponseFailed(ctx, fiber.StatusUnauthorized, "not pong!", fiber.ErrUnauthorized)
	})
}
