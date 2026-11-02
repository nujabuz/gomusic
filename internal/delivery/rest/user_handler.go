package rest

import (
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/nujabuz/gomusic/internal/entity"
)

type ResponseUser struct {
	ID        int64     `json:"id,omitempty"`
	Username  string    `json:"username,omitempty"`
	Token     string    `json:"token,omitempty"`
	CreatedAt time.Time `json:"created_at,omitzero"`
	UpdatedAt time.Time `json:"updated_at,omitzero"`
}
type UserHandler struct {
	UserUseCase entity.UserUseCase
}

func NewUserHandler(userUseCase entity.UserUseCase) *UserHandler {
	return &UserHandler{
		UserUseCase: userUseCase,
	}
}

func (h *UserHandler) Routes(api fiber.Router) {
	users := api.Group("/users")
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
