package dtos

import (
	"time"

	"github.com/mvcris/maya-guessr/backend/internal/core/entities"
)

type CreateSinglePlayerGameRequest struct {
	MapId                 string `json:"map_id" binding:"required,uuid"`
	Mode                  string `json:"mode" binding:"required,oneof=move no_move nmpz"`
	RoundSecondsDuration  int    `json:"round_seconds_duration" binding:"required,min=10,max=300"`
}

type CreateSinglePlayerGameResponse struct {
	ID        string                      `json:"id"`
	UserId    string                      `json:"user_id"`
	MapId     string                      `json:"map_id"`
	Mode      entities.SinglePlayerGameMode `json:"mode"`
	CreatedAt time.Time                   `json:"created_at"`
}
