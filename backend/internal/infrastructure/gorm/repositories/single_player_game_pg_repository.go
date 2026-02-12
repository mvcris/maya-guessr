package repositories

import (
	"context"
	"errors"

	"github.com/mvcris/maya-guessr/backend/internal/core/entities"
	"github.com/mvcris/maya-guessr/backend/internal/core/repositories"
	localgorm "github.com/mvcris/maya-guessr/backend/internal/infrastructure/gorm"
	"gorm.io/gorm"
)

type SinglePlayerGamePgRepository struct {
	db *gorm.DB
}

func NewSinglePlayerGamePgRepository(db *gorm.DB) repositories.SinglePlayerGameRepository {
	return &SinglePlayerGamePgRepository{db: db}
}

func (r *SinglePlayerGamePgRepository) getDB(ctx context.Context) *gorm.DB {
	if tx, ok := localgorm.ExtractTx(ctx); ok {
		return tx.WithContext(ctx)
	}
	return r.db.WithContext(ctx)
}	

func (r *SinglePlayerGamePgRepository) Create(ctx context.Context, game *entities.SinglePlayerGame) error {
	return r.getDB(ctx).Create(game).Error
}

func (r *SinglePlayerGamePgRepository) FindByUserIdAndStatuses(ctx context.Context, userId string, statuses []entities.SinglePlayerGameStatus) (*entities.SinglePlayerGame, error) {
	var game entities.SinglePlayerGame
	if err := r.getDB(ctx).Where("user_id = ? AND status IN (?)", userId, statuses).First(&game).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &game, nil
}

func (r *SinglePlayerGamePgRepository) Update(ctx context.Context, game *entities.SinglePlayerGame) error {
	return r.getDB(ctx).Save(game).Error
}