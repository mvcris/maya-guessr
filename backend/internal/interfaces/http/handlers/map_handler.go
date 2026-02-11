package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	coreerrors "github.com/mvcris/maya-guessr/backend/internal/core/errors"
	"github.com/mvcris/maya-guessr/backend/internal/core/services"
	mapuc "github.com/mvcris/maya-guessr/backend/internal/core/use_cases/map"
	localgorm "github.com/mvcris/maya-guessr/backend/internal/infrastructure/gorm"
	"github.com/mvcris/maya-guessr/backend/internal/infrastructure/gorm/repositories"
	httppkg "github.com/mvcris/maya-guessr/backend/internal/interfaces/http"
	"github.com/mvcris/maya-guessr/backend/internal/interfaces/http/dtos"
	"github.com/mvcris/maya-guessr/backend/internal/interfaces/http/middleware"
	"gorm.io/gorm"
)

type MapHandler struct {
	createMapUseCase *mapuc.CreateMapUseCase
	jwtService       *services.JwtService
	router           *gin.Engine
}

func NewMapHandler(db *gorm.DB, router *gin.Engine) *MapHandler {
	mapRepository := repositories.NewMapPgRepository(db)
	locationRepository := repositories.NewLocationPgRepository(db)
	txManager := localgorm.NewGormTransactionManager(db)
	jwtService := services.NewJwtService()
	return &MapHandler{
		createMapUseCase: mapuc.NewCreateMapUseCase(mapRepository, locationRepository, txManager),
		jwtService:       jwtService,
		router:           router,
	}
}

func (h *MapHandler) CreateMap(c *gin.Context) {
	userID, ok := middleware.GetAuthenticatedUserID(c)
	if !ok {
		httppkg.RespondError(c, coreerrors.Unauthorized("invalid or missing token"))
		return
	}

	var input dtos.CreateMapRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	locations := make([]mapuc.LocationInput, len(input.Locations))
	for i, loc := range input.Locations {
		locations[i] = mapuc.LocationInput{
			PanoId:   loc.PanoId,
			Latitude: loc.Latitude,
			Longitude: loc.Longitude,
			Heading:  loc.Heading,
			Pitch:    loc.Pitch,
		}
	}

	output, err := h.createMapUseCase.Execute(mapuc.CreateMapInput{
		Name:        input.Name,
		Description: input.Description,
		OwnerId:     userID,
		Locations:   locations,
	})
	if err != nil {
		httppkg.RespondError(c, err)
		return
	}

	locationDTOs := make([]dtos.LocationOutputDTO, len(output.Locations))
	for i, loc := range output.Locations {
		locationDTOs[i] = dtos.LocationOutputDTO{
			ID:        loc.ID,
			PanoId:    loc.PanoId,
			Latitude:  loc.Latitude,
			Longitude: loc.Longitude,
			Heading:   loc.Heading,
			Pitch:     loc.Pitch,
		}
	}

	c.JSON(http.StatusOK, dtos.CreateMapResponse{
		ID:          output.ID,
		Name:        output.Name,
		Description: output.Description,
		OwnerId:     output.OwnerId,
		Locations:   locationDTOs,
		CreatedAt:   output.CreatedAt,
	})
}

func (h *MapHandler) SetupRoutes() {
	h.router.POST("/maps", middleware.AuthMiddleware(h.jwtService), h.CreateMap)
}
