package entities

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type MapSuite struct {
	suite.Suite
}

func TestMapSuite(t *testing.T) {
	suite.Run(t, new(MapSuite))
}

func (s *MapSuite) TestTableName() {
	tableName := (Map{}).TableName()
	s.Equal("maps", tableName)
}

func (s *MapSuite) TestNewMap() {
	name := "My Cool Map"
	description := "A map with cool locations"
	ownerId := "owner-uuid-123"

	m := NewMap(name, description, ownerId)

	s.NotNil(m)
	s.Empty(m.ID)
	s.Equal(name, m.Name)
	s.Equal(description, m.Description)
	s.Equal(ownerId, m.OwnerId)
}

func (s *MapSuite) TestRestoreMap() {
	id := "map-uuid-456"
	name := "Restored Map"
	description := "Previously saved map"
	ownerId := "owner-uuid-789"

	m := RestoreMap(id, name, description, ownerId)

	s.NotNil(m)
	s.Equal(id, m.ID)
	s.Equal(name, m.Name)
	s.Equal(description, m.Description)
	s.Equal(ownerId, m.OwnerId)
}
