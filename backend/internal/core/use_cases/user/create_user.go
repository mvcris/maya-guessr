package user

import (
	"context"
	"time"

	"github.com/mvcris/maya-guessr/backend/internal/core/entities"
	coreerrors "github.com/mvcris/maya-guessr/backend/internal/core/errors"
	"github.com/mvcris/maya-guessr/backend/internal/core/repositories"
)


type CreateUserInput struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type CreateUserOutput struct {
	ID string `json:"id"`
	Name string `json:"name"`
	Email string `json:"email"`
	Username string `json:"username"`
	CreatedAt time.Time `json:"created_at"`
}

type CreateUserUseCase struct {
	userRepository repositories.UserRepository
}

func NewCreateUserUseCase(userRepository repositories.UserRepository) *CreateUserUseCase {
	return &CreateUserUseCase{userRepository: userRepository}
}

func (uc *CreateUserUseCase) Execute(input CreateUserInput) (CreateUserOutput, error) {
	ctx := context.Background()

	existingUserByEmail, err := uc.userRepository.FindByEmail(ctx, input.Email)
	if err != nil {
		return CreateUserOutput{}, err
	}
	if existingUserByEmail != nil {
		return CreateUserOutput{}, coreerrors.Conflict("user with email already exists")
	}

	existingUserByUsername, err := uc.userRepository.FindByUsername(ctx, input.Username)
	if err != nil {
		return CreateUserOutput{}, err
	}
	if existingUserByUsername != nil {
		return CreateUserOutput{}, coreerrors.Conflict("user with username already exists")
	}

	newUser := entities.NewUser(input.Name, input.Email, input.Username, input.Password)
	if err := newUser.EncryptPassword(); err != nil {
		return CreateUserOutput{}, err
	}
	if err := uc.userRepository.Create(ctx, newUser); err != nil {
		return CreateUserOutput{}, err
	}
	return CreateUserOutput{
		ID:        newUser.ID,
		Name:      newUser.Name,
		Email:     newUser.Email,
		Username:  newUser.Username,
		CreatedAt: newUser.CreatedAt,
	}, nil
}