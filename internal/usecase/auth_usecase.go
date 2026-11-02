package usecase

import (
	"context"
	"errors"

	"github.com/nujabuz/gomusic/internal/entity"
)

// PasswordHasher port
type PasswordHasher interface {
	Hash(plain string) (string, error)
	Compare(hashedPassword string, password string) error
}

// TokenGenerator port
type TokenGenerator interface {
	Generate(user *entity.User) (string, error)
}

type authUseCase struct {
	hasher         PasswordHasher
	token          TokenGenerator
	userRepository entity.UserRepository
	userUseCase    entity.UserUseCase
}

func NewAuthUseCase(hasher PasswordHasher, token TokenGenerator, userRepository entity.UserRepository, userUseCase entity.UserUseCase) entity.AuthUseCase {
	return &authUseCase{
		hasher:         hasher,
		token:          token,
		userRepository: userRepository,
		userUseCase:    userUseCase,
	}
}

func (u *authUseCase) Register(ctx context.Context, in *entity.RegisterInput) (*entity.UserOutput, error) {
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

	out := u.userUseCase.ToUserOutput(user)
	return out, nil
}

func (u *authUseCase) Login(ctx context.Context, in *entity.LoginInput) (string, error) {
	user, err := u.userRepository.Find(ctx, &entity.User{Username: in.Username})
	if err != nil {
		return "", err
	}

	if err := u.hasher.Compare(user.Password, in.Password); err != nil {
		return "", err
	}

	token, err := u.token.Generate(user)
	if err != nil {
		return "", err
	}

	return token, nil
}
