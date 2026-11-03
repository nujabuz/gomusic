// Package port
package port

import "context"

type AuthUseCase interface {
	Register(ctx context.Context, in *RegisterInput) (*RegisterOutput, error)
	Login(ctx context.Context, in *LoginInput) (*LoginOutput, error)
}

type RegisterInput struct {
	Username string
	Password string
}

type RegisterOutput struct {
	Username string
}

type LoginInput struct {
	Username string
	Password string
}

type LoginOutput struct {
	Token string
}
