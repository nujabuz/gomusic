package entity

import (
	"context"
	"time"
)

type User struct {
	ID        int64
	Username  string
	Password  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type UserRepository interface {
	Create(ctx context.Context, user *User) error
	Find(ctx context.Context, user *User) (*User, error)
}

type UserUseCase interface {
	Register(ctx context.Context, in *RegisterInput) (*UserOutput, error)
	FindUserByUsername(ctx context.Context, username string) (*UserOutput, error)
}

type RegisterInput struct {
	Username string
	Password string
}

type UserOutput struct {
	ID        int64
	Username  string
	CreatedAt time.Time
	UpdatedAt time.Time
}
