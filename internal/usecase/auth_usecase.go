package usecase

import (
	"context"
	"errors"

	"github.com/nujabuz/gomusic/internal/entity"
	"github.com/nujabuz/gomusic/internal/entity/port"
)

type authUseCase struct {
	hasher         port.PasswordHasher
	token          port.TokenGenerator
	userRepository port.UserRepository
	userUseCase    port.UserUseCase
}

func NewAuthUseCase(hasher port.PasswordHasher, token port.TokenGenerator, userRepository port.UserRepository, userUseCase port.UserUseCase) port.AuthUseCase {
	return &authUseCase{
		hasher:         hasher,
		token:          token,
		userRepository: userRepository,
		userUseCase:    userUseCase,
	}
}

func (u *authUseCase) Register(ctx context.Context, in *port.RegisterInput) (*port.RegisterOutput, error) {
	if _, err := u.userRepository.Find(ctx, &entity.User{Username: in.Username}); err != nil {
		if !errors.Is(err, entity.ErrUserNotFound) {
			return nil, err
		}
	}

	hashed, err := u.hasher.Hash(in.Password)
	if err != nil {
		return nil, err
	}

	user := new(entity.User)
	user.Username = in.Username
	user.Password = hashed

	if err := u.userRepository.Create(ctx, user); err != nil {
		return nil, err
	}

	out := new(port.RegisterOutput)
	out.Username = user.Username

	return out, nil
}

func (u *authUseCase) Login(ctx context.Context, in *port.LoginInput) (*port.LoginOutput, error) {
	user, err := u.userRepository.Find(ctx, &entity.User{Username: in.Username})
	if err != nil {
		return nil, err
	}

	if err := u.hasher.Compare(user.Password, in.Password); err != nil {
		return nil, err
	}

	token, err := u.token.Generate(user)
	if err != nil {
		return nil, err
	}

	out := new(port.LoginOutput)
	out.Token = token

	return out, nil
}
