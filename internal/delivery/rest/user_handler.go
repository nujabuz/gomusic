package rest

import (
	"github.com/gofiber/fiber/v3"
	"github.com/nujabuz/music-cloud/internal/delivery/rest/dto"
	"github.com/nujabuz/music-cloud/internal/domain"
)

type UserHandler interface {
	Register(ctx fiber.Ctx) error
	Login(ctx fiber.Ctx) error
}

type userHandler struct {
	UserUseCase domain.UserUseCase
}

func NewUserHandler(userUseCase domain.UserUseCase) UserHandler {
	return &userHandler{
		UserUseCase: userUseCase,
	}
}

func (h *userHandler) Register(ctx fiber.Ctx) error {
	request := new(dto.RequestRegister)
	if err := ctx.Bind().Body(request); err != nil {
		return dto.BindValidation(ctx, err)
	}

	in := new(domain.RegisterInput)
	in.Username = request.Username
	in.Password = request.Password

	out, err := h.UserUseCase.Register(ctx, in)
	if err != nil {
		return dto.TranslateDomainError(ctx, err)
	}

	res := new(dto.ResponseUser)
	res.Username = out.Username
	res.CreatedAt = out.CreatedAt

	return dto.ResponseSuccess(ctx, fiber.StatusOK, "user registration was successful", res)
}

func (h *userHandler) Login(ctx fiber.Ctx) error {
	request := new(dto.RequestLogin)
	if err := ctx.Bind().Body(request); err != nil {
		return dto.BindValidation(ctx, err)
	}

	in := new(domain.LoginInput)
	in.Username = request.Username
	in.Password = request.Password

	out, err := h.UserUseCase.Login(ctx, in)
	if err != nil {
		return dto.TranslateDomainError(ctx, err)
	}

	res := new(dto.ResponseUser)
	res.ID = out.ID
	res.Username = out.Username

	return dto.ResponseSuccess(ctx, fiber.StatusOK, "user login was successful", res)
}
