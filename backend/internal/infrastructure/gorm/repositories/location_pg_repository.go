package repositories

import (
	"context"

	"github.com/mvcris/maya-guessr/backend/internal/core/entities"
	"github.com/mvcris/maya-guessr/backend/internal/core/repositories"
	localgorm "github.com/mvcris/maya-guessr/backend/internal/infrastructure/gorm"
	"gorm.io/gorm"
)

type LocationPgRepository struct {
	db *gorm.DB
}

func NewLocationPgRepository(db *gorm.DB) repositories.LocationRepository {
	return &LocationPgRepository{db: db}
}

func (r *LocationPgRepository) getDB(ctx context.Context) *gorm.DB {
	if tx, ok := localgorm.ExtractTx(ctx); ok {
		return tx
	}
	return r.db
}

func (r *LocationPgRepository) Create(ctx context.Context, l *entities.Location) error {
	return r.getDB(ctx).Create(l).Error
}

func (r *LocationPgRepository) CountByMapId(ctx context.Context, mapId string) (int64, error) {
	var count int64
	if err := r.getDB(ctx).Model(&entities.Location{}).Where("map_id = ?", mapId).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}
