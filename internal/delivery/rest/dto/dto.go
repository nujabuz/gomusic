package dto

import (
	"errors"
	"log/slog"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
	"github.com/nujabuz/music-cloud/internal/domain"
)

type ApiResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
	Error   any    `json:"errors,omitempty"`
}

func ResponseSuccess(ctx fiber.Ctx, status int, message string, data any) error {
	return ctx.Status(status).JSON(ApiResponse{Success: true, Message: message, Data: data})
}

func ResponseFailed(ctx fiber.Ctx, status int, message string, err any) error {
	return ctx.Status(status).JSON(ApiResponse{Success: false, Message: message, Error: err})
}

func TranslateDomainError(ctx fiber.Ctx, err error) error {
	switch {
	case errors.Is(err, domain.ErrInvalidUser):
		return ResponseFailed(ctx, fiber.StatusBadRequest, err.Error(), fiber.ErrBadRequest)
	case errors.Is(err, domain.ErrUserAlreadyExists):
		return ResponseFailed(ctx, fiber.StatusConflict, err.Error(), fiber.ErrConflict)
	case errors.Is(err, domain.ErrUserNotFound):
		return ResponseFailed(ctx, fiber.StatusNotFound, err.Error(), fiber.ErrNotFound)
	default:
		slog.Error("an unexpected error occurred", "error", err)
		return ResponseFailed(ctx, fiber.StatusInternalServerError, "an unexpected error occurred", fiber.ErrInternalServerError)
	}
}

func BindValidation(ctx fiber.Ctx, err error) error {
	validationErrors := new(validator.ValidationErrors)
	if errors.As(err, validationErrors) {
		out := make([]fiber.Map, 0, len(*validationErrors))
		for _, e := range *validationErrors {
			// e.Field() - field name, e.Tag() - failed rule,
			// e.Param() - rule parameter, e.Value() - invalid value
			out = append(out, fiber.Map{
				"field": e.Field(),
				"rule":  e.Tag(),
			})
		}
		return ResponseFailed(ctx, fiber.ErrBadRequest.Code, "failed to validate request body", out)
	}
	return ResponseFailed(ctx, fiber.ErrBadRequest.Code, "failed to parse request body", err)
}
