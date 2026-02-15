package singleplayer

import (
	"context"
	"fmt"
	"strings"

	coreerrors "github.com/mvcris/maya-guessr/backend/internal/core/errors"
	"github.com/mvcris/maya-guessr/backend/internal/core/repositories"
	"github.com/mvcris/maya-guessr/backend/internal/core/services"
	"github.com/mvcris/maya-guessr/backend/internal/core/transactions"
)

type SinglePlayerGuessInput struct {
	GameId        string
	RoundId       string
	UserId        string
	GuessLatitude float64
	GuessLongitude float64
}

// Validate checks that all required fields are present and that coordinates are within valid ranges.
// Returns a BadRequest error if validation fails.
func (i SinglePlayerGuessInput) Validate() error {
	if strings.TrimSpace(i.GameId) == "" {
		return coreerrors.BadRequest("game id is required")
	}
	if strings.TrimSpace(i.RoundId) == "" {
		return coreerrors.BadRequest("round id is required")
	}
	if strings.TrimSpace(i.UserId) == "" {
		return coreerrors.BadRequest("user id is required")
	}
	if i.GuessLatitude < -90 || i.GuessLatitude > 90 {
		return coreerrors.BadRequest(fmt.Sprintf("guess latitude must be between -90 and 90, got %f", i.GuessLatitude))
	}
	if i.GuessLongitude < -180 || i.GuessLongitude > 180 {
		return coreerrors.BadRequest(fmt.Sprintf("guess longitude must be between -180 and 180, got %f", i.GuessLongitude))
	}
	return nil
}

type SinglePlayerGuessOutput struct {
	Score       int
	TotalScore  int
	GameEnded   bool
}

type SinglePlayerGuessUseCase struct {
	gameRepository  repositories.SinglePlayerGameRepository
	roundRepository repositories.SinglePlayerRoundRepository
	txManager       transactions.TransactionManager
	geoService      *services.GeoService
}

func NewSinglePlayerGuessUseCase(
	gameRepository repositories.SinglePlayerGameRepository,
	roundRepository repositories.SinglePlayerRoundRepository,
	txManager transactions.TransactionManager,
	geoService *services.GeoService,
) *SinglePlayerGuessUseCase {
	return &SinglePlayerGuessUseCase{
		gameRepository:  gameRepository,
		roundRepository: roundRepository,
		txManager:       txManager,
		geoService:      geoService,
	}
}

func (uc *SinglePlayerGuessUseCase) Execute(ctx context.Context, input SinglePlayerGuessInput) (SinglePlayerGuessOutput, error) {
	if err := input.Validate(); err != nil {
		return SinglePlayerGuessOutput{}, err
	}

	var output SinglePlayerGuessOutput
	err := uc.txManager.RunInTransaction(ctx, func(ctx context.Context) error {
		game, err := uc.gameRepository.FindByIdAndUserIdWithLock(ctx, input.GameId, input.UserId)
		if err != nil {
			return err
		}
		if game == nil {
			return coreerrors.NotFound("game not found")
		}
		if !game.IsInProgress() {
			return coreerrors.BadRequest("game is not in progress")
		}

		round, err := uc.roundRepository.FindByIdAndGameIdWithLock(ctx, input.RoundId, input.GameId)
		if err != nil {
			return err
		}
		if round == nil {
			return coreerrors.NotFound("round not found")
		}
		if !round.IsInProgress() {
			return coreerrors.BadRequest("round is not in progress")
		}

		var distance float64
		var score int
		if round.Location != nil {
			distance = uc.geoService.CalculateDistance(round.Location.Latitude, round.Location.Longitude, input.GuessLatitude, input.GuessLongitude)
			score = uc.geoService.CalculateScoreFromDistance(distance)
		}
		round.ApplyGuess(input.GuessLatitude, input.GuessLongitude, distance, score)
		if err := round.Finish(); err != nil {
			return err
		}
		if err := uc.roundRepository.Update(ctx, round); err != nil {
			return err
		}

		game.AddScore(round.Score)

		if game.HasNextRound() {
			nextRoundNumber := game.CurrentRound + 1
			nextRound, err := uc.roundRepository.FindByGameIdAndRoundNumberWithLock(ctx, game.ID, nextRoundNumber)
			if err != nil {
				return err
			}
			if nextRound == nil {
				return coreerrors.InternalServerError("next round not found")
			}
			if err := nextRound.Start(); err != nil {
				return err
			}
			if err := game.AdvanceRound(); err != nil {
				return err
			}
			if err := uc.roundRepository.Update(ctx, nextRound); err != nil {
				return err
			}
		} else {
			if err := game.Complete(); err != nil {
				return err
			}
			output.GameEnded = true
		}

		if err := uc.gameRepository.Update(ctx, game); err != nil {
			return err
		}

		output.Score = round.Score
		output.TotalScore = game.Score
		return nil
	})
	if err != nil {
		return SinglePlayerGuessOutput{}, err
	}

	return output, nil
}
