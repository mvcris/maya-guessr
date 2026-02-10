package entities

import (
	"time"

	"gorm.io/gorm"
)

type Map struct {
	ID string `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	Name string `json:"name" gorm:"not null;unique"`
	Description string `json:"description" gorm:"not null"`
	OwnerId string `json:"owner_id" gorm:"not null;type:uuid"`
	Owner *User `json:"owner" gorm:"foreignKey:OwnerId"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime;type:timestamptz"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime;type:timestamptz"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index;type:timestamptz"`
}

func (Map) TableName() string {
	return "maps"
}

func NewMap(name, description, ownerId string) *Map {
	return &Map{
		Name: name,
		Description: description,
		OwnerId: ownerId,
	}
}

func RestoreMap(id, name, description, ownerId string) *Map {
	return &Map{
		ID: id,
		Name: name,
		Description: description,
		OwnerId: ownerId,
	}
}