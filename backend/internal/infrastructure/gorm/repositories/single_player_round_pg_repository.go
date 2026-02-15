package repositories

import (
	"context"
	"errors"

	"github.com/mvcris/maya-guessr/backend/internal/core/entities"
	"github.com/mvcris/maya-guessr/backend/internal/core/repositories"
	localgorm "github.com/mvcris/maya-guessr/backend/internal/infrastructure/gorm"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
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

func (r *SinglePlayerRoundPgRepository) FindByIdAndGameIdWithLock(ctx context.Context, id, gameId string) (*entities.SinglePlayerRound, error) {
	var round entities.SinglePlayerRound
	if err := r.getDB(ctx).
		Clauses(clause.Locking{Strength: "UPDATE", Table: clause.Table{Name: clause.CurrentTable}}).
		Joins("Location").
		Where("single_player_rounds.id = ? AND single_player_rounds.game_id = ?", id, gameId).
		First(&round).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &round, nil
}

func (r *SinglePlayerRoundPgRepository) FindByGameIdAndRoundNumberWithLock(ctx context.Context, gameId string, roundNumber int) (*entities.SinglePlayerRound, error) {
	var round entities.SinglePlayerRound
	if err := r.getDB(ctx).
		Clauses(clause.Locking{Strength: "UPDATE", Table: clause.Table{Name: clause.CurrentTable}}).
		Joins("Location").
		Where("single_player_rounds.game_id = ? AND single_player_rounds.round_number = ?", gameId, roundNumber).
		First(&round).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &round, nil
}

func (r *SinglePlayerRoundPgRepository) Update(ctx context.Context, round *entities.SinglePlayerRound) error {
	return r.getDB(ctx).Save(round).Error
}