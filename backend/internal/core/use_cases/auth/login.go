package auth

import (
	"errors"
	"time"

	"github.com/mvcris/maya-guessr/backend/internal/core/entities"
	"github.com/mvcris/maya-guessr/backend/internal/core/repositories"
	"github.com/mvcris/maya-guessr/backend/internal/core/services"
)

var ErrInvalidCredentials = errors.New("invalid email or password")

type LoginInput struct {
	Email string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type LoginOutput struct {
	AccessToken string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type LoginUseCase struct {
	userRepository repositories.UserRepository
	refreshTokenRepository repositories.RefreshTokenRepository
	jwtService *services.JwtService
}

func NewLoginUseCase(userRepository repositories.UserRepository, refreshTokenRepository repositories.RefreshTokenRepository, jwtService *services.JwtService) *LoginUseCase {
	return &LoginUseCase{userRepository: userRepository, refreshTokenRepository: refreshTokenRepository, jwtService: jwtService}
}

func (uc *LoginUseCase) Execute(input LoginInput) (LoginOutput, error) {
	user, err := uc.userRepository.FindByEmail(input.Email)
	if err != nil {
		return LoginOutput{}, err
	}
	if user == nil {
		return LoginOutput{}, ErrInvalidCredentials
	}
	if err := user.ComparePassword(input.Password); err != nil {
		return LoginOutput{}, ErrInvalidCredentials
	}
	refreshTokenEntity := entities.NewRefreshToken(user.ID, time.Now().Add(time.Hour * 24 * 7))
	if err := uc.refreshTokenRepository.Create(refreshTokenEntity); err != nil {
		return LoginOutput{}, err
	}
	accessToken, err := uc.jwtService.GenerateAccessToken(user.ID)
	if err != nil {
		return LoginOutput{}, err
	}
	refreshToken, err := uc.jwtService.GenerateRefreshToken(user.ID, refreshTokenEntity.ID)
	if err != nil {
		return LoginOutput{}, err
	}
	
	return LoginOutput{
		AccessToken: accessToken,
		RefreshToken: refreshToken,
	}, nil
}
