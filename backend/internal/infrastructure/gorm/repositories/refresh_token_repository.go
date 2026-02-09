package repositories

import (
	"errors"

	"github.com/mvcris/maya-guessr/backend/internal/core/entities"
	"github.com/mvcris/maya-guessr/backend/internal/core/repositories"
	"gorm.io/gorm"
)

type RefreshTokenPgRepository struct {
	db *gorm.DB
}

func NewRefreshTokenPgRepository(db *gorm.DB) repositories.RefreshTokenRepository {
	return &RefreshTokenPgRepository{db: db}
}

func (r *RefreshTokenPgRepository) Create(refreshToken *entities.RefreshToken) error {
	return r.db.Create(refreshToken).Error
}

func (r *RefreshTokenPgRepository) FindById(id string) (*entities.RefreshToken, error) {
	var refreshToken entities.RefreshToken
	if err := r.db.Where("id = ?", id).First(&refreshToken).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &refreshToken, nil
}

func (r *RefreshTokenPgRepository) Update(refreshToken *entities.RefreshToken) error {
	return r.db.Save(refreshToken).Error
}