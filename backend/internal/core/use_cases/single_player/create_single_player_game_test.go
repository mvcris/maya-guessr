package singleplayer

import (
	"context"
	"errors"
	"testing"

	"github.com/mvcris/maya-guessr/backend/internal/core/entities"
	"github.com/mvcris/maya-guessr/backend/internal/core/repositories"
	repomocks "github.com/mvcris/maya-guessr/backend/internal/core/repositories/mocks"
	txmocks "github.com/mvcris/maya-guessr/backend/internal/core/transactions/mocks"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

var errMock = errors.New("mock error")

type CreateSinglePlayerGameSuite struct {
	suite.Suite
}

func TestCreateSinglePlayerGameSuite(t *testing.T) {
	suite.Run(t, new(CreateSinglePlayerGameSuite))
}

func passThroughTx(mockTx *txmocks.MockTransactionManager) {
	mockTx.EXPECT().
		RunInTransaction(mock.Anything, mock.Anything).
		RunAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
			return fn(ctx)
		})
}

func makeLocations(n int) []*entities.Location {
	locations := make([]*entities.Location, n)
	for i := 0; i < n; i++ {
		locations[i] = entities.RestoreLocation(
			"loc-uuid-"+string(rune('1'+i)),
			"pano-id",
			"map-uuid",
			20.0+float64(i),
			-100.0+float64(i),
			0, 0,
		)
	}
	return locations
}

func defaultInput() CreateSinglePlayerGameInput {
	return CreateSinglePlayerGameInput{
		UserId:               "user-uuid",
		MapId:                "map-uuid",
		Mode:                 entities.SinglePlayerGameModeMove,
		RoundSecondsDuration: 60,
	}
}

func (s *CreateSinglePlayerGameSuite) TestNewCreateSinglePlayerGameUseCase() {
	var gameRepo repositories.SinglePlayerGameRepository = repomocks.NewMockSinglePlayerGameRepository(s.T())
	var roundRepo repositories.SinglePlayerRoundRepository = repomocks.NewMockSinglePlayerRoundRepository(s.T())
	var locationRepo repositories.LocationRepository = repomocks.NewMockLocationRepository(s.T())
	mockTx := txmocks.NewMockTransactionManager(s.T())

	uc := NewCreateSinglePlayerGameUseCase(gameRepo, roundRepo, locationRepo, mockTx)
	s.NotNil(uc)
}

func (s *CreateSinglePlayerGameSuite) TestExecute_WhenSuccess_CreatesGameAndStartsFirstRound() {
	mockGameRepo := repomocks.NewMockSinglePlayerGameRepository(s.T())
	mockRoundRepo := repomocks.NewMockSinglePlayerRoundRepository(s.T())
	mockLocationRepo := repomocks.NewMockLocationRepository(s.T())
	mockTx := txmocks.NewMockTransactionManager(s.T())
	uc := NewCreateSinglePlayerGameUseCase(mockGameRepo, mockRoundRepo, mockLocationRepo, mockTx)

	input := defaultInput()
	locations := makeLocations(5)

	passThroughTx(mockTx)
	mockGameRepo.EXPECT().
		FindByUserIdAndStatuses(mock.Anything, input.UserId, []entities.SinglePlayerGameStatus{
			entities.SinglePlayerGameStatusInProgress,
			entities.SinglePlayerGameStatusPending,
		}).
		Return((*entities.SinglePlayerGame)(nil), nil)
	mockLocationRepo.EXPECT().
		FindRandomLocationByMapId(mock.Anything, input.MapId, 5).
		Return(locations, nil)
	mockGameRepo.EXPECT().
		Create(mock.Anything, mock.MatchedBy(func(g *entities.SinglePlayerGame) bool {
			return g.UserId == input.UserId &&
				g.MapId == input.MapId &&
				g.Mode == input.Mode &&
				g.RoundSecondsDuration == input.RoundSecondsDuration &&
				len(g.Rounds) == 5
		})).
		Return(nil)
	mockGameRepo.EXPECT().
		Update(mock.Anything, mock.MatchedBy(func(g *entities.SinglePlayerGame) bool {
			return g.Status == entities.SinglePlayerGameStatusInProgress &&
				g.CurrentRound == 1
		})).
		Return(nil)
	mockRoundRepo.EXPECT().
		Update(mock.Anything, mock.MatchedBy(func(r *entities.SinglePlayerRound) bool {
			return r.RoundStatus == entities.SinglePlayerRoundStatusInProgress &&
				r.RoundNumber == 1
		})).
		Return(nil)

	output, err := uc.Execute(input)

	s.Require().NoError(err)
	s.Equal(input.UserId, output.UserId)
	s.Equal(input.MapId, output.MapId)
	s.Equal(input.Mode, output.Mode)
	s.NotEmpty(output.ID)
	s.NotZero(output.CreatedAt)
}

