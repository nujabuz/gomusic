package domain

import (
	"context"
	"errors"
	"time"
)

var (
	ErrUserNotFound      = errors.New("user not found")
	ErrInvalidUser       = errors.New("invalid username or password")
	ErrUserAlreadyExists = errors.New("username is already in use")
)

type User struct {
	ID        int64
	Username  string
	Password  string // hashed Password
	CreatedAt time.Time
	UpdatedAt time.Time
}

type UserUseCase interface {
	Register(ctx context.Context, in *RegisterInput) (*UserOutput, error)
	Login(ctx context.Context, in *LoginInput) (*UserOutput, error)
	ChangePassword(ctx context.Context, in *ChangePasswordInput) (*UserOutput, error)
}

type RegisterInput struct {
	Username string
	Password string
}

type LoginInput struct {
	Username string
	Password string
}

type ChangePasswordInput struct {
	Username         string
	PreviousPassword string
	NewPassword      string
}

type UserOutput struct {
	ID        int64
	Username  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type UserRepository interface {
	Create(ctx context.Context, user *User) (*User, error)
	Update(ctx context.Context, user *User) (*User, error)
	FindByUsername(ctx context.Context, username string) (*User, error)
}

type PasswordHasher interface {
	Hash(plain string) (string, error)
	Compare(plain string, hashed string) error
}
