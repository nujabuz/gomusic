package usecase_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/nujabuz/music-cloud/internal/domain"
	"github.com/nujabuz/music-cloud/internal/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockUserRepository implements domain.UserRepository
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Create(ctx context.Context, user *domain.User) (*domain.User, error) {
	args := m.Called(ctx, user)
	if res := args.Get(0); res != nil {
		return res.(*domain.User), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockUserRepository) Update(ctx context.Context, user *domain.User) (*domain.User, error) {
	args := m.Called(ctx, user)
	if res := args.Get(0); res != nil {
		return res.(*domain.User), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockUserRepository) Delete() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockUserRepository) FindByUsername(ctx context.Context, username string) (*domain.User, error) {
	args := m.Called(ctx, username)
	if res := args.Get(0); res != nil {
		return res.(*domain.User), args.Error(1)
	}
	return nil, args.Error(1)
}

// MockPasswordHasher implements domain.PasswordHasher
type MockPasswordHasher struct {
	mock.Mock
}

func (m *MockPasswordHasher) Hash(plain string) (string, error) {
	args := m.Called(plain)
	return args.String(0), args.Error(1)
}

func (m *MockPasswordHasher) Compare(plain string, hashed string) error {
	args := m.Called(plain, hashed)
	return args.Error(0)
}

func TestUserUseCase_Register(t *testing.T) {
	ctx := context.Background()
	now := time.Now()

	t.Run("success", func(t *testing.T) {
		repo := new(MockUserRepository)
		hasher := new(MockPasswordHasher)
		uc := usecase.NewUserUseCase(repo, hasher)

		repo.On("FindByUsername", ctx, "john").Return(nil, domain.ErrUserNotFound)
		hasher.On("Hash", "secret123").Return("hashed_secret", nil)
		repo.On("Create", ctx, &domain.User{
			Username: "john",
			Password: "hashed_secret",
		}).Return(&domain.User{
			ID:        1,
			Username:  "john",
			Password:  "hashed_secret",
			CreatedAt: now,
			UpdatedAt: now,
		}, nil)

		res, err := uc.Register(ctx, &domain.RegisterInput{
			Username: "john",
			Password: "secret123",
		})

		assert.NoError(t, err)
		assert.Equal(t, &domain.UserOutput{
			ID:        1,
			Username:  "john",
			CreatedAt: now,
			UpdatedAt: now,
		}, res)
		repo.AssertExpectations(t)
		hasher.AssertExpectations(t)
	})

	t.Run("username already exists", func(t *testing.T) {
		repo := new(MockUserRepository)
		hasher := new(MockPasswordHasher)
		uc := usecase.NewUserUseCase(repo, hasher)

		repo.On("FindByUsername", ctx, "john").Return(&domain.User{ID: 1, Username: "john"}, nil)

		res, err := uc.Register(ctx, &domain.RegisterInput{
			Username: "john",
			Password: "secret123",
		})

		assert.Nil(t, res)
		assert.ErrorIs(t, err, domain.ErrUserAlreadyExists)
		repo.AssertExpectations(t)
	})

	t.Run("repo find error", func(t *testing.T) {
		repo := new(MockUserRepository)
		hasher := new(MockPasswordHasher)
		uc := usecase.NewUserUseCase(repo, hasher)

		dbErr := errors.New("db connection failure")
		repo.On("FindByUsername", ctx, "john").Return(nil, dbErr)

		res, err := uc.Register(ctx, &domain.RegisterInput{
			Username: "john",
			Password: "secret123",
		})

		assert.Nil(t, res)
		assert.ErrorIs(t, err, dbErr)
		repo.AssertExpectations(t)
	})

	t.Run("hasher invalid user error", func(t *testing.T) {
		repo := new(MockUserRepository)
		hasher := new(MockPasswordHasher)
		uc := usecase.NewUserUseCase(repo, hasher)

		repo.On("FindByUsername", ctx, "john").Return(nil, domain.ErrUserNotFound)
		hasher.On("Hash", "short").Return("", domain.ErrInvalidUser)

		res, err := uc.Register(ctx, &domain.RegisterInput{
			Username: "john",
			Password: "short",
		})

		assert.Nil(t, res)
		assert.ErrorIs(t, err, domain.ErrInvalidUser)
		repo.AssertExpectations(t)
		hasher.AssertExpectations(t)
	})

	t.Run("hasher internal error", func(t *testing.T) {
		repo := new(MockUserRepository)
		hasher := new(MockPasswordHasher)
		uc := usecase.NewUserUseCase(repo, hasher)

		hashErr := errors.New("hash calculation failed")
		repo.On("FindByUsername", ctx, "john").Return(nil, domain.ErrUserNotFound)
		hasher.On("Hash", "secret123").Return("", hashErr)

		res, err := uc.Register(ctx, &domain.RegisterInput{
			Username: "john",
			Password: "secret123",
		})

		assert.Nil(t, res)
		assert.ErrorIs(t, err, hashErr)
		repo.AssertExpectations(t)
		hasher.AssertExpectations(t)
	})

	t.Run("repo create error", func(t *testing.T) {
		repo := new(MockUserRepository)
		hasher := new(MockPasswordHasher)
		uc := usecase.NewUserUseCase(repo, hasher)

		dbErr := errors.New("insert failed")
		repo.On("FindByUsername", ctx, "john").Return(nil, domain.ErrUserNotFound)
		hasher.On("Hash", "secret123").Return("hashed_secret", nil)
		repo.On("Create", ctx, &domain.User{
			Username: "john",
			Password: "hashed_secret",
		}).Return(nil, dbErr)

		res, err := uc.Register(ctx, &domain.RegisterInput{
			Username: "john",
			Password: "secret123",
		})

		assert.Nil(t, res)
		assert.ErrorIs(t, err, dbErr)
		repo.AssertExpectations(t)
		hasher.AssertExpectations(t)
	})
}

func TestUserUseCase_Login(t *testing.T) {
	ctx := context.Background()
	now := time.Now()

	t.Run("success", func(t *testing.T) {
		repo := new(MockUserRepository)
		hasher := new(MockPasswordHasher)
		uc := usecase.NewUserUseCase(repo, hasher)

		user := &domain.User{
			ID:        1,
			Username:  "john",
			Password:  "hashed_secret",
			CreatedAt: now,
			UpdatedAt: now,
		}

		repo.On("FindByUsername", ctx, "john").Return(user, nil)
		hasher.On("Compare", "secret123", "hashed_secret").Return(nil)

		res, err := uc.Login(ctx, &domain.LoginInput{
			Username: "john",
			Password: "secret123",
		})

		assert.NoError(t, err)
		assert.Equal(t, &domain.UserOutput{
			ID:        1,
			Username:  "john",
			CreatedAt: now,
			UpdatedAt: now,
		}, res)
		repo.AssertExpectations(t)
		hasher.AssertExpectations(t)
	})

	t.Run("user not found returns ErrInvalidUser", func(t *testing.T) {
		repo := new(MockUserRepository)
		hasher := new(MockPasswordHasher)
		uc := usecase.NewUserUseCase(repo, hasher)

		repo.On("FindByUsername", ctx, "unknown").Return(nil, domain.ErrUserNotFound)

		res, err := uc.Login(ctx, &domain.LoginInput{
			Username: "unknown",
			Password: "secret123",
		})

		assert.Nil(t, res)
		assert.ErrorIs(t, err, domain.ErrInvalidUser)
		repo.AssertExpectations(t)
	})

	t.Run("repo find error", func(t *testing.T) {
		repo := new(MockUserRepository)
		hasher := new(MockPasswordHasher)
		uc := usecase.NewUserUseCase(repo, hasher)

		dbErr := errors.New("db error")
		repo.On("FindByUsername", ctx, "john").Return(nil, dbErr)

		res, err := uc.Login(ctx, &domain.LoginInput{
			Username: "john",
			Password: "secret123",
		})

		assert.Nil(t, res)
		assert.ErrorIs(t, err, dbErr)
		repo.AssertExpectations(t)
	})

	t.Run("wrong password returns ErrInvalidUser", func(t *testing.T) {
		repo := new(MockUserRepository)
		hasher := new(MockPasswordHasher)
		uc := usecase.NewUserUseCase(repo, hasher)

		user := &domain.User{
			ID:       1,
			Username: "john",
			Password: "hashed_secret",
		}
		repo.On("FindByUsername", ctx, "john").Return(user, nil)
		hasher.On("Compare", "wrongpass", "hashed_secret").Return(domain.ErrInvalidUser)

		res, err := uc.Login(ctx, &domain.LoginInput{
			Username: "john",
			Password: "wrongpass",
		})

		assert.Nil(t, res)
		assert.ErrorIs(t, err, domain.ErrInvalidUser)
		repo.AssertExpectations(t)
		hasher.AssertExpectations(t)
	})

	t.Run("compare internal error", func(t *testing.T) {
		repo := new(MockUserRepository)
		hasher := new(MockPasswordHasher)
		uc := usecase.NewUserUseCase(repo, hasher)

		user := &domain.User{
			ID:       1,
			Username: "john",
			Password: "hashed_secret",
		}
		compareErr := errors.New("bcrypt failure")
		repo.On("FindByUsername", ctx, "john").Return(user, nil)
		hasher.On("Compare", "secret123", "hashed_secret").Return(compareErr)

		res, err := uc.Login(ctx, &domain.LoginInput{
			Username: "john",
			Password: "secret123",
		})

		assert.Nil(t, res)
		assert.ErrorIs(t, err, compareErr)
		repo.AssertExpectations(t)
		hasher.AssertExpectations(t)
	})
}

func TestUserUseCase_ChangePassword(t *testing.T) {
	ctx := context.Background()
	now := time.Now()

	existingUser := &domain.User{
		ID:        1,
		Username:  "john",
		Password:  "old_hashed",
		CreatedAt: now,
		UpdatedAt: now,
	}

	t.Run("success", func(t *testing.T) {
		repo := new(MockUserRepository)
		hasher := new(MockPasswordHasher)
		uc := usecase.NewUserUseCase(repo, hasher)

		userCopy := *existingUser
		repo.On("FindByUsername", ctx, "john").Return(&userCopy, nil)
		hasher.On("Compare", "old_pass", "old_hashed").Return(nil)
		hasher.On("Hash", "new_pass").Return("new_hashed", nil)
		repo.On("Update", ctx, &userCopy).Return(&domain.User{
			ID:        1,
			Username:  "john",
			Password:  "new_hashed",
			CreatedAt: now,
			UpdatedAt: now.Add(time.Minute),
		}, nil)

		res, err := uc.ChangePassword(ctx, &domain.ChangePasswordInput{
			Username:         "john",
			PreviousPassword: "old_pass",
			NewPassword:      "new_pass",
		})

		assert.NoError(t, err)
		assert.Equal(t, &domain.UserOutput{
			ID:        1,
			Username:  "john",
			CreatedAt: now,
			UpdatedAt: now.Add(time.Minute),
		}, res)
		repo.AssertExpectations(t)
		hasher.AssertExpectations(t)
	})

	t.Run("user not found returns ErrInvalidUser", func(t *testing.T) {
		repo := new(MockUserRepository)
		hasher := new(MockPasswordHasher)
		uc := usecase.NewUserUseCase(repo, hasher)

		repo.On("FindByUsername", ctx, "unknown").Return(nil, domain.ErrUserNotFound)

		res, err := uc.ChangePassword(ctx, &domain.ChangePasswordInput{
			Username:         "unknown",
			PreviousPassword: "old_pass",
			NewPassword:      "new_pass",
		})

		assert.Nil(t, res)
		assert.ErrorIs(t, err, domain.ErrInvalidUser)
		repo.AssertExpectations(t)
	})

	t.Run("repo find error", func(t *testing.T) {
		repo := new(MockUserRepository)
		hasher := new(MockPasswordHasher)
		uc := usecase.NewUserUseCase(repo, hasher)

		dbErr := errors.New("db error")
		repo.On("FindByUsername", ctx, "john").Return(nil, dbErr)

		res, err := uc.ChangePassword(ctx, &domain.ChangePasswordInput{
			Username:         "john",
			PreviousPassword: "old_pass",
			NewPassword:      "new_pass",
		})

		assert.Nil(t, res)
		assert.ErrorIs(t, err, dbErr)
		repo.AssertExpectations(t)
	})

	t.Run("wrong previous password returns ErrInvalidUser", func(t *testing.T) {
		repo := new(MockUserRepository)
		hasher := new(MockPasswordHasher)
		uc := usecase.NewUserUseCase(repo, hasher)

		userCopy := *existingUser
		repo.On("FindByUsername", ctx, "john").Return(&userCopy, nil)
		hasher.On("Compare", "wrong_pass", "old_hashed").Return(domain.ErrInvalidUser)

		res, err := uc.ChangePassword(ctx, &domain.ChangePasswordInput{
			Username:         "john",
			PreviousPassword: "wrong_pass",
			NewPassword:      "new_pass",
		})

		assert.Nil(t, res)
		assert.ErrorIs(t, err, domain.ErrInvalidUser)
		repo.AssertExpectations(t)
		hasher.AssertExpectations(t)
	})

	t.Run("hash new password error", func(t *testing.T) {
		repo := new(MockUserRepository)
		hasher := new(MockPasswordHasher)
		uc := usecase.NewUserUseCase(repo, hasher)

		userCopy := *existingUser
		hashErr := errors.New("hash failed")
		repo.On("FindByUsername", ctx, "john").Return(&userCopy, nil)
		hasher.On("Compare", "old_pass", "old_hashed").Return(nil)
		hasher.On("Hash", "new_pass").Return("", hashErr)

		res, err := uc.ChangePassword(ctx, &domain.ChangePasswordInput{
			Username:         "john",
			PreviousPassword: "old_pass",
			NewPassword:      "new_pass",
		})

		assert.Nil(t, res)
		assert.ErrorIs(t, err, hashErr)
		repo.AssertExpectations(t)
		hasher.AssertExpectations(t)
	})

	t.Run("repo update error", func(t *testing.T) {
		repo := new(MockUserRepository)
		hasher := new(MockPasswordHasher)
		uc := usecase.NewUserUseCase(repo, hasher)

		userCopy := *existingUser
		updateErr := errors.New("update failed")
		repo.On("FindByUsername", ctx, "john").Return(&userCopy, nil)
		hasher.On("Compare", "old_pass", "old_hashed").Return(nil)
		hasher.On("Hash", "new_pass").Return("new_hashed", nil)
		repo.On("Update", ctx, &userCopy).Return(nil, updateErr)

		res, err := uc.ChangePassword(ctx, &domain.ChangePasswordInput{
			Username:         "john",
			PreviousPassword: "old_pass",
			NewPassword:      "new_pass",
		})

		assert.Nil(t, res)
		assert.ErrorIs(t, err, updateErr)
		repo.AssertExpectations(t)
		hasher.AssertExpectations(t)
	})
}
