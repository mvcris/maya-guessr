package repositories

import (
	"context"
	"errors"

	"github.com/mvcris/maya-guessr/backend/internal/core/entities"
	"github.com/mvcris/maya-guessr/backend/internal/core/repositories"
	localgorm "github.com/mvcris/maya-guessr/backend/internal/infrastructure/gorm"
	"gorm.io/gorm"
)

type UserPgRepository struct {
	db *gorm.DB
}

func NewUserPgRepository(db *gorm.DB) repositories.UserRepository {
	return &UserPgRepository{db: db}
}

func (r *UserPgRepository) getDB(ctx context.Context) *gorm.DB {
	if tx, ok := localgorm.ExtractTx(ctx); ok {
		return tx.WithContext(ctx)
	}
	return r.db.WithContext(ctx)
}

func (r *UserPgRepository) Create(ctx context.Context, user *entities.User) error {
	return r.getDB(ctx).Create(user).Error
}

func (r *UserPgRepository) FindByEmail(ctx context.Context, email string) (*entities.User, error) {
	var user entities.User
	if err := r.getDB(ctx).Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *UserPgRepository) FindByUsername(ctx context.Context, username string) (*entities.User, error) {
	var user entities.User
	if err := r.getDB(ctx).Where("username = ?", username).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *UserPgRepository) FindById(ctx context.Context, id string) (*entities.User, error) {
	var user entities.User
	if err := r.getDB(ctx).Where("id = ?", id).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}
