package repositories

import "github.com/mvcris/maya-guessr/backend/internal/core/entities"


type UserRepository interface {
	Create(user *entities.User) error
	FindByEmail(email string) (*entities.User, error)
	FindByUsername(username string) (*entities.User, error)
	FindById(id string) (*entities.User, error)
}