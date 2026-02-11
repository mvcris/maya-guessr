package gorm

import (
	"github.com/mvcris/maya-guessr/backend/internal/core/entities"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewConnection(dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		SkipDefaultTransaction: true,
	})
	if err != nil {
		return nil, err
	}
	err = db.AutoMigrate(&entities.User{}, &entities.RefreshToken{}, &entities.Map{}, &entities.Location{})
	if err != nil {
		return nil, err
	}
	return db, nil
}