package user

import (
	"errors"
	"testing"
	"time"

	"github.com/mvcris/maya-guessr/backend/internal/core/entities"
	"github.com/mvcris/maya-guessr/backend/internal/core/repositories"
	repomocks "github.com/mvcris/maya-guessr/backend/internal/core/repositories/mocks"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type CreateUserSuite struct {
	suite.Suite
}

func TestCreateUserSuite(t *testing.T) {
	suite.Run(t, new(CreateUserSuite))
}

func (s *CreateUserSuite) TestExecute_WhenUserDoesNotExist_CreatesAndReturnsUser() {
	mockRepo := repomocks.NewMockUserRepository(s.T())
	uc := NewCreateUserUseCase(mockRepo)

	input := CreateUserInput{
		Name:     "John Doe",
		Email:    "john@example.com",
		Username: "johndoe",
		Password: "secret123",
	}

	mockRepo.EXPECT().
		FindByEmail(input.Email).
		Return((*entities.User)(nil), nil)
	mockRepo.EXPECT().
		FindByUsername(input.Username).
		Return((*entities.User)(nil), nil)
	mockRepo.EXPECT().
		Create(mock.MatchedBy(func(u *entities.User) bool {
			return u.Name == input.Name && u.Email == input.Email &&
				u.Username == input.Username && u.Password != input.Password
		})).
		RunAndReturn(func(u *entities.User) error {
			u.ID = "generated-uuid"
			u.CreatedAt = time.Date(2025, 2, 9, 12, 0, 0, 0, time.UTC)
			return nil
		})

	output, err := uc.Execute(input)

	s.Require().NoError(err)
	s.Equal("generated-uuid", output.ID)
	s.Equal(input.Name, output.Name)
	s.Equal(input.Email, output.Email)
	s.Equal(input.Username, output.Username)
	s.Equal(time.Date(2025, 2, 9, 12, 0, 0, 0, time.UTC), output.CreatedAt)
}

func (s *CreateUserSuite) TestExecute_WhenEmailAlreadyExists_ReturnsError() {
	mockRepo := repomocks.NewMockUserRepository(s.T())
	uc := NewCreateUserUseCase(mockRepo)

	input := CreateUserInput{
		Name:     "John Doe",
		Email:    "existing@example.com",
		Username: "johndoe",
		Password: "secret123",
	}
	existingUser := entities.RestoreUser("id-1", "Existing", "existing@example.com", "existing", "hash")

	mockRepo.EXPECT().
		FindByEmail(input.Email).
		Return(existingUser, nil)

	output, err := uc.Execute(input)

	s.Require().Error(err)
	s.Equal("user with email already exists", err.Error())
	s.Equal(CreateUserOutput{}, output)
}

func (s *CreateUserSuite) TestExecute_WhenUsernameAlreadyExists_ReturnsError() {
	mockRepo := repomocks.NewMockUserRepository(s.T())
	uc := NewCreateUserUseCase(mockRepo)

	input := CreateUserInput{
		Name:     "John Doe",
		Email:    "john@example.com",
		Username: "taken",
		Password: "secret123",
	}
	existingUser := entities.RestoreUser("id-1", "Other", "other@example.com", "taken", "hash")

	mockRepo.EXPECT().
		FindByEmail(input.Email).
		Return((*entities.User)(nil), nil)
	mockRepo.EXPECT().
		FindByUsername(input.Username).
		Return(existingUser, nil)

	output, err := uc.Execute(input)

	s.Require().Error(err)
	s.Equal("user with username already exists", err.Error())
	s.Equal(CreateUserOutput{}, output)
}

func (s *CreateUserSuite) TestExecute_WhenFindByEmailFails_ReturnsError() {
	mockRepo := repomocks.NewMockUserRepository(s.T())
	uc := NewCreateUserUseCase(mockRepo)

	input := CreateUserInput{
		Name:     "John Doe",
		Email:    "john@example.com",
		Username: "johndoe",
		Password: "secret123",
	}
	findErr := errMock

	mockRepo.EXPECT().
		FindByEmail(input.Email).
		Return((*entities.User)(nil), findErr)

	output, err := uc.Execute(input)

	s.Require().Error(err)
	s.ErrorIs(err, findErr)
	s.Equal(CreateUserOutput{}, output)
}

func (s *CreateUserSuite) TestExecute_WhenFindByUsernameFails_ReturnsError() {
	mockRepo := repomocks.NewMockUserRepository(s.T())
	uc := NewCreateUserUseCase(mockRepo)

	input := CreateUserInput{
		Name:     "John Doe",
		Email:    "john@example.com",
		Username: "johndoe",
		Password: "secret123",
	}
	findErr := errMock

	mockRepo.EXPECT().
		FindByEmail(input.Email).
		Return((*entities.User)(nil), nil)
	mockRepo.EXPECT().
		FindByUsername(input.Username).
		Return((*entities.User)(nil), findErr)

	output, err := uc.Execute(input)

	s.Require().Error(err)
	s.ErrorIs(err, findErr)
	s.Equal(CreateUserOutput{}, output)
}

func (s *CreateUserSuite) TestExecute_WhenCreateFails_ReturnsError() {
	mockRepo := repomocks.NewMockUserRepository(s.T())
	uc := NewCreateUserUseCase(mockRepo)

	input := CreateUserInput{
		Name:     "John Doe",
		Email:    "john@example.com",
		Username: "johndoe",
		Password: "secret123",
	}
	createErr := errMock

	mockRepo.EXPECT().
		FindByEmail(input.Email).
		Return((*entities.User)(nil), nil)
	mockRepo.EXPECT().
		FindByUsername(input.Username).
		Return((*entities.User)(nil), nil)
	mockRepo.EXPECT().
		Create(mock.MatchedBy(func(u *entities.User) bool {
			return u.Email == input.Email && u.Username == input.Username
		})).
		Return(createErr)

	output, err := uc.Execute(input)

	s.Require().Error(err)
	s.ErrorIs(err, createErr)
	s.Equal(CreateUserOutput{}, output)
}

func (s *CreateUserSuite) TestNewCreateUserUseCase() {
	var repo repositories.UserRepository = repomocks.NewMockUserRepository(s.T())
	uc := NewCreateUserUseCase(repo)
	s.NotNil(uc)
}

var errMock = errors.New("mock error")
