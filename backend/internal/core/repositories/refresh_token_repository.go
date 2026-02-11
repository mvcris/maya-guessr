package repositories

import (
	"context"

	"github.com/mvcris/maya-guessr/backend/internal/core/entities"
)

type RefreshTokenRepository interface {
	Create(ctx context.Context, refreshToken *entities.RefreshToken) error
	FindById(ctx context.Context, id string) (*entities.RefreshToken, error)
	Update(ctx context.Context, refreshToken *entities.RefreshToken) error
}