func (s *CreateSinglePlayerGameSuite) TestExecute_WhenFindByUserIdAndStatusesFails_ReturnsError() {
	mockGameRepo := repomocks.NewMockSinglePlayerGameRepository(s.T())
	mockRoundRepo := repomocks.NewMockSinglePlayerRoundRepository(s.T())
	mockLocationRepo := repomocks.NewMockLocationRepository(s.T())
	mockTx := txmocks.NewMockTransactionManager(s.T())
	uc := NewCreateSinglePlayerGameUseCase(mockGameRepo, mockRoundRepo, mockLocationRepo, mockTx)

	input := defaultInput()

	passThroughTx(mockTx)
	mockGameRepo.EXPECT().
		FindByUserIdAndStatuses(mock.Anything, input.UserId, mock.Anything).
		Return((*entities.SinglePlayerGame)(nil), errMock)

	output, err := uc.Execute(input)

	s.Require().Error(err)
	s.Contains(err.Error(), "failed to find user in game")
	s.Equal(CreateSinglePlayerGameOutput{}, output)
}

func (s *CreateSinglePlayerGameSuite) TestExecute_WhenUserAlreadyInGame_ReturnsConflict() {
	mockGameRepo := repomocks.NewMockSinglePlayerGameRepository(s.T())
	mockRoundRepo := repomocks.NewMockSinglePlayerRoundRepository(s.T())
	mockLocationRepo := repomocks.NewMockLocationRepository(s.T())
	mockTx := txmocks.NewMockTransactionManager(s.T())
	uc := NewCreateSinglePlayerGameUseCase(mockGameRepo, mockRoundRepo, mockLocationRepo, mockTx)

	input := defaultInput()
	existingGame := entities.NewSinglePlayerGame(input.UserId, input.MapId, entities.SinglePlayerGameModeMove, 60)

	passThroughTx(mockTx)
	mockGameRepo.EXPECT().
		FindByUserIdAndStatuses(mock.Anything, input.UserId, mock.Anything).
		Return(existingGame, nil)

	output, err := uc.Execute(input)

	s.Require().Error(err)
	s.Contains(err.Error(), "user already in a game")
	s.Equal(CreateSinglePlayerGameOutput{}, output)
}

func (s *CreateSinglePlayerGameSuite) TestExecute_WhenFindRandomLocationsFails_ReturnsError() {
	mockGameRepo := repomocks.NewMockSinglePlayerGameRepository(s.T())
	mockRoundRepo := repomocks.NewMockSinglePlayerRoundRepository(s.T())
	mockLocationRepo := repomocks.NewMockLocationRepository(s.T())
	mockTx := txmocks.NewMockTransactionManager(s.T())
	uc := NewCreateSinglePlayerGameUseCase(mockGameRepo, mockRoundRepo, mockLocationRepo, mockTx)

	input := defaultInput()

	passThroughTx(mockTx)
	mockGameRepo.EXPECT().
		FindByUserIdAndStatuses(mock.Anything, input.UserId, mock.Anything).
		Return((*entities.SinglePlayerGame)(nil), nil)
	mockLocationRepo.EXPECT().
		FindRandomLocationByMapId(mock.Anything, input.MapId, 5).
		Return(nil, errMock)

	output, err := uc.Execute(input)

	s.Require().Error(err)
	s.Contains(err.Error(), "failed to find random locations")
	s.Equal(CreateSinglePlayerGameOutput{}, output)
}

