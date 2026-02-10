package entities

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type LocationSuite struct {
	suite.Suite
}

func TestLocationSuite(t *testing.T) {
	suite.Run(t, new(LocationSuite))
}

func (s *LocationSuite) TestTableName() {
	tableName := (Location{}).TableName()
	s.Equal("locations", tableName)
}

func (s *LocationSuite) TestNewLocation() {
	panoId := "pano-123"
	mapId := "map-uuid-456"
	latitude := 20.5
	longitude := -100.25
	heading := 45.0
	pitch := 10.0

	loc := NewLocation(panoId, mapId, latitude, longitude, heading, pitch)

	s.NotNil(loc)
	s.Empty(loc.ID)
	s.Equal(panoId, loc.PanoId)
	s.Equal(mapId, loc.MapId)
	s.Equal(latitude, loc.Latitude)
	s.Equal(longitude, loc.Longitude)
	s.Equal(heading, loc.Heading)
	s.Equal(pitch, loc.Pitch)
}

func (s *LocationSuite) TestRestoreLocation() {
	id := "loc-uuid-123"
	panoId := "pano-456"
	mapId := "map-uuid-789"
	latitude := 0.0
	longitude := 180.0
	heading := 90.0
	pitch := -5.0

	loc := RestoreLocation(id, panoId, mapId, latitude, longitude, heading, pitch)

	s.NotNil(loc)
	s.Equal(id, loc.ID)
	s.Equal(panoId, loc.PanoId)
	s.Equal(mapId, loc.MapId)
	s.Equal(latitude, loc.Latitude)
	s.Equal(longitude, loc.Longitude)
	s.Equal(heading, loc.Heading)
	s.Equal(pitch, loc.Pitch)
}

func (s *LocationSuite) TestIsValidLatitude_WhenValid() {
	loc := RestoreLocation("id", "pano", "map", 0, 0, 0, 0)
	s.True(loc.IsValidLatitude())

	loc.Latitude = 90
	s.True(loc.IsValidLatitude())

	loc.Latitude = -90
	s.True(loc.IsValidLatitude())

	loc.Latitude = 45.5
	s.True(loc.IsValidLatitude())
}

func (s *LocationSuite) TestIsValidLatitude_WhenInvalid() {
	loc := RestoreLocation("id", "pano", "map", 100, 0, 0, 0)
	s.False(loc.IsValidLatitude())

	loc.Latitude = -91
	s.False(loc.IsValidLatitude())

	loc.Latitude = 200
	s.False(loc.IsValidLatitude())
}

func (s *LocationSuite) TestIsValidLongitude_WhenValid() {
	loc := RestoreLocation("id", "pano", "map", 0, 0, 0, 0)
	s.True(loc.IsValidLongitude())

	loc.Longitude = 180
	s.True(loc.IsValidLongitude())

	loc.Longitude = -180
	s.True(loc.IsValidLongitude())

	loc.Longitude = -100.5
	s.True(loc.IsValidLongitude())
}

func (s *LocationSuite) TestIsValidLongitude_WhenInvalid() {
	loc := RestoreLocation("id", "pano", "map", 0, 200, 0, 0)
	s.False(loc.IsValidLongitude())

	loc.Longitude = -181
	s.False(loc.IsValidLongitude())

	loc.Longitude = 250
	s.False(loc.IsValidLongitude())
}

func (s *LocationSuite) TestValidate_WhenValid() {
	loc := NewLocation("pano", "map", 20.0, -100.0, 0, 0)

	err := loc.Validate()

	s.NoError(err)
}

func (s *LocationSuite) TestValidate_WhenInvalidLatitude() {
	loc := NewLocation("pano", "map", 95.0, 0, 0, 0)

	err := loc.Validate()

	s.Error(err)
	s.Contains(err.Error(), "invalid latitude")
	s.Contains(err.Error(), "95")
	s.Contains(err.Error(), "between -90 and 90")
}

func (s *LocationSuite) TestValidate_WhenInvalidLongitude() {
	loc := NewLocation("pano", "map", 0, 200.0, 0, 0)

	err := loc.Validate()

	s.Error(err)
	s.Contains(err.Error(), "invalid longitude")
	s.Contains(err.Error(), "200")
	s.Contains(err.Error(), "between -180 and 180")
}

func (s *LocationSuite) TestValidate_WhenBoundaryValues() {
	loc := NewLocation("pano", "map", 90, 180, 0, 0)
	err := loc.Validate()
	s.NoError(err)

	loc.Latitude = -90
	loc.Longitude = -180
	err = loc.Validate()
	s.NoError(err)
}
