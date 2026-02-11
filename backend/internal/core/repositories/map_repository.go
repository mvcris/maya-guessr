package repositories

import (
	"context"

	"github.com/mvcris/maya-guessr/backend/internal/core/entities"
)

type MapRepository interface {
	Create(ctx context.Context, m *entities.Map) error
	FindByName(ctx context.Context, name string) (*entities.Map, error)
}
