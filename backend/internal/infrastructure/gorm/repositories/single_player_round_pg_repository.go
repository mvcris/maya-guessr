package repositories

import (
	"context"

	"github.com/mvcris/maya-guessr/backend/internal/core/entities"
	"github.com/mvcris/maya-guessr/backend/internal/core/repositories"
	localgorm "github.com/mvcris/maya-guessr/backend/internal/infrastructure/gorm"
	"gorm.io/gorm"
)

type SinglePlayerRoundPgRepository struct {
	db *gorm.DB
}

func NewSinglePlayerRoundPgRepository(db *gorm.DB) repositories.SinglePlayerRoundRepository {
	return &SinglePlayerRoundPgRepository{db: db}
}

func (r *SinglePlayerRoundPgRepository) getDB(ctx context.Context) *gorm.DB {
	if tx, ok := localgorm.ExtractTx(ctx); ok {
		return tx.WithContext(ctx)
	}
	return r.db.WithContext(ctx)
}

func (r *SinglePlayerRoundPgRepository) Update(ctx context.Context, round *entities.SinglePlayerRound) error {
	return r.getDB(ctx).Save(round).Error
}