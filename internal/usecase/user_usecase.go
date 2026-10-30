package usecase

import (
	"context"
	"errors"
	"fmt"

	"github.com/nujabuz/music-cloud/internal/domain"
)

type userUseCase struct {
	repo   domain.UserRepository
	hasher domain.PasswordHasher
}

func NewUserUseCase(repo domain.UserRepository, hasher domain.PasswordHasher) domain.UserUseCase {
	return &userUseCase{
		repo:   repo,
		hasher: hasher,
	}
}

func (uc *userUseCase) Register(ctx context.Context, in *domain.RegisterInput) (*domain.UserOutput, error) {

	_, err := uc.repo.FindByUsername(ctx, in.Username)
	if err != nil {
		if !errors.Is(err, domain.ErrUserNotFound) {
			return nil, fmt.Errorf("register.find: look for username: %w", err)
		}
	} else {
		return nil, domain.ErrUserAlreadyExists
	}

	hashed, err := uc.hasher.Hash(in.Password)
	if err != nil {
		if errors.Is(err, domain.ErrInvalidUser) {
			return nil, err
		}
		return nil, fmt.Errorf("register.hash: hash password: %w", err)
	}

	userRegister := new(domain.User)
	userRegister.Username = in.Username
	userRegister.Password = hashed

	user, err := uc.repo.Create(ctx, userRegister)
	if err != nil {
		return nil, fmt.Errorf("register.create: create user: %w", err)
	}

	out := new(domain.UserOutput)
	out.ID = user.ID
	out.Username = user.Username
	out.CreatedAt = user.CreatedAt
	out.UpdatedAt = user.UpdatedAt

	return out, nil
}

func (uc *userUseCase) Login(ctx context.Context, in *domain.LoginInput) (*domain.UserOutput, error) {

	user, err := uc.repo.FindByUsername(ctx, in.Username)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			return nil, domain.ErrInvalidUser
		}
		return nil, fmt.Errorf("login.find: look for username: %w", err)
	}

	if err := uc.hasher.Compare(in.Password, user.Password); err != nil {
		if errors.Is(err, domain.ErrInvalidUser) {
			return nil, err
		}
		return nil, fmt.Errorf("login.compare: compare password: %w", err)
	}

	out := new(domain.UserOutput)
	out.ID = user.ID
	out.Username = user.Username
	out.CreatedAt = user.CreatedAt
	out.UpdatedAt = user.UpdatedAt

	return out, nil
}

func (uc *userUseCase) ChangePassword(ctx context.Context, in *domain.ChangePasswordInput) (*domain.UserOutput, error) {
	user, err := uc.repo.FindByUsername(ctx, in.Username)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			return nil, domain.ErrInvalidUser
		}
		return nil, fmt.Errorf("changePassword.find: look for username: %w", err)
	}

	if err := uc.hasher.Compare(in.PreviousPassword, user.Password); err != nil {
		if errors.Is(err, domain.ErrInvalidUser) {
			return nil, err
		}
		return nil, fmt.Errorf("changePassword.compare: compare password: %w", err)
	}

	hashed, err := uc.hasher.Hash(in.NewPassword)
	if err != nil {
		if errors.Is(err, domain.ErrInvalidUser) {
			return nil, err
		}
		return nil, fmt.Errorf("changePassword.hash: hash password: %w", err)
	}

	user.Password = hashed

	updatedUser, err := uc.repo.Update(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("changePassword.update: update password: %w", err)
	}

	out := new(domain.UserOutput)
	out.ID = updatedUser.ID
	out.Username = updatedUser.Username
	out.UpdatedAt = updatedUser.UpdatedAt
	out.CreatedAt = updatedUser.CreatedAt

	return out, nil
}
