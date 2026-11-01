package usecase

import (
	"context"
	"errors"

	"github.com/nujabuz/gomusic/internal/entity"
)

// Ports For UserUseCase
type PasswordHasher interface {
	Hash(plain string) (string, error)
}

type userUseCase struct {
	repo   entity.UserRepository
	hasher PasswordHasher
}

func NewUserUseCase(repo entity.UserRepository, hasher PasswordHasher) entity.UserUseCase {
	return &userUseCase{
		repo:   repo,
		hasher: hasher,
	}
}

func (u *userUseCase) Register(ctx context.Context, in *entity.RegisterInput) (*entity.UserOutput, error) {
	if _, err := u.repo.Find(ctx, &entity.User{Username: in.Username}); err != nil {
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

	if err := u.repo.Create(ctx, user); err != nil {
		return nil, err
	}

	out := u.toOutput(user)
	return out, nil
}

func (u *userUseCase) FindUserByUsername(ctx context.Context, username string) (*entity.UserOutput, error) {
	user, err := u.repo.Find(ctx, &entity.User{Username: username})
	if err != nil {
		return nil, err
	}

	out := u.toOutput(user)
	return out, nil
}

// Helper
func (u *userUseCase) toOutput(user *entity.User) *entity.UserOutput {
	return &entity.UserOutput{ID: user.ID, Username: user.Username, CreatedAt: user.CreatedAt, UpdatedAt: user.UpdatedAt}
}
