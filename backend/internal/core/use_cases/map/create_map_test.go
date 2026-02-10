package mapuc

import (
	"errors"
	"testing"
	"time"

	"github.com/mvcris/maya-guessr/backend/internal/core/entities"
	"github.com/mvcris/maya-guessr/backend/internal/core/repositories"
	repomocks "github.com/mvcris/maya-guessr/backend/internal/core/repositories/mocks"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type CreateMapSuite struct {
	suite.Suite
}

func TestCreateMapSuite(t *testing.T) {
	suite.Run(t, new(CreateMapSuite))
}

func (s *CreateMapSuite) TestExecute_WhenMapDoesNotExist_CreatesAndReturnsMap() {
	mockMapRepo := repomocks.NewMockMapRepository(s.T())
	mockLocationRepo := repomocks.NewMockLocationRepository(s.T())
	uc := NewCreateMapUseCase(mockMapRepo, mockLocationRepo)

	input := CreateMapInput{
		Name:        "My Map",
		Description: "A cool map",
		OwnerId:     "owner-123",
		Locations: []LocationInput{
			{PanoId: "pano-1", Latitude: 20.0, Longitude: -100.0, Heading: 0, Pitch: 0},
			{PanoId: "pano-2", Latitude: 21.0, Longitude: -101.0, Heading: 90, Pitch: 10},
		},
	}

	createdAt := time.Date(2025, 2, 9, 12, 0, 0, 0, time.UTC)

	mockMapRepo.EXPECT().
		FindByName(input.Name).
		Return((*entities.Map)(nil), nil)
	mockMapRepo.EXPECT().
		Create(mock.MatchedBy(func(m *entities.Map) bool {
			return m.Name == input.Name && m.Description == input.Description && m.OwnerId == input.OwnerId
		})).
		RunAndReturn(func(m *entities.Map) error {
			m.ID = "map-uuid-123"
			m.CreatedAt = createdAt
			return nil
		})
	mockLocationRepo.EXPECT().
		Create(mock.MatchedBy(func(l *entities.Location) bool {
			return l.PanoId == "pano-1" && l.MapId == "map-uuid-123"
		})).
		RunAndReturn(func(l *entities.Location) error {
			l.ID = "loc-uuid-1"
			return nil
		})
	mockLocationRepo.EXPECT().
		Create(mock.MatchedBy(func(l *entities.Location) bool {
			return l.PanoId == "pano-2" && l.MapId == "map-uuid-123"
		})).
		RunAndReturn(func(l *entities.Location) error {
			l.ID = "loc-uuid-2"
			return nil
		})

	output, err := uc.Execute(input)

	s.Require().NoError(err)
	s.Equal("map-uuid-123", output.ID)
	s.Equal(input.Name, output.Name)
	s.Equal(input.Description, output.Description)
	s.Equal(input.OwnerId, output.OwnerId)
	s.Equal(createdAt, output.CreatedAt)
	s.Len(output.Locations, 2)
	s.Equal("loc-uuid-1", output.Locations[0].ID)
	s.Equal("pano-1", output.Locations[0].PanoId)
	s.Equal(20.0, output.Locations[0].Latitude)
	s.Equal("loc-uuid-2", output.Locations[1].ID)
	s.Equal("pano-2", output.Locations[1].PanoId)
}

func (s *CreateMapSuite) TestExecute_WhenNoLocations_CreatesMapWithEmptyLocations() {
	mockMapRepo := repomocks.NewMockMapRepository(s.T())
	mockLocationRepo := repomocks.NewMockLocationRepository(s.T())
	uc := NewCreateMapUseCase(mockMapRepo, mockLocationRepo)

	input := CreateMapInput{
		Name:        "Empty Map",
		Description: "No locations",
		OwnerId:     "owner-123",
		Locations:   []LocationInput{},
	}

	createdAt := time.Date(2025, 2, 9, 12, 0, 0, 0, time.UTC)

	mockMapRepo.EXPECT().
		FindByName(input.Name).
		Return((*entities.Map)(nil), nil)
	mockMapRepo.EXPECT().
		Create(mock.MatchedBy(func(m *entities.Map) bool {
			return m.Name == input.Name
		})).
		RunAndReturn(func(m *entities.Map) error {
			m.ID = "map-uuid"
			m.CreatedAt = createdAt
			return nil
		})

	output, err := uc.Execute(input)

	s.Require().NoError(err)
	s.Equal("map-uuid", output.ID)
	s.Empty(output.Locations)
}

func (s *CreateMapSuite) TestExecute_WhenTooManyLocations_ReturnsBadRequest() {
	mockMapRepo := repomocks.NewMockMapRepository(s.T())
	mockLocationRepo := repomocks.NewMockLocationRepository(s.T())
	uc := NewCreateMapUseCase(mockMapRepo, mockLocationRepo)

	locations := make([]LocationInput, 51)
	for i := 0; i < 51; i++ {
		locations[i] = LocationInput{PanoId: "pano", Latitude: 0, Longitude: 0, Heading: 0, Pitch: 0}
	}
	input := CreateMapInput{
		Name:        "Map",
		Description: "Desc",
		OwnerId:     "owner",
		Locations:   locations,
	}

	output, err := uc.Execute(input)

	s.Require().Error(err)
	s.Equal("map cannot have more than 50 locations", err.Error())
	s.Equal(CreateMapOutput{}, output)
}

func (s *CreateMapSuite) TestExecute_WhenMapNameAlreadyExists_ReturnsConflict() {
	mockMapRepo := repomocks.NewMockMapRepository(s.T())
	mockLocationRepo := repomocks.NewMockLocationRepository(s.T())
	uc := NewCreateMapUseCase(mockMapRepo, mockLocationRepo)

	input := CreateMapInput{
		Name:        "Existing Map",
		Description: "Desc",
		OwnerId:     "owner",
		Locations:   []LocationInput{},
	}
	existingMap := entities.RestoreMap("id-1", "Existing Map", "desc", "other-owner")

	mockMapRepo.EXPECT().
		FindByName(input.Name).
		Return(existingMap, nil)

	output, err := uc.Execute(input)

	s.Require().Error(err)
	s.Equal("map with this name already exists", err.Error())
	s.Equal(CreateMapOutput{}, output)
}

func (s *CreateMapSuite) TestExecute_WhenFindByNameFails_ReturnsError() {
	mockMapRepo := repomocks.NewMockMapRepository(s.T())
	mockLocationRepo := repomocks.NewMockLocationRepository(s.T())
	uc := NewCreateMapUseCase(mockMapRepo, mockLocationRepo)

	input := CreateMapInput{
		Name:        "Map",
		Description: "Desc",
		OwnerId:     "owner",
		Locations:   []LocationInput{},
	}
	findErr := errMock

	mockMapRepo.EXPECT().
		FindByName(input.Name).
		Return((*entities.Map)(nil), findErr)

	output, err := uc.Execute(input)

	s.Require().Error(err)
	s.ErrorIs(err, findErr)
	s.Equal(CreateMapOutput{}, output)
}

func (s *CreateMapSuite) TestExecute_WhenMapCreateFails_ReturnsError() {
	mockMapRepo := repomocks.NewMockMapRepository(s.T())
	mockLocationRepo := repomocks.NewMockLocationRepository(s.T())
	uc := NewCreateMapUseCase(mockMapRepo, mockLocationRepo)

	input := CreateMapInput{
		Name:        "Map",
		Description: "Desc",
		OwnerId:     "owner",
		Locations:   []LocationInput{},
	}
	createErr := errMock

	mockMapRepo.EXPECT().
		FindByName(input.Name).
		Return((*entities.Map)(nil), nil)
	mockMapRepo.EXPECT().
		Create(mock.MatchedBy(func(m *entities.Map) bool {
			return m.Name == input.Name
		})).
		Return(createErr)

	output, err := uc.Execute(input)

	s.Require().Error(err)
	s.ErrorIs(err, createErr)
	s.Equal(CreateMapOutput{}, output)
}

func (s *CreateMapSuite) TestExecute_WhenLocationHasInvalidLatitude_ReturnsError() {
	mockMapRepo := repomocks.NewMockMapRepository(s.T())
	mockLocationRepo := repomocks.NewMockLocationRepository(s.T())
	uc := NewCreateMapUseCase(mockMapRepo, mockLocationRepo)

	input := CreateMapInput{
		Name:        "Map",
		Description: "Desc",
		OwnerId:     "owner",
		Locations: []LocationInput{
			{PanoId: "pano", Latitude: 200.0, Longitude: 0, Heading: 0, Pitch: 0},
		},
	}

	mockMapRepo.EXPECT().
		FindByName(input.Name).
		Return((*entities.Map)(nil), nil)
	mockMapRepo.EXPECT().
		Create(mock.MatchedBy(func(m *entities.Map) bool {
			return m.Name == input.Name
		})).
		RunAndReturn(func(m *entities.Map) error {
			m.ID = "map-uuid"
			m.CreatedAt = time.Now()
			return nil
		})

	output, err := uc.Execute(input)

	s.Require().Error(err)
	s.Contains(err.Error(), "invalid latitude")
	s.Equal(CreateMapOutput{}, output)
}

func (s *CreateMapSuite) TestExecute_WhenLocationHasInvalidLongitude_ReturnsError() {
	mockMapRepo := repomocks.NewMockMapRepository(s.T())
	mockLocationRepo := repomocks.NewMockLocationRepository(s.T())
	uc := NewCreateMapUseCase(mockMapRepo, mockLocationRepo)

	input := CreateMapInput{
		Name:        "Map",
		Description: "Desc",
		OwnerId:     "owner",
		Locations: []LocationInput{
			{PanoId: "pano", Latitude: 0, Longitude: 200.0, Heading: 0, Pitch: 0},
		},
	}

	mockMapRepo.EXPECT().
		FindByName(input.Name).
		Return((*entities.Map)(nil), nil)
	mockMapRepo.EXPECT().
		Create(mock.MatchedBy(func(m *entities.Map) bool {
			return m.Name == input.Name
		})).
		RunAndReturn(func(m *entities.Map) error {
			m.ID = "map-uuid"
			m.CreatedAt = time.Now()
			return nil
		})

	output, err := uc.Execute(input)

	s.Require().Error(err)
	s.Contains(err.Error(), "invalid longitude")
	s.Equal(CreateMapOutput{}, output)
}

func (s *CreateMapSuite) TestExecute_WhenLocationCreateFails_ReturnsError() {
	mockMapRepo := repomocks.NewMockMapRepository(s.T())
	mockLocationRepo := repomocks.NewMockLocationRepository(s.T())
	uc := NewCreateMapUseCase(mockMapRepo, mockLocationRepo)

	input := CreateMapInput{
		Name:        "Map",
		Description: "Desc",
		OwnerId:     "owner",
		Locations: []LocationInput{
			{PanoId: "pano", Latitude: 20.0, Longitude: -100.0, Heading: 0, Pitch: 0},
		},
	}
	createErr := errMock

	mockMapRepo.EXPECT().
		FindByName(input.Name).
		Return((*entities.Map)(nil), nil)
	mockMapRepo.EXPECT().
		Create(mock.MatchedBy(func(m *entities.Map) bool {
			return m.Name == input.Name
		})).
		RunAndReturn(func(m *entities.Map) error {
			m.ID = "map-uuid"
			m.CreatedAt = time.Now()
			return nil
		})
	mockLocationRepo.EXPECT().
		Create(mock.MatchedBy(func(l *entities.Location) bool {
			return l.PanoId == "pano"
		})).
		Return(createErr)

	output, err := uc.Execute(input)

	s.Require().Error(err)
	s.ErrorIs(err, createErr)
	s.Equal(CreateMapOutput{}, output)
}

func (s *CreateMapSuite) TestNewCreateMapUseCase() {
	var mapRepo repositories.MapRepository = repomocks.NewMockMapRepository(s.T())
	var locationRepo repositories.LocationRepository = repomocks.NewMockLocationRepository(s.T())
	uc := NewCreateMapUseCase(mapRepo, locationRepo)
	s.NotNil(uc)
}

var errMock = errors.New("mock error")
