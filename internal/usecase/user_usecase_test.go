package usecase_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/nujabuz/music-cloud/internal/domain"
	"github.com/nujabuz/music-cloud/internal/usecase"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

// MockUserRepository mocks domain.UserRepository
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Create(ctx context.Context, user *domain.User) (*domain.User, error) {
	args := m.Called(ctx, user)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserRepository) Update(ctx context.Context, user *domain.User) (*domain.User, error) {
	args := m.Called(ctx, user)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserRepository) FindByUsername(ctx context.Context, username string) (*domain.User, error) {
	args := m.Called(ctx, username)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

// MockPasswordHasher mocks domain.PasswordHasher
type MockPasswordHasher struct {
	mock.Mock
}

func (m *MockPasswordHasher) Hash(password string) (string, error) {
	args := m.Called(password)
	return args.String(0), args.Error(1)
}

func (m *MockPasswordHasher) Compare(password, hash string) error {
	args := m.Called(password, hash)
	return args.Error(0)
}

type UserUseCaseTestSuite struct {
	suite.Suite
	mockRepo   *MockUserRepository
	mockHasher *MockPasswordHasher
	useCase    domain.UserUseCase
}

// SetupTest runs before each test case
func (s *UserUseCaseTestSuite) SetupTest() {
	s.mockRepo = new(MockUserRepository)
	s.mockHasher = new(MockPasswordHasher)
	s.useCase = usecase.NewUserUseCase(s.mockRepo, s.mockHasher)
}

// TestRegister_Success tests successful registration
func (s *UserUseCaseTestSuite) TestRegister_Success() {
	in := &domain.RegisterInput{
		Username: "newuser",
		Password: "plainpassword",
	}
	hashedPassword := "hashedpassword"
	now := time.Now()

	s.mockRepo.On("FindByUsername", mock.Anything, in.Username).
		Return(nil, domain.ErrUserNotFound)
	s.mockHasher.On("Hash", in.Password).
		Return(hashedPassword, nil)
	s.mockRepo.On("Create", mock.Anything, &domain.User{
		Username: in.Username,
		Password: hashedPassword,
	}).Return(&domain.User{
		ID:        1,
		Username:  in.Username,
		Password:  hashedPassword,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil)

	out, err := s.useCase.Register(context.Background(), in)

	s.NoError(err)
	s.NotNil(out)
	s.Equal(int64(1), out.ID)
	s.Equal(in.Username, out.Username)
	s.mockRepo.AssertExpectations(s.T())
	s.mockHasher.AssertExpectations(s.T())
}

// TestRegister_UsernameExists tests failure when username already exists
func (s *UserUseCaseTestSuite) TestRegister_UsernameExists() {
	in := &domain.RegisterInput{
		Username: "existinguser",
		Password: "plainpassword",
	}

	s.mockRepo.On("FindByUsername", mock.Anything, in.Username).
		Return(&domain.User{Username: in.Username}, nil)

	out, err := s.useCase.Register(context.Background(), in)

	s.ErrorIs(err, domain.ErrUserAlreadyExists)
	s.Nil(out)
	s.mockHasher.AssertNotCalled(s.T(), "Hash")
	s.mockRepo.AssertNotCalled(s.T(), "Create")
}

// TestLogin_Success tests successful login
func (s *UserUseCaseTestSuite) TestLogin_Success() {
	in := &domain.LoginInput{
		Username: "validuser",
		Password: "plainpassword",
	}
	user := &domain.User{
		ID:        5,
		Username:  in.Username,
		Password:  "hashedpassword",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	s.mockRepo.On("FindByUsername", mock.Anything, in.Username).
		Return(user, nil)
	s.mockHasher.On("Compare", in.Password, user.Password).
		Return(nil)

	out, err := s.useCase.Login(context.Background(), in)

	s.NoError(err)
	s.NotNil(out)
	s.Equal(user.ID, out.ID)
	s.Equal(user.Username, out.Username)
	s.mockRepo.AssertExpectations(s.T())
	s.mockHasher.AssertExpectations(s.T())
}

// TestLogin_UserNotFound tests login with non-existent user
func (s *UserUseCaseTestSuite) TestLogin_UserNotFound() {
	in := &domain.LoginInput{
		Username: "unknownuser",
		Password: "plainpassword",
	}

	s.mockRepo.On("FindByUsername", mock.Anything, in.Username).
		Return(nil, domain.ErrUserNotFound)

	out, err := s.useCase.Login(context.Background(), in)

	s.ErrorIs(err, domain.ErrInvalidUser)
	s.Nil(out)
	s.mockHasher.AssertNotCalled(s.T(), "Compare")
}

// TestLogin_WrongPassword tests login with incorrect password
func (s *UserUseCaseTestSuite) TestLogin_WrongPassword() {
	in := &domain.LoginInput{
		Username: "validuser",
		Password: "wrongpassword",
	}
	user := &domain.User{
		ID:       5,
		Username: in.Username,
		Password: "hashedpassword",
	}

	s.mockRepo.On("FindByUsername", mock.Anything, in.Username).
		Return(user, nil)
	s.mockHasher.On("Compare", in.Password, user.Password).
		Return(domain.ErrInvalidUser)

	out, err := s.useCase.Login(context.Background(), in)

	s.ErrorIs(err, domain.ErrInvalidUser)
	s.Nil(out)
	s.mockRepo.AssertExpectations(s.T())
	s.mockHasher.AssertExpectations(s.T())
}

// TestChangePassword_Success tests successful password change
func (s *UserUseCaseTestSuite) TestChangePassword_Success() {
	in := &domain.ChangePasswordInput{
		Username:         "johndoe",
		PreviousPassword: "oldpassword",
		NewPassword:      "newpassword",
	}
	oldHashed := "oldhashed"
	newHashed := "newhashed"
	now := time.Now()

	user := &domain.User{
		ID:        1,
		Username:  in.Username,
		Password:  oldHashed,
		CreatedAt: now,
		UpdatedAt: now,
	}

	updatedUser := &domain.User{
		ID:        1,
		Username:  in.Username,
		Password:  newHashed,
		CreatedAt: now,
		UpdatedAt: now,
	}

	s.mockRepo.On("FindByUsername", mock.Anything, in.Username).
		Return(user, nil)
	s.mockHasher.On("Compare", in.PreviousPassword, oldHashed).
		Return(nil)
	s.mockHasher.On("Hash", in.NewPassword).
		Return(newHashed, nil)
	s.mockRepo.On("Update", mock.Anything, mock.MatchedBy(func(u *domain.User) bool {
		return u.Username == in.Username && u.Password == newHashed
	})).Return(updatedUser, nil)

	out, err := s.useCase.ChangePassword(context.Background(), in)

	s.NoError(err)
	s.NotNil(out)
	s.Equal(updatedUser.ID, out.ID)
	s.Equal(updatedUser.Username, out.Username)
	s.mockRepo.AssertExpectations(s.T())
	s.mockHasher.AssertExpectations(s.T())
}

// TestChangePassword_WrongPreviousPassword tests password change with wrong old password
func (s *UserUseCaseTestSuite) TestChangePassword_WrongPreviousPassword() {
	in := &domain.ChangePasswordInput{
		Username:         "johndoe",
		PreviousPassword: "wrongpassword",
		NewPassword:      "newpassword",
	}
	user := &domain.User{
		Username: in.Username,
		Password: "currentpasswordhash",
	}

	s.mockRepo.On("FindByUsername", mock.Anything, in.Username).
		Return(user, nil)
	s.mockHasher.On("Compare", in.PreviousPassword, user.Password).
		Return(errors.New("password mismatch"))

	out, err := s.useCase.ChangePassword(context.Background(), in)

	s.Error(err)
	s.Nil(out)
	s.mockHasher.AssertNotCalled(s.T(), "Hash")
	s.mockRepo.AssertNotCalled(s.T(), "Update")
}

func TestUserUseCaseTestSuite(t *testing.T) {
	suite.Run(t, new(UserUseCaseTestSuite))
}
