package repositories

import (
	"context"

	"github.com/mvcris/maya-guessr/backend/internal/core/entities"
)

type SinglePlayerRoundRepository interface {
	Update(ctx context.Context, round *entities.SinglePlayerRound) error
	FindByIdAndGameIdWithLock(ctx context.Context, id, gameId string) (*entities.SinglePlayerRound, error)
	FindByGameIdAndRoundNumberWithLock(ctx context.Context, gameId string, roundNumber int) (*entities.SinglePlayerRound, error)
}