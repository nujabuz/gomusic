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
	FindUserByUsername(ctx context.Context, username string) (*UserOutput, error)

	// helper
	ToUserOutput(user *User) *UserOutput
}

type AuthUseCase interface {
	Register(ctx context.Context, in *RegisterInput) (*UserOutput, error)
	Login(ctx context.Context, in *LoginInput) (string, error)
}

type RegisterInput struct {
	Username string
	Password string
}

type LoginInput struct {
	Username string
	Password string
}

type UserOutput struct {
	ID        int64
	Username  string
	CreatedAt time.Time
	UpdatedAt time.Time
}
