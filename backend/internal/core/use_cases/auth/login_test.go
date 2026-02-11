package auth

import (
	"errors"
	"os"
	"testing"

	"github.com/mvcris/maya-guessr/backend/internal/core/entities"
	"github.com/mvcris/maya-guessr/backend/internal/core/repositories"
	repomocks "github.com/mvcris/maya-guessr/backend/internal/core/repositories/mocks"
	"github.com/mvcris/maya-guessr/backend/internal/core/services"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"golang.org/x/crypto/bcrypt"
)

type LoginSuite struct {
	suite.Suite
}

func TestLoginSuite(t *testing.T) {
	suite.Run(t, new(LoginSuite))
}

func (s *LoginSuite) TestExecute_WhenCredentialsValid_ReturnsTokens() {
	// JwtService requires JWT_SECRET_KEY
	prev := os.Getenv("JWT_SECRET_KEY")
	os.Setenv("JWT_SECRET_KEY", "test-secret-for-login-tests")
	defer func() { _ = os.Setenv("JWT_SECRET_KEY", prev) }()

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	s.Require().NoError(err)
	user := entities.RestoreUser("user-id", "John", "john@example.com", "johndoe", string(hashedPassword))

	mockUserRepo := repomocks.NewMockUserRepository(s.T())
	mockRefreshRepo := repomocks.NewMockRefreshTokenRepository(s.T())
	jwtService := services.NewJwtService()
	uc := NewLoginUseCase(mockUserRepo, mockRefreshRepo, jwtService)

	input := LoginInput{
		Email:    "john@example.com",
		Password: "password123",
	}

	mockUserRepo.EXPECT().
		FindByEmail(mock.Anything, input.Email).
		Return(user, nil)
	mockRefreshRepo.EXPECT().
		Create(mock.Anything, mock.MatchedBy(func(rt *entities.RefreshToken) bool {
			return rt.UserId == user.ID && !rt.ExpiresAt.IsZero()
		})).
		Return(nil)

	output, err := uc.Execute(input)

	s.Require().NoError(err)
	s.NotEmpty(output.AccessToken)
	s.NotEmpty(output.RefreshToken)
}

func (s *LoginSuite) TestExecute_WhenUserNotFound_ReturnsInvalidCredentials() {
	mockUserRepo := repomocks.NewMockUserRepository(s.T())
	mockRefreshRepo := repomocks.NewMockRefreshTokenRepository(s.T())
	uc := NewLoginUseCase(mockUserRepo, mockRefreshRepo, nil) // jwt not called

	input := LoginInput{
		Email:    "unknown@example.com",
		Password: "any",
	}

	mockUserRepo.EXPECT().
		FindByEmail(mock.Anything, input.Email).
		Return((*entities.User)(nil), nil)

	output, err := uc.Execute(input)

	s.Require().Error(err)
	s.Equal("invalid email or password", err.Error())
	s.Equal(LoginOutput{}, output)
}

func (s *LoginSuite) TestExecute_WhenPasswordInvalid_ReturnsInvalidCredentials() {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("correct"), bcrypt.DefaultCost)
	s.Require().NoError(err)
	user := entities.RestoreUser("user-id", "John", "john@example.com", "johndoe", string(hashedPassword))

	mockUserRepo := repomocks.NewMockUserRepository(s.T())
	mockRefreshRepo := repomocks.NewMockRefreshTokenRepository(s.T())
	uc := NewLoginUseCase(mockUserRepo, mockRefreshRepo, nil)

	input := LoginInput{
		Email:    "john@example.com",
		Password: "wrongpassword",
	}

	mockUserRepo.EXPECT().
		FindByEmail(mock.Anything, input.Email).
		Return(user, nil)

	output, err := uc.Execute(input)

	s.Require().Error(err)
	s.Equal("invalid email or password", err.Error())
	s.Equal(LoginOutput{}, output)
}

func (s *LoginSuite) TestExecute_WhenFindByEmailFails_ReturnsError() {
	mockUserRepo := repomocks.NewMockUserRepository(s.T())
	mockRefreshRepo := repomocks.NewMockRefreshTokenRepository(s.T())
	uc := NewLoginUseCase(mockUserRepo, mockRefreshRepo, nil)

	input := LoginInput{
		Email:    "john@example.com",
		Password: "password123",
	}
	findErr := errMock

	mockUserRepo.EXPECT().
		FindByEmail(mock.Anything, input.Email).
		Return((*entities.User)(nil), findErr)

	output, err := uc.Execute(input)

	s.Require().Error(err)
	s.ErrorIs(err, findErr)
	s.Equal(LoginOutput{}, output)
}

func (s *LoginSuite) TestExecute_WhenRefreshTokenCreateFails_ReturnsError() {
	prev := os.Getenv("JWT_SECRET_KEY")
	os.Setenv("JWT_SECRET_KEY", "test-secret")
	defer func() { _ = os.Setenv("JWT_SECRET_KEY", prev) }()

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	s.Require().NoError(err)
	user := entities.RestoreUser("user-id", "John", "john@example.com", "johndoe", string(hashedPassword))

	mockUserRepo := repomocks.NewMockUserRepository(s.T())
	mockRefreshRepo := repomocks.NewMockRefreshTokenRepository(s.T())
	jwtService := services.NewJwtService()
	uc := NewLoginUseCase(mockUserRepo, mockRefreshRepo, jwtService)

	input := LoginInput{
		Email:    "john@example.com",
		Password: "password123",
	}
	createErr := errMock

	mockUserRepo.EXPECT().
		FindByEmail(mock.Anything, input.Email).
		Return(user, nil)
	mockRefreshRepo.EXPECT().
		Create(mock.Anything, mock.MatchedBy(func(rt *entities.RefreshToken) bool {
			return rt.UserId == user.ID
		})).
		Return(createErr)

	output, err := uc.Execute(input)

	s.Require().Error(err)
	s.ErrorIs(err, createErr)
	s.Equal(LoginOutput{}, output)
}

func (s *LoginSuite) TestNewLoginUseCase() {
	var userRepo repositories.UserRepository = repomocks.NewMockUserRepository(s.T())
	var refreshRepo repositories.RefreshTokenRepository = repomocks.NewMockRefreshTokenRepository(s.T())
	uc := NewLoginUseCase(userRepo, refreshRepo, &services.JwtService{})
	s.NotNil(uc)
}

var errMock = errors.New("mock error")
