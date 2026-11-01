package rest

import (
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/nujabuz/gomusic/internal/entity"
)

type RequestUserRegister struct {
	Username string `json:"username" validate:"required,min=3,max=16"`
	Password string `json:"password" validate:"required,min=8,max=24"`
}

type ResponseUser struct {
	ID        int64     `json:"id,omitempty"`
	Username  string    `json:"username,omitempty"`
	CreatedAt time.Time `json:"created_at,omitzero"`
	UpdatedAt time.Time `json:"updated_at,omitzero"`
}

type UserHandler struct {
	UserUseCase entity.UserUseCase
}

func NewUserHandler(
	userUseCase entity.UserUseCase,
) *UserHandler {
	return &UserHandler{
		UserUseCase: userUseCase,
	}
}

func (h *UserHandler) Routes(api fiber.Router) {
	auth := api.Group("/auth")
	auth.Post("/register", h.register)

	users := api.Group("/users")
	users.Get("/:username<betweenLen(3,16)>", h.findByUsername)
}

func (h *UserHandler) register(ctx fiber.Ctx) error {
	request := new(RequestUserRegister)
	if err := ctx.Bind().Body(request); err != nil {
		return ResponseFailedParseRequest(ctx, err)
	}

	out, err := h.UserUseCase.Register(ctx, &entity.RegisterInput{Username: request.Username, Password: request.Password})
	if err != nil {
		return EntityErrorTranslator(ctx, err)
	}

	res := new(ResponseUser)
	res.Username = out.Username
	res.CreatedAt = out.CreatedAt

	return ResponseSuccess(ctx, fiber.StatusCreated, entity.SuccessRegister, res)
}

// helper
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
