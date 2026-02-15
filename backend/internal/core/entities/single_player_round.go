package entities

import (
	"time"

	coreerrors "github.com/mvcris/maya-guessr/backend/internal/core/errors"
	"gorm.io/gorm"
)

type SinglePlayerRoundStatus string

const (
	SinglePlayerRoundStatusPending SinglePlayerRoundStatus = "pending"
	SinglePlayerRoundStatusInProgress SinglePlayerRoundStatus = "in_progress"
	SinglePlayerRoundStatusCompleted SinglePlayerRoundStatus = "completed"
)

type SinglePlayerRound struct {
	ID string `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	GameId string `json:"game_id" gorm:"not null;type:uuid"`
	Game *SinglePlayerGame `json:"game" gorm:"foreignKey:GameId"`
	LocationId string `json:"location_id" gorm:"not null;type:uuid"`
	Location *Location `json:"location" gorm:"foreignKey:LocationId"`
	Distance float64 `json:"distance"`
	Score int `json:"score"`
	GuessLatitude float64 `json:"guess_latitude"`
	GuessLongitude float64 `json:"guess_longitude"`
	StartedAt *time.Time `json:"started_at" gorm:"type:timestamptz;default:null"`
	EndedAt *time.Time `json:"ended_at" gorm:"type:timestamptz;default:null"`
	RoundNumber int `json:"round_number" gorm:"not null"`
	TotalRoundSecondsDuration int `json:"total_round_seconds_duration" gorm:"not null"`
	RoundStatus SinglePlayerRoundStatus `json:"round_status" gorm:"not null;default:pending"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime;type:timestamptz"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime;type:timestamptz"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index;type:timestamptz"`
}

func NewSinglePlayerRound(gameId, locationId string, roundNumber int, totalRoundSecondsDuration int) *SinglePlayerRound {
	return &SinglePlayerRound{
		GameId: gameId,
		LocationId: locationId,
		RoundNumber: roundNumber,
		TotalRoundSecondsDuration: totalRoundSecondsDuration,
		RoundStatus: SinglePlayerRoundStatusPending,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

func (r *SinglePlayerRound) Start() error {
	if r.RoundStatus != SinglePlayerRoundStatusPending {
		return coreerrors.BadRequest("round is not pending")
	}

	r.RoundStatus = SinglePlayerRoundStatusInProgress
	now := time.Now()
	r.StartedAt = &now
	return nil
}

func (r *SinglePlayerRound) Finish() error {
	if r.RoundStatus != SinglePlayerRoundStatusInProgress {
		return coreerrors.BadRequest("round is not in progress")
	}

	r.RoundStatus = SinglePlayerRoundStatusCompleted
	now := time.Now()
	r.EndedAt = &now
	return nil
}

func (r *SinglePlayerRound) IsInProgress() bool {
	return r.RoundStatus == SinglePlayerRoundStatusInProgress
}

// ApplyGuess records the guess coordinates and the precomputed distance (meters) and score.
// The caller (e.g., use case) must compute distance and score using its geo logic and pass them in.
func (r *SinglePlayerRound) ApplyGuess(guessLatitude, guessLongitude float64, distance float64, score int) {
	r.GuessLatitude = guessLatitude
	r.GuessLongitude = guessLongitude
	r.Distance = distance
	r.Score = score
}