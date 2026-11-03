package port

import (
	"context"
	"time"

	"github.com/nujabuz/gomusic/internal/entity"
)

type UserRepository interface {
	Create(ctx context.Context, user *entity.User) error
	Find(ctx context.Context, user *entity.User) (*entity.User, error)
}

type UserUseCase interface {
	FindUserByUsername(ctx context.Context, username string) (*UserOutput, error)
}

type UserOutput struct {
	ID        int64
	Username  string
	CreatedAt time.Time
	UpdatedAt time.Time
}
