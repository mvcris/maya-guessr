package entities

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type UserSuite struct {
	suite.Suite
}

func TestUserSuite(t *testing.T) {
	suite.Run(t, new(UserSuite))
}

func (s *UserSuite) TestTableName() {
	tableName := (User{}).TableName()
	s.Equal("users", tableName)
}

func (s *UserSuite) TestNewUser() {
	name := "John Doe"
	email := "john@example.com"
	username := "johndoe"
	password := "secret123"

	u := NewUser(name, email, username, password)

	s.NotNil(u)
	s.Empty(u.ID)
	s.Equal(name, u.Name)
	s.Equal(email, u.Email)
	s.Equal(username, u.Username)
	s.Equal(password, u.Password)
}

func (s *UserSuite) TestRestoreUser() {
	id := "user-uuid-123"
	name := "John Doe"
	email := "john@example.com"
	username := "johndoe"
	password := "hashedpassword"

	u := RestoreUser(id, name, email, username, password)

	s.NotNil(u)
	s.Equal(id, u.ID)
	s.Equal(name, u.Name)
	s.Equal(email, u.Email)
	s.Equal(username, u.Username)
	s.Equal(password, u.Password)
}

func (s *UserSuite) TestEncryptPassword() {
	u := NewUser("John", "john@example.com", "johndoe", "plainpassword")

	err := u.EncryptPassword()

	s.NoError(err)
	s.NotEqual("plainpassword", u.Password)
	s.NotEmpty(u.Password)
}

func (s *UserSuite) TestComparePassword_WhenCorrect() {
	u := NewUser("John", "john@example.com", "johndoe", "correctpassword")
	s.Require().NoError(u.EncryptPassword())

	err := u.ComparePassword("correctpassword")

	s.NoError(err)
}

func (s *UserSuite) TestComparePassword_WhenWrong() {
	u := NewUser("John", "john@example.com", "johndoe", "correctpassword")
	s.Require().NoError(u.EncryptPassword())

	err := u.ComparePassword("wrongpassword")

	s.Error(err)
}
