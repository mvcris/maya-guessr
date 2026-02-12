package repositories

import (
	"context"

	"github.com/mvcris/maya-guessr/backend/internal/core/entities"
)

type SinglePlayerGameRepository interface {
	Create(ctx context.Context, game *entities.SinglePlayerGame) error
	FindByUserIdAndStatuses(ctx context.Context, userId string, statuses []entities.SinglePlayerGameStatus) (*entities.SinglePlayerGame, error)
	Update(ctx context.Context, game *entities.SinglePlayerGame) error
}