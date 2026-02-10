package repositories

import "github.com/mvcris/maya-guessr/backend/internal/core/entities"

type LocationRepository interface {
	Create(l *entities.Location) error
	CountByMapId(mapId string) (int64, error)
}
