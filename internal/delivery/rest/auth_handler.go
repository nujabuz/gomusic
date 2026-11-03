package rest

import (
	"github.com/gofiber/fiber/v3"
	"github.com/nujabuz/gomusic/internal/entity"
	"github.com/nujabuz/gomusic/internal/entity/port"
)

type RequestUserRegister struct {
	Username string `json:"username" validate:"required,min=3,max=16"`
	Password string `json:"password" validate:"required,min=8,max=24"`
}

type RequestUserLogin struct {
	Username string `json:"username" validate:"required,min=3,max=16"`
	Password string `json:"password" validate:"required,min=8,max=24"`
}

type AuthHandler interface {
	Routes(api fiber.Router)
}

type authHandler struct {
	authUseCase port.AuthUseCase
}

func NewAuthHandler(authUseCase port.AuthUseCase) AuthHandler {
	return &authHandler{
		authUseCase: authUseCase,
	}
}

func (h *authHandler) Routes(api fiber.Router) {
	auth := api.Group("/auth")

	auth.Post("/register", h.register)
	auth.Post("/login", h.login)
}

func (h *authHandler) register(ctx fiber.Ctx) error {
	request := new(RequestUserRegister)
	if err := ctx.Bind().Body(request); err != nil {
		return ResponseFailedParseRequest(ctx, err)
	}

	out, err := h.authUseCase.Register(ctx, &port.RegisterInput{Username: request.Username, Password: request.Password})
	if err != nil {
		return EntityErrorTranslator(ctx, err)
	}

	res := new(ResponseUser)
	res.Username = out.Username

	return ResponseSuccess(ctx, fiber.StatusCreated, entity.SuccessRegister, res)
}

func (h *authHandler) login(ctx fiber.Ctx) error {
	request := new(RequestUserLogin)
	if err := ctx.Bind().Body(request); err != nil {
		return ResponseFailedParseRequest(ctx, err)
	}

	out, err := h.authUseCase.Login(ctx, &port.LoginInput{Username: request.Username, Password: request.Password})
	if err != nil {
		return EntityErrorTranslator(ctx, err)
	}

	res := new(ResponseUser)
	res.Token = out.Token

	return ResponseSuccess(ctx, fiber.StatusCreated, entity.SuccessRegister, res)
}
