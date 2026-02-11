package auth

import (
	"context"
	"time"

	"github.com/mvcris/maya-guessr/backend/internal/core/entities"
	coreerrors "github.com/mvcris/maya-guessr/backend/internal/core/errors"
	"github.com/mvcris/maya-guessr/backend/internal/core/repositories"
	"github.com/mvcris/maya-guessr/backend/internal/core/services"
)

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
	ctx := context.Background()

	user, err := uc.userRepository.FindByEmail(ctx, input.Email)
	if err != nil {
		return LoginOutput{}, err
	}
	if user == nil {
		return LoginOutput{}, coreerrors.Unauthorized("invalid email or password")
	}
	if err := user.ComparePassword(input.Password); err != nil {
		return LoginOutput{}, coreerrors.Unauthorized("invalid email or password")
	}
	refreshTokenEntity := entities.NewRefreshToken(user.ID, time.Now().Add(time.Hour*24*7))
	if err := uc.refreshTokenRepository.Create(ctx, refreshTokenEntity); err != nil {
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
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
