package mapuc

import (
	"context"
	"time"

	"github.com/mvcris/maya-guessr/backend/internal/core/entities"
	coreerrors "github.com/mvcris/maya-guessr/backend/internal/core/errors"
	"github.com/mvcris/maya-guessr/backend/internal/core/repositories"
	"github.com/mvcris/maya-guessr/backend/internal/core/transactions"
)

const maxLocationsPerMap = 50

type LocationInput struct {
	PanoId    string
	Latitude  float64
	Longitude float64
	Heading   float64
	Pitch     float64
}

type LocationOutput struct {
	ID        string  `json:"id"`
	PanoId    string  `json:"pano_id"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Heading   float64 `json:"heading"`
	Pitch     float64 `json:"pitch"`
}

type CreateMapInput struct {
	Name        string
	Description string
	OwnerId     string
	Locations   []LocationInput
}

type CreateMapOutput struct {
	ID          string           `json:"id"`
	Name        string           `json:"name"`
	Description string           `json:"description"`
	OwnerId     string           `json:"owner_id"`
	Locations   []LocationOutput `json:"locations"`
	CreatedAt   time.Time        `json:"created_at"`
}

type CreateMapUseCase struct {
	mapRepository      repositories.MapRepository
	locationRepository repositories.LocationRepository
	txManager          transactions.TransactionManager
}

func NewCreateMapUseCase(mapRepository repositories.MapRepository, locationRepository repositories.LocationRepository, txManager transactions.TransactionManager) *CreateMapUseCase {
	return &CreateMapUseCase{
		mapRepository:      mapRepository,
		locationRepository: locationRepository,
		txManager:          txManager,
	}
}

func (uc *CreateMapUseCase) Execute(input CreateMapInput) (CreateMapOutput, error) {
	if len(input.Locations) > maxLocationsPerMap {
		return CreateMapOutput{}, coreerrors.BadRequest("map cannot have more than 50 locations")
	}

	ctx := context.Background()

	existingMap, err := uc.mapRepository.FindByName(ctx, input.Name)
	if err != nil {
		return CreateMapOutput{}, err
	}
	if existingMap != nil {
		return CreateMapOutput{}, coreerrors.Conflict("map with this name already exists")
	}

	var output CreateMapOutput
	err = uc.txManager.RunInTransaction(ctx, func(ctx context.Context) error {
		newMap := entities.NewMap(input.Name, input.Description, input.OwnerId)
		if err := uc.mapRepository.Create(ctx, newMap); err != nil {
			return err
		}

		locationOutputs := make([]LocationOutput, 0, len(input.Locations))
		for _, loc := range input.Locations {
			location := entities.NewLocation(loc.PanoId, newMap.ID, loc.Latitude, loc.Longitude, loc.Heading, loc.Pitch)
			if err := location.Validate(); err != nil {
				return err
			}
			if err := uc.locationRepository.Create(ctx, location); err != nil {
				return err
			}
			locationOutputs = append(locationOutputs, LocationOutput{
				ID:        location.ID,
				PanoId:    location.PanoId,
				Latitude:  location.Latitude,
				Longitude: location.Longitude,
				Heading:   location.Heading,
				Pitch:     location.Pitch,
			})
		}

		output = CreateMapOutput{
			ID:          newMap.ID,
			Name:        newMap.Name,
			Description: newMap.Description,
			OwnerId:     newMap.OwnerId,
			Locations:   locationOutputs,
			CreatedAt:   newMap.CreatedAt,
		}
		return nil
	})

	return output, err
}
