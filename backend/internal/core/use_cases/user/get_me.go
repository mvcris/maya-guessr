package user

import (
	"github.com/mvcris/maya-guessr/backend/internal/core/entities"
	coreerrors "github.com/mvcris/maya-guessr/backend/internal/core/errors"
	"github.com/mvcris/maya-guessr/backend/internal/core/repositories"
)

type GetMeUseCase struct {
	userRepository repositories.UserRepository
}

func NewGetMeUseCase(userRepository repositories.UserRepository) *GetMeUseCase {
	return &GetMeUseCase{userRepository: userRepository}
}

func (uc *GetMeUseCase) Execute(userID string) (*entities.User, error) {
	user, err := uc.userRepository.FindById(userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, coreerrors.NotFound("user not found")
	}
	return user, nil
}
