package entities

import (
	"time"

	"github.com/google/uuid"
	coreerrors "github.com/mvcris/maya-guessr/backend/internal/core/errors"
	"gorm.io/gorm"
)

type SinglePlayerGameStatus string

const (
	SinglePlayerGameStatusPending SinglePlayerGameStatus = "pending"
	SinglePlayerGameStatusInProgress SinglePlayerGameStatus = "in_progress"
	SinglePlayerGameStatusCompleted SinglePlayerGameStatus = "completed"
)

type SinglePlayerGameMode string

const (
	SinglePlayerGameModeMove SinglePlayerGameMode = "move"
	SinglePlayerGameModeNoMove SinglePlayerGameMode = "no_move"
	SinglePlayerGameModeNMPZ SinglePlayerGameMode = "nmpz"
)

type SinglePlayerGame struct {
	ID string `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	UserId string `json:"user_id" gorm:"not null;type:uuid"`
	User *User `json:"user" gorm:"foreignKey:UserId"`
	MapId string `json:"map_id" gorm:"not null;type:uuid"`
	Map *Map `json:"map" gorm:"foreignKey:MapId"`
	Score int `json:"score"`
	Status SinglePlayerGameStatus `json:"status" gorm:"not null;default:pending"`
	StartedAt *time.Time `json:"started_at" gorm:"type:timestamptz;default:null"`
	EndedAt *time.Time `json:"ended_at" gorm:"type:timestamptz;default:null"`
	Mode SinglePlayerGameMode `json:"mode" gorm:"not null"`
	TotalRounds int `json:"total_rounds" gorm:"not null"`
	CurrentRound int `json:"current_round" gorm:"not null"`
	RoundSecondsDuration int `json:"round_seconds_duration" gorm:"not null"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime;type:timestamptz"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime;type:timestamptz"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index;type:timestamptz"`
	Rounds []*SinglePlayerRound `json:"rounds" gorm:"foreignKey:GameId"`
}

func NewSinglePlayerGame(userId, mapId string, mode SinglePlayerGameMode, roundSecondsDuration int) *SinglePlayerGame {
	return &SinglePlayerGame{
		ID:     uuid.New().String(),
		UserId: userId,
		MapId:  mapId,
		Mode: mode,
		Status: SinglePlayerGameStatusPending,
		RoundSecondsDuration: roundSecondsDuration,
		TotalRounds: 5,
		CurrentRound: 0,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Rounds: make([]*SinglePlayerRound, 0),
	}
}

func (g *SinglePlayerGame) AddRoundsFromLocations(locations []*Location) {
	for i, location := range locations {
		g.Rounds = append(g.Rounds, NewSinglePlayerRound(g.ID, location.ID, i+1, g.RoundSecondsDuration))
	}
}

func (g *SinglePlayerGame) Start() error {
	if g.Status != SinglePlayerGameStatusPending {
		return coreerrors.BadRequest("game is not pending")
	}

	g.Status = SinglePlayerGameStatusInProgress
	now := time.Now()
	g.StartedAt = &now
	return nil
}

func (g *SinglePlayerGame) IsInProgress() bool {
	return g.Status == SinglePlayerGameStatusInProgress
}

func (g *SinglePlayerGame) HasNextRound() bool {
	return g.CurrentRound < g.TotalRounds
}

func (g *SinglePlayerGame) AdvanceRound() error {
	if g.Status != SinglePlayerGameStatusInProgress {
		return coreerrors.BadRequest("game is not in progress")
	}
	if g.CurrentRound >= g.TotalRounds {
		return coreerrors.BadRequest("game is already completed")
	}
	g.CurrentRound++
	return nil
}

func (g *SinglePlayerGame) AddScore(score int) {
	if score < 0 {
		return
	}
	g.Score += score
}

func (g *SinglePlayerGame) Complete() error {
	if g.Status != SinglePlayerGameStatusInProgress {
		return coreerrors.BadRequest("game is not in progress")
	}

	g.Status = SinglePlayerGameStatusCompleted
	now := time.Now()
	g.EndedAt = &now
	return nil
}

func (g *SinglePlayerGame) StartNextRound() (*SinglePlayerRound, error) {
	if !g.IsInProgress() {
		return nil, coreerrors.BadRequest("game is not in progress")
	}

	if g.CurrentRound >= g.TotalRounds {
		return nil, coreerrors.BadRequest("game is already completed")
	}

	if g.CurrentRound >= len(g.Rounds) {
		return nil, coreerrors.InternalServerError("round not found")
	}

	g.CurrentRound++
	round := g.Rounds[g.CurrentRound-1]
	if err := round.Start(); err != nil {
		return nil, err
	}
	return round, nil
}