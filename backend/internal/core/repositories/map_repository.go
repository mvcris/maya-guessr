package repositories

import "github.com/mvcris/maya-guessr/backend/internal/core/entities"

type MapRepository interface {
	Create(m *entities.Map) error
	FindByName(name string) (*entities.Map, error)
}
