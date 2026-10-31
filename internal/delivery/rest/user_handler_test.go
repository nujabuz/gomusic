package rest_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
	"github.com/nujabuz/music-cloud/internal/delivery/rest"
	"github.com/nujabuz/music-cloud/internal/delivery/rest/dto"
	"github.com/nujabuz/music-cloud/internal/domain"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

// MockUserUseCase mocks domain.UserUseCase
type MockUserUseCase struct {
	mock.Mock
}

func (m *MockUserUseCase) Register(ctx context.Context, in *domain.RegisterInput) (*domain.UserOutput, error) {
	args := m.Called(ctx, in)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.UserOutput), args.Error(1)
}

func (m *MockUserUseCase) Login(ctx context.Context, in *domain.LoginInput) (*domain.UserOutput, error) {
	args := m.Called(ctx, in)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.UserOutput), args.Error(1)
}

func (m *MockUserUseCase) ChangePassword(ctx context.Context, in *domain.ChangePasswordInput) (*domain.UserOutput, error) {
	args := m.Called(ctx, in)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.UserOutput), args.Error(1)
}

type structValidator struct {
	validate *validator.Validate
}

func (v *structValidator) Validate(out any) error {
	return v.validate.Struct(out)
}

type UserHandlerTestSuite struct {
	suite.Suite
	mockUseCase *MockUserUseCase
	app         *fiber.App
}

// SetupTest runs before each test case
func (s *UserHandlerTestSuite) SetupTest() {
	s.mockUseCase = new(MockUserUseCase)
	handler := rest.NewUserHandler(s.mockUseCase)

	s.app = fiber.New(fiber.Config{
		StructValidator: &structValidator{validate: validator.New()},
	})

	s.app.Post("/api/users/register", handler.Register)
	s.app.Post("/api/users/login", handler.Login)
}

