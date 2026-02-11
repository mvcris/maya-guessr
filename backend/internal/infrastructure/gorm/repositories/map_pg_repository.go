package repositories

import (
	"context"
	"errors"

	"github.com/mvcris/maya-guessr/backend/internal/core/entities"
	"github.com/mvcris/maya-guessr/backend/internal/core/repositories"
	localgorm "github.com/mvcris/maya-guessr/backend/internal/infrastructure/gorm"
	"gorm.io/gorm"
)

type MapPgRepository struct {
	db *gorm.DB
}

func NewMapPgRepository(db *gorm.DB) repositories.MapRepository {
	return &MapPgRepository{db: db}
}

func (r *MapPgRepository) getDB(ctx context.Context) *gorm.DB {
	if tx, ok := localgorm.ExtractTx(ctx); ok {
		return tx.WithContext(ctx)
	}
	return r.db.WithContext(ctx)
}

func (r *MapPgRepository) Create(ctx context.Context, m *entities.Map) error {
	return r.getDB(ctx).Create(m).Error
}

func (r *MapPgRepository) FindByName(ctx context.Context, name string) (*entities.Map, error) {
	var m entities.Map
	if err := r.getDB(ctx).Where("name = ?", name).First(&m).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &m, nil
}
