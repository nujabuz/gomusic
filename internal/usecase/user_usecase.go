// Package usecase
package usecase

import (
	"context"

	"github.com/nujabuz/gomusic/internal/entity"
	"github.com/nujabuz/gomusic/internal/entity/port"
)

type userUseCase struct {
	repo   port.UserRepository
	hasher port.PasswordHasher
}

func NewUserUseCase(repo port.UserRepository) port.UserUseCase {
	return &userUseCase{
		repo: repo,
	}
}

func (u *userUseCase) FindUserByUsername(ctx context.Context, username string) (*port.UserOutput, error) {
	user, err := u.repo.Find(ctx, &entity.User{Username: username})
	if err != nil {
		return nil, err
	}

	out := u.ToUserOutput(user)
	return out, nil
}

// Helper
func (u *userUseCase) ToUserOutput(user *entity.User) *port.UserOutput {
	return &port.UserOutput{ID: user.ID, Username: user.Username, CreatedAt: user.CreatedAt, UpdatedAt: user.UpdatedAt}
}