// TestRegister_Success tests successful registration
func (s *UserHandlerTestSuite) TestRegister_Success() {
	reqBody := dto.RequestRegister{
		Username: "newuser",
		Password: "password123",
	}
	body, _ := json.Marshal(reqBody)

	expectedOutput := &domain.UserOutput{
		ID:        1,
		Username:  reqBody.Username,
		CreatedAt: time.Now(),
	}

	s.mockUseCase.On("Register", mock.Anything, &domain.RegisterInput{
		Username: reqBody.Username,
		Password: reqBody.Password,
	}).Return(expectedOutput, nil)

	req := httptest.NewRequest(http.MethodPost, "/api/users/register", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.app.Test(req)
	s.NoError(err)
	s.Equal(fiber.StatusOK, resp.StatusCode)

	var apiResp dto.ApiResponse
	respBody, _ := io.ReadAll(resp.Body)
	s.NoError(json.Unmarshal(respBody, &apiResp))
	s.True(apiResp.Success)
	s.Equal("user registration was successful", apiResp.Message)
	s.mockUseCase.AssertExpectations(s.T())
}

// TestRegister_ValidationError tests failure when request payload is invalid
func (s *UserHandlerTestSuite) TestRegister_ValidationError() {
	// Password is under the minimum length (min=8)
	reqBody := dto.RequestRegister{
		Username: "ab",
		Password: "123",
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/api/users/register", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.app.Test(req)
	s.NoError(err)
	s.Equal(fiber.StatusBadRequest, resp.StatusCode)

	var apiResp dto.ApiResponse
	respBody, _ := io.ReadAll(resp.Body)
	s.NoError(json.Unmarshal(respBody, &apiResp))
	s.False(apiResp.Success)
	s.Equal("failed to validate request body", apiResp.Message)
	s.mockUseCase.AssertNotCalled(s.T(), "Register")
}

// TestRegister_AlreadyExists tests conflict when user is already registered
func (s *UserHandlerTestSuite) TestRegister_AlreadyExists() {
	reqBody := dto.RequestRegister{
		Username: "existinguser",
		Password: "password123",
	}
	body, _ := json.Marshal(reqBody)

	s.mockUseCase.On("Register", mock.Anything, &domain.RegisterInput{
		Username: reqBody.Username,
		Password: reqBody.Password,
	}).Return(nil, domain.ErrUserAlreadyExists)

	req := httptest.NewRequest(http.MethodPost, "/api/users/register", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.app.Test(req)
	s.NoError(err)
	s.Equal(fiber.StatusConflict, resp.StatusCode)

	var apiResp dto.ApiResponse
	respBody, _ := io.ReadAll(resp.Body)
	s.NoError(json.Unmarshal(respBody, &apiResp))
	s.False(apiResp.Success)
	s.Equal(domain.ErrUserAlreadyExists.Error(), apiResp.Message)
	s.mockUseCase.AssertExpectations(s.T())
}

// TestLogin_Success tests successful login
func (s *UserHandlerTestSuite) TestLogin_Success() {
	reqBody := dto.RequestLogin{
		Username: "validuser",
		Password: "password123",
	}
	body, _ := json.Marshal(reqBody)

	expectedOutput := &domain.UserOutput{
		ID:        10,
		Username:  reqBody.Username,
		CreatedAt: time.Now(),
	}

	s.mockUseCase.On("Login", mock.Anything, &domain.LoginInput{
		Username: reqBody.Username,
		Password: reqBody.Password,
	}).Return(expectedOutput, nil)

	req := httptest.NewRequest(http.MethodPost, "/api/users/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.app.Test(req)
	s.NoError(err)
	s.Equal(fiber.StatusOK, resp.StatusCode)

	var apiResp dto.ApiResponse
	respBody, _ := io.ReadAll(resp.Body)
	s.NoError(json.Unmarshal(respBody, &apiResp))
	s.True(apiResp.Success)
	s.Equal("user login was successful", apiResp.Message)
	s.mockUseCase.AssertExpectations(s.T())
}

// TestLogin_InvalidCredentials tests login with invalid username/password
func (s *UserHandlerTestSuite) TestLogin_InvalidCredentials() {
	reqBody := dto.RequestLogin{
		Username: "validuser",
		Password: "wrongpassword",
	}
	body, _ := json.Marshal(reqBody)

	s.mockUseCase.On("Login", mock.Anything, &domain.LoginInput{
		Username: reqBody.Username,
		Password: reqBody.Password,
	}).Return(nil, domain.ErrInvalidUser)

	req := httptest.NewRequest(http.MethodPost, "/api/users/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.app.Test(req)
	s.NoError(err)
	s.Equal(fiber.StatusBadRequest, resp.StatusCode)

	var apiResp dto.ApiResponse
	respBody, _ := io.ReadAll(resp.Body)
	s.NoError(json.Unmarshal(respBody, &apiResp))
	s.False(apiResp.Success)
	s.Equal(domain.ErrInvalidUser.Error(), apiResp.Message)
	s.mockUseCase.AssertExpectations(s.T())
}

// TestLogin_InternalError tests unexpected usecase errors
func (s *UserHandlerTestSuite) TestLogin_InternalError() {
	reqBody := dto.RequestLogin{
		Username: "validuser",
		Password: "password123",
	}
	body, _ := json.Marshal(reqBody)

	s.mockUseCase.On("Login", mock.Anything, &domain.LoginInput{
		Username: reqBody.Username,
		Password: reqBody.Password,
	}).Return(nil, errors.New("database connection down"))

	req := httptest.NewRequest(http.MethodPost, "/api/users/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.app.Test(req)
	s.NoError(err)
	s.Equal(fiber.StatusInternalServerError, resp.StatusCode)

	var apiResp dto.ApiResponse
	respBody, _ := io.ReadAll(resp.Body)
	s.NoError(json.Unmarshal(respBody, &apiResp))
	s.False(apiResp.Success)
	s.Equal("an unexpected error occurred", apiResp.Message)
	s.mockUseCase.AssertExpectations(s.T())
}

func TestUserHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(UserHandlerTestSuite))
}
