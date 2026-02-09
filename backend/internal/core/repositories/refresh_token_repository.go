package repositories

import "github.com/mvcris/maya-guessr/backend/internal/core/entities"

type RefreshTokenRepository interface {
	Create(refreshToken *entities.RefreshToken) error
	FindById(id string) (*entities.RefreshToken, error)
	Update(refreshToken *entities.RefreshToken) error
}