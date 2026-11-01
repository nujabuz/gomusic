package config

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
	"github.com/nujabuz/gomusic/internal/delivery/rest"
	"github.com/nujabuz/gomusic/internal/entity"
)

type ConfigApp struct {
	Addr string `mapstructure:"addr"`
}

type structValidator struct {
	validate *validator.Validate
}

func (v *structValidator) Validate(out any) error {
	return v.validate.Struct(out)
}

func NewApp(userHandler *rest.UserHandler) *fiber.App {
	app := fiber.New(fiber.Config{
		StructValidator: &structValidator{validate: validator.New()},
	})

	api := app.Group("/api/v1")

	userHandler.Routes(api)

	app.Use(func(ctx fiber.Ctx) error {
		return rest.ResponseFailed(ctx, fiber.StatusNotFound, entity.ErrNotFound.Error(), fiber.ErrNotFound) // => 404 "Not Found"
	})

	return app
}
