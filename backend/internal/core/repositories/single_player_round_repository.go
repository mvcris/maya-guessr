package repositories

import (
	"context"

	"github.com/mvcris/maya-guessr/backend/internal/core/entities"
)

type SinglePlayerRoundRepository interface {
	Update(ctx context.Context, round *entities.SinglePlayerRound) error
}