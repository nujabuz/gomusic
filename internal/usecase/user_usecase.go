// Package usecase
package usecase

import (
	"context"

	"github.com/nujabuz/gomusic/internal/entity"
)

type userUseCase struct {
	repo   entity.UserRepository
	hasher PasswordHasher
}

func NewUserUseCase(repo entity.UserRepository) entity.UserUseCase {
	return &userUseCase{
		repo: repo,
	}
}

func (u *userUseCase) FindUserByUsername(ctx context.Context, username string) (*entity.UserOutput, error) {
	user, err := u.repo.Find(ctx, &entity.User{Username: username})
	if err != nil {
		return nil, err
	}

	out := u.ToUserOutput(user)
	return out, nil
}

// Helper
func (u *userUseCase) ToUserOutput(user *entity.User) *entity.UserOutput {
	return &entity.UserOutput{ID: user.ID, Username: user.Username, CreatedAt: user.CreatedAt, UpdatedAt: user.UpdatedAt}
}
