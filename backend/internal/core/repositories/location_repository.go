package repositories

import (
	"context"

	"github.com/mvcris/maya-guessr/backend/internal/core/entities"
)

type LocationRepository interface {
	Create(ctx context.Context, l *entities.Location) error
	CountByMapId(ctx context.Context, mapId string) (int64, error)
}
