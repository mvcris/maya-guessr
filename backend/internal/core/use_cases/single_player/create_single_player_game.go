package singleplayer

import (
	"context"
	"time"

	"github.com/mvcris/maya-guessr/backend/internal/core/entities"
	coreerrors "github.com/mvcris/maya-guessr/backend/internal/core/errors"
	"github.com/mvcris/maya-guessr/backend/internal/core/repositories"
	"github.com/mvcris/maya-guessr/backend/internal/core/transactions"
)

type CreateSinglePlayerGameInput struct {
	UserId string
	MapId string
	Mode entities.SinglePlayerGameMode
	RoundSecondsDuration int
}

type CreateSinglePlayerGameOutput struct {
	ID string
	UserId string
	MapId string
	Mode entities.SinglePlayerGameMode
	CreatedAt time.Time
}

type CreateSinglePlayerGameUseCase struct {
	singlePlayerGameRepository repositories.SinglePlayerGameRepository
	singlePlayerRoundRepository repositories.SinglePlayerRoundRepository
	locationRepository        repositories.LocationRepository
	txManager                 transactions.TransactionManager
}

func NewCreateSinglePlayerGameUseCase(
	singlePlayerGameRepository repositories.SinglePlayerGameRepository,
	singlePlayerRoundRepository repositories.SinglePlayerRoundRepository,
	locationRepository repositories.LocationRepository,
	txManager transactions.TransactionManager,
) *CreateSinglePlayerGameUseCase {
	return &CreateSinglePlayerGameUseCase{
		singlePlayerGameRepository: singlePlayerGameRepository,
		locationRepository:        locationRepository,
		singlePlayerRoundRepository: singlePlayerRoundRepository,
		txManager:                 txManager,
	}
}

func (uc *CreateSinglePlayerGameUseCase) Execute(input CreateSinglePlayerGameInput) (CreateSinglePlayerGameOutput, error) {
	ctx := context.Background()

	var output CreateSinglePlayerGameOutput
	err := uc.txManager.RunInTransaction(ctx, func(ctx context.Context) error {
		userAlreadyInGame, err := uc.singlePlayerGameRepository.FindByUserIdAndStatuses(ctx, input.UserId, []entities.SinglePlayerGameStatus{entities.SinglePlayerGameStatusInProgress, entities.SinglePlayerGameStatusPending})
		if err != nil {
			return coreerrors.InternalServerError("failed to find user in game")
		}

		if userAlreadyInGame != nil {
			return coreerrors.Conflict("user already in a game")
		}

		randomLocations, err := uc.locationRepository.FindRandomLocationByMapId(ctx, input.MapId, 5)
		if err != nil {
			return coreerrors.InternalServerError("failed to find random locations")
		}

		newGame := entities.NewSinglePlayerGame(input.UserId, input.MapId, input.Mode, input.RoundSecondsDuration)
		newGame.AddRoundsFromLocations(randomLocations)
		if err := uc.singlePlayerGameRepository.Create(ctx, newGame); err != nil {
			return coreerrors.InternalServerError("failed to create game")
		}

		if err := newGame.Start(); err != nil {
			return err
		}
		nextRound, err := newGame.StartNextRound()
		if err != nil {
			return err
		}


		if err := uc.singlePlayerGameRepository.Update(ctx, newGame); err != nil {
			return coreerrors.InternalServerError("failed to update game")
		}
		if err := uc.singlePlayerRoundRepository.Update(ctx, nextRound); err != nil {
			return coreerrors.InternalServerError("failed to update round")
		}

		output = CreateSinglePlayerGameOutput{
			ID:        newGame.ID,
			UserId:    newGame.UserId,
			MapId:     newGame.MapId,
			Mode:      newGame.Mode,
			CreatedAt: newGame.CreatedAt,
		}
		return nil
	})
	if err != nil {
		return CreateSinglePlayerGameOutput{}, err
	}

	return output, nil
}