func (s *CreateSinglePlayerGameSuite) TestExecute_WhenGameCreateFails_ReturnsError() {
	mockGameRepo := repomocks.NewMockSinglePlayerGameRepository(s.T())
	mockRoundRepo := repomocks.NewMockSinglePlayerRoundRepository(s.T())
	mockLocationRepo := repomocks.NewMockLocationRepository(s.T())
	mockTx := txmocks.NewMockTransactionManager(s.T())
	uc := NewCreateSinglePlayerGameUseCase(mockGameRepo, mockRoundRepo, mockLocationRepo, mockTx)

	input := defaultInput()
	locations := makeLocations(5)

	passThroughTx(mockTx)
	mockGameRepo.EXPECT().
		FindByUserIdAndStatuses(mock.Anything, input.UserId, mock.Anything).
		Return((*entities.SinglePlayerGame)(nil), nil)
	mockLocationRepo.EXPECT().
		FindRandomLocationByMapId(mock.Anything, input.MapId, 5).
		Return(locations, nil)
	mockGameRepo.EXPECT().
		Create(mock.Anything, mock.Anything).
		Return(errMock)

	output, err := uc.Execute(input)

	s.Require().Error(err)
	s.Contains(err.Error(), "failed to create game")
	s.Equal(CreateSinglePlayerGameOutput{}, output)
}

func (s *CreateSinglePlayerGameSuite) TestExecute_WhenGameUpdateFails_ReturnsError() {
	mockGameRepo := repomocks.NewMockSinglePlayerGameRepository(s.T())
	mockRoundRepo := repomocks.NewMockSinglePlayerRoundRepository(s.T())
	mockLocationRepo := repomocks.NewMockLocationRepository(s.T())
	mockTx := txmocks.NewMockTransactionManager(s.T())
	uc := NewCreateSinglePlayerGameUseCase(mockGameRepo, mockRoundRepo, mockLocationRepo, mockTx)

	input := defaultInput()
	locations := makeLocations(5)

	passThroughTx(mockTx)
	mockGameRepo.EXPECT().
		FindByUserIdAndStatuses(mock.Anything, input.UserId, mock.Anything).
		Return((*entities.SinglePlayerGame)(nil), nil)
	mockLocationRepo.EXPECT().
		FindRandomLocationByMapId(mock.Anything, input.MapId, 5).
		Return(locations, nil)
	mockGameRepo.EXPECT().
		Create(mock.Anything, mock.Anything).
		Return(nil)
	mockGameRepo.EXPECT().
		Update(mock.Anything, mock.Anything).
		Return(errMock)

	output, err := uc.Execute(input)

	s.Require().Error(err)
	s.Contains(err.Error(), "failed to update game")
	s.Equal(CreateSinglePlayerGameOutput{}, output)
}

func (s *CreateSinglePlayerGameSuite) TestExecute_WhenRoundUpdateFails_ReturnsError() {
	mockGameRepo := repomocks.NewMockSinglePlayerGameRepository(s.T())
	mockRoundRepo := repomocks.NewMockSinglePlayerRoundRepository(s.T())
	mockLocationRepo := repomocks.NewMockLocationRepository(s.T())
	mockTx := txmocks.NewMockTransactionManager(s.T())
	uc := NewCreateSinglePlayerGameUseCase(mockGameRepo, mockRoundRepo, mockLocationRepo, mockTx)

	input := defaultInput()
	locations := makeLocations(5)

	passThroughTx(mockTx)
	mockGameRepo.EXPECT().
		FindByUserIdAndStatuses(mock.Anything, input.UserId, mock.Anything).
		Return((*entities.SinglePlayerGame)(nil), nil)
	mockLocationRepo.EXPECT().
		FindRandomLocationByMapId(mock.Anything, input.MapId, 5).
		Return(locations, nil)
	mockGameRepo.EXPECT().
		Create(mock.Anything, mock.Anything).
		Return(nil)
	mockGameRepo.EXPECT().
		Update(mock.Anything, mock.Anything).
		Return(nil)
	mockRoundRepo.EXPECT().
		Update(mock.Anything, mock.Anything).
		Return(errMock)

	output, err := uc.Execute(input)

	s.Require().Error(err)
	s.Contains(err.Error(), "failed to update round")
	s.Equal(CreateSinglePlayerGameOutput{}, output)
}
