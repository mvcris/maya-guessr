package repositories

import (
	"errors"

	"github.com/mvcris/maya-guessr/backend/internal/core/entities"
	"github.com/mvcris/maya-guessr/backend/internal/core/repositories"
	"gorm.io/gorm"
)

type MapPgRepository struct {
	db *gorm.DB
}

func NewMapPgRepository(db *gorm.DB) repositories.MapRepository {
	return &MapPgRepository{db: db}
}

func (r *MapPgRepository) Create(m *entities.Map) error {
	return r.db.Create(m).Error
}

func (r *MapPgRepository) FindByName(name string) (*entities.Map, error) {
	var m entities.Map
	if err := r.db.Where("name = ?", name).First(&m).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &m, nil
}
