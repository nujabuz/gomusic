package rest

import (
	"strings"

	"github.com/gofiber/fiber/v3"
	"github.com/nujabuz/gomusic/internal/entity"
	"github.com/nujabuz/gomusic/internal/entity/port"
)

type newMiddleware struct {
	token port.TokenGenerator
}

func NewMiddleware(tokenGenerator port.TokenGenerator) *newMiddleware {
	return &newMiddleware{
		token: tokenGenerator,
	}
}

func (m *newMiddleware) TokenMiddleware(ctx fiber.Ctx) error {
	authHeader := ctx.Get("Authorization")
	if !strings.HasPrefix(authHeader, "Bearer") {
		return ResponseFailed(ctx, fiber.StatusUnauthorized, entity.ErrUnauthorized.Error(), fiber.ErrUnauthorized)
	}

	rawToken := strings.Split(authHeader, " ")
	if len(rawToken) != 2 {
		return ResponseFailed(ctx, fiber.StatusUnauthorized, entity.ErrUnauthorized.Error(), fiber.ErrUnauthorized)
	}

	token := rawToken[1]
	claims, err := m.token.Parse(token)
	if err != nil {
		return EntityErrorTranslator(ctx, err)
	}

	ctx.SetLocal("userId", claims.ID)
	ctx.SetLocal("username", claims.Subject)

	return ctx.Next()
}
