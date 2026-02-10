package entities

import (
	"fmt"
	"time"

	coreerrors "github.com/mvcris/maya-guessr/backend/internal/core/errors"
	"gorm.io/gorm"
)

type Location struct {
	ID string `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	PanoId string `json:"pano_id" gorm:"not null;uniqueIndex:idx_location_pano_map"`
	MapId  string `json:"map_id" gorm:"not null;type:uuid;foreignKey:id;uniqueIndex:idx_location_pano_map"`
	Map *Map `json:"map" gorm:"foreignKey:MapId"`
	Latitude float64 `json:"latitude" gorm:"not null"`
	Longitude float64 `json:"longitude" gorm:"not null"`
	Heading float64 `json:"heading" gorm:"not null"`
	Pitch float64 `json:"pitch" gorm:"not null"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime;type:timestamptz"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime;type:timestamptz"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index;type:timestamptz"`
}

func (Location) TableName() string {
	return "locations"
}

func NewLocation(panoId, mapId string, latitude, longitude, heading, pitch float64) *Location {
	return &Location{
		PanoId: panoId,
		MapId: mapId,
		Latitude: latitude,
		Longitude: longitude,
		Heading: heading,
		Pitch: pitch,
	}	
}

func RestoreLocation(id, panoId, mapId string, latitude, longitude, heading, pitch float64) *Location {
	return &Location{
		ID: id,
		PanoId: panoId,
		MapId: mapId,
		Latitude: latitude,
		Longitude: longitude,
		Heading: heading,
		Pitch: pitch,
	}
}

func (l *Location) IsValidLatitude() bool {
	return l.Latitude >= -90 && l.Latitude <= 90
}

func (l *Location) IsValidLongitude() bool {
	return l.Longitude >= -180 && l.Longitude <= 180
}

func (l *Location) Validate() error {
	if !l.IsValidLatitude() {
		return coreerrors.BadRequest(fmt.Sprintf("invalid latitude: %f (must be between -90 and 90)", l.Latitude))
	}
	if !l.IsValidLongitude() {
		return coreerrors.BadRequest(fmt.Sprintf("invalid longitude: %f (must be between -180 and 180)", l.Longitude))
	}
	return nil
}