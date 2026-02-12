package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mvcris/maya-guessr/backend/internal/core/entities"
	coreerrors "github.com/mvcris/maya-guessr/backend/internal/core/errors"
	"github.com/mvcris/maya-guessr/backend/internal/core/services"
	singleplayer "github.com/mvcris/maya-guessr/backend/internal/core/use_cases/single_player"
	localgorm "github.com/mvcris/maya-guessr/backend/internal/infrastructure/gorm"
	"github.com/mvcris/maya-guessr/backend/internal/infrastructure/gorm/repositories"
	httppkg "github.com/mvcris/maya-guessr/backend/internal/interfaces/http"
	"github.com/mvcris/maya-guessr/backend/internal/interfaces/http/dtos"
	"github.com/mvcris/maya-guessr/backend/internal/interfaces/http/middleware"
	"gorm.io/gorm"
)

type SinglePlayerHandler struct {
	createSinglePlayerGameUseCase *singleplayer.CreateSinglePlayerGameUseCase
	jwtService                    *services.JwtService
	router                        *gin.Engine
}

func NewSinglePlayerHandler(db *gorm.DB, router *gin.Engine) *SinglePlayerHandler {
	singlePlayerGameRepository := repositories.NewSinglePlayerGamePgRepository(db)
	singlePlayerRoundRepository := repositories.NewSinglePlayerRoundPgRepository(db)
	locationRepository := repositories.NewLocationPgRepository(db)
	txManager := localgorm.NewGormTransactionManager(db)
	jwtService := services.NewJwtService()
	return &SinglePlayerHandler{
		createSinglePlayerGameUseCase: singleplayer.NewCreateSinglePlayerGameUseCase(singlePlayerGameRepository, singlePlayerRoundRepository, locationRepository, txManager),
		jwtService:                    jwtService,
		router:                        router,
	}
}

func (h *SinglePlayerHandler) CreateGame(c *gin.Context) {
	userID, ok := middleware.GetAuthenticatedUserID(c)
	if !ok {
		httppkg.RespondError(c, coreerrors.Unauthorized("invalid or missing token"))
		return
	}

	var input dtos.CreateSinglePlayerGameRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	output, err := h.createSinglePlayerGameUseCase.Execute(singleplayer.CreateSinglePlayerGameInput{
		UserId:                userID,
		MapId:                 input.MapId,
		Mode:                  entities.SinglePlayerGameMode(input.Mode),
		RoundSecondsDuration:  input.RoundSecondsDuration,
	})
	if err != nil {
		httppkg.RespondError(c, err)
		return
	}

	c.JSON(http.StatusCreated, dtos.CreateSinglePlayerGameResponse{
		ID:        output.ID,
		UserId:    output.UserId,
		MapId:     output.MapId,
		Mode:      output.Mode,
		CreatedAt: output.CreatedAt,
	})
}

func (h *SinglePlayerHandler) SetupRoutes() {
	h.router.POST("/single-player/games", middleware.AuthMiddleware(h.jwtService), h.CreateGame)
}
