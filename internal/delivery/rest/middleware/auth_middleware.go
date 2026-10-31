package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v3"
	"github.com/nujabuz/music-cloud/internal/delivery/rest/dto"
	"github.com/nujabuz/music-cloud/internal/domain"
)

func NewAuthMiddleware(tg domain.TokenGenerator) fiber.Handler {
	return func(ctx fiber.Ctx) error {
		auth := strings.TrimSpace(ctx.Get("Authorization"))

		if auth == "" {
			return dto.ResponseFailed(ctx, fiber.StatusUnauthorized, "authorization token not found", fiber.ErrUnauthorized)
		}
		if !strings.HasPrefix(auth, "Bearer") {
			return dto.ResponseFailed(ctx, fiber.StatusUnauthorized, "authentication schema not found", fiber.ErrUnauthorized)
		}

		rawToken := strings.Split(auth, " ")
		if len(rawToken) != 2 {
			return dto.ResponseFailed(ctx, fiber.StatusUnauthorized, "invalid authentication token", fiber.ErrUnauthorized)
		}

		token := rawToken[1]
		username, err := tg.ParseTokenUser(token)
		if err != nil {
			return dto.ResponseFailed(ctx, fiber.StatusUnauthorized, err.Error(), fiber.ErrUnauthorized)
		}

		fiber.StoreInContext(ctx, "username", username)
		return ctx.Next()
	}
}
