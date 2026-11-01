package rest

import (
	"errors"
	"log/slog"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
	"github.com/nujabuz/gomusic/internal/entity"
)

type envelope struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
	Error   any    `json:"error,omitempty"`
}

func ResponseSuccess(ctx fiber.Ctx, status int, message string, data any) error {
	return ctx.Status(status).JSON(&envelope{Success: true, Message: message, Data: data})
}

func ResponseFailed(ctx fiber.Ctx, status int, message string, err any) error {
	return ctx.Status(status).JSON(&envelope{Success: false, Message: message, Error: err})
}

func EntityErrorTranslator(ctx fiber.Ctx, err error) error {
	switch {
	// user error
	case errors.Is(err, entity.ErrUserNotFound):
		return ResponseFailed(ctx, fiber.StatusNotFound, err.Error(), fiber.ErrNotFound)
	case errors.Is(err, entity.ErrUserAlreadyExists):
		return ResponseFailed(ctx, fiber.StatusConflict, err.Error(), fiber.ErrConflict)

	// unexpected error
	default:
		slog.Error("an unexpected error occurred", "error", err)
		return ResponseFailed(ctx, fiber.StatusInternalServerError, entity.ErrUnexpected.Error(), fiber.ErrInternalServerError)
	}
}

func ResponseFailedParseRequest(ctx fiber.Ctx, err error) error {
	var validationErrors validator.ValidationErrors
	if errors.As(err, &validationErrors) {
		out := make([]fiber.Map, 0, len(validationErrors))
		for _, e := range validationErrors {
			// e.Field() - field name, e.Tag() - failed rule,
			// e.Param() - rule parameter, e.Value() - invalid value
			out = append(out, fiber.Map{
				"field": e.Field(),
				"rule":  e.Tag(),
			})
		}
		return ResponseFailed(ctx, fiber.StatusBadRequest, entity.ErrParseRequestBody.Error(), out)
	}
	return ResponseFailed(ctx, fiber.StatusBadRequest, entity.ErrParseRequestBody.Error(), err)
}
