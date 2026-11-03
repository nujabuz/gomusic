package rest

import (
	"strings"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/nujabuz/gomusic/internal/entity"
	"github.com/nujabuz/gomusic/internal/entity/port"
)

type ResponseUser struct {
	ID        int64     `json:"id,omitempty"`
	Username  string    `json:"username,omitempty"`
	Token     string    `json:"token,omitempty"`
	CreatedAt time.Time `json:"created_at,omitzero"`
	UpdatedAt time.Time `json:"updated_at,omitzero"`
}
type UserHandler struct {
	UserUseCase port.UserUseCase
	Token       port.TokenGenerator
}

func NewUserHandler(userUseCase port.UserUseCase, token port.TokenGenerator) *UserHandler {
	return &UserHandler{
		UserUseCase: userUseCase,
		Token:       token,
	}
}

func (h *UserHandler) Routes(api fiber.Router) {
	users := api.Group("/users")

	users.Use(h.tokenMiddleware)

	users.Get("/:username<betweenLen(3,16)>", h.findByUsername)
}

func (h *UserHandler) findByUsername(ctx fiber.Ctx) error {
	username := fiber.Params[string](ctx, "username")

	out, err := h.UserUseCase.FindUserByUsername(ctx, username)
	if err != nil {
		return EntityErrorTranslator(ctx, err)
	}

	res := new(ResponseUser)
	res.ID = out.ID
	res.Username = out.Username
	res.CreatedAt = out.CreatedAt
	res.UpdatedAt = out.UpdatedAt

	return ResponseSuccess(ctx, fiber.StatusFound, entity.SuccessFindbyUsername, res)
}

// Middleware
func (h *UserHandler) tokenMiddleware(ctx fiber.Ctx) error {
	authHeader := ctx.Get("Authorization")
	if !strings.HasPrefix(authHeader, "Bearer") {
		return ResponseFailed(ctx, fiber.StatusUnauthorized, entity.ErrUnauthorized.Error(), fiber.ErrUnauthorized)
	}

	rawToken := strings.Split(authHeader, " ")
	if len(rawToken) != 2 {
		return ResponseFailed(ctx, fiber.StatusUnauthorized, entity.ErrUnauthorized.Error(), fiber.ErrUnauthorized)
	}

	token := rawToken[1]
	claims, err := h.Token.Parse(token)
	if err != nil {
		EntityErrorTranslator(ctx, err)
	}

	ctx.SetLocal("userId", claims.ID)
	ctx.SetLocal("username", claims.Subject)

	return ctx.Next()
}
