package repositories

import (
	"context"
	"errors"

	"github.com/mvcris/maya-guessr/backend/internal/core/entities"
	"github.com/mvcris/maya-guessr/backend/internal/core/repositories"
	localgorm "github.com/mvcris/maya-guessr/backend/internal/infrastructure/gorm"
	"gorm.io/gorm"
)

type RefreshTokenPgRepository struct {
	db *gorm.DB
}

func NewRefreshTokenPgRepository(db *gorm.DB) repositories.RefreshTokenRepository {
	return &RefreshTokenPgRepository{db: db}
}

func (r *RefreshTokenPgRepository) getDB(ctx context.Context) *gorm.DB {
	if tx, ok := localgorm.ExtractTx(ctx); ok {
		return tx.WithContext(ctx)
	}
	return r.db.WithContext(ctx)
}

func (r *RefreshTokenPgRepository) Create(ctx context.Context, refreshToken *entities.RefreshToken) error {
	return r.getDB(ctx).Create(refreshToken).Error
}

func (r *RefreshTokenPgRepository) FindById(ctx context.Context, id string) (*entities.RefreshToken, error) {
	var refreshToken entities.RefreshToken
	if err := r.getDB(ctx).Where("id = ?", id).First(&refreshToken).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &refreshToken, nil
}

func (r *RefreshTokenPgRepository) Update(ctx context.Context, refreshToken *entities.RefreshToken) error {
	return r.getDB(ctx).Save(refreshToken).Error
}
