package repositories

import (
	"github.com/mvcris/maya-guessr/backend/internal/core/entities"
	"github.com/mvcris/maya-guessr/backend/internal/core/repositories"
	"gorm.io/gorm"
)

type LocationPgRepository struct {
	db *gorm.DB
}

func NewLocationPgRepository(db *gorm.DB) repositories.LocationRepository {
	return &LocationPgRepository{db: db}
}

func (r *LocationPgRepository) Create(l *entities.Location) error {
	return r.db.Create(l).Error
}

func (r *LocationPgRepository) CountByMapId(mapId string) (int64, error) {
	var count int64
	if err := r.db.Model(&entities.Location{}).Where("map_id = ?", mapId).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}
