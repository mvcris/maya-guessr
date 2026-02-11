package repositories

import (
	"context"

	"github.com/mvcris/maya-guessr/backend/internal/core/entities"
)

type UserRepository interface {
	Create(ctx context.Context, user *entities.User) error
	FindByEmail(ctx context.Context, email string) (*entities.User, error)
	FindByUsername(ctx context.Context, username string) (*entities.User, error)
	FindById(ctx context.Context, id string) (*entities.User, error)
}