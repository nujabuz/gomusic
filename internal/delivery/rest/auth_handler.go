package rest

import (
	"github.com/gofiber/fiber/v3"
	"github.com/nujabuz/gomusic/internal/entity"
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
	authUseCase entity.AuthUseCase
}

func NewAuthHandler(authUseCase entity.AuthUseCase) AuthHandler {
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

	out, err := h.authUseCase.Register(ctx, &entity.RegisterInput{Username: request.Username, Password: request.Password})
	if err != nil {
		return EntityErrorTranslator(ctx, err)
	}

	res := new(ResponseUser)
	res.Username = out.Username
	res.CreatedAt = out.CreatedAt

	return ResponseSuccess(ctx, fiber.StatusCreated, entity.SuccessRegister, res)
}

func (h *authHandler) login(ctx fiber.Ctx) error {
	request := new(RequestUserLogin)
	if err := ctx.Bind().Body(request); err != nil {
		return ResponseFailedParseRequest(ctx, err)
	}

	token, err := h.authUseCase.Login(ctx, &entity.LoginInput{Username: request.Username, Password: request.Password})
	if err != nil {
		return EntityErrorTranslator(ctx, err)
	}

	res := new(ResponseUser)
	res.Token = token

	return ResponseSuccess(ctx, fiber.StatusCreated, entity.SuccessRegister, res)
}
