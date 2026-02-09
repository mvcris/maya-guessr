package coreerrors

import (
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/suite"
)

type ErrorsSuite struct {
	suite.Suite
}

func TestErrorsSuite(t *testing.T) {
	suite.Run(t, new(ErrorsSuite))
}

func (s *ErrorsSuite) TestConflict() {
	err := Conflict("resource already exists")
	s.Equal("resource already exists", err.Error())

	status, ok := Status(err)
	s.True(ok)
	s.Equal(http.StatusConflict, status)
}

func (s *ErrorsSuite) TestForbidden() {
	err := Forbidden("access denied")
	s.Equal("access denied", err.Error())

	status, ok := Status(err)
	s.True(ok)
	s.Equal(http.StatusForbidden, status)
}

func (s *ErrorsSuite) TestNotFound() {
	err := NotFound("user not found")
	s.Equal("user not found", err.Error())

	status, ok := Status(err)
	s.True(ok)
	s.Equal(http.StatusNotFound, status)
}

func (s *ErrorsSuite) TestUnauthorized() {
	err := Unauthorized("invalid email or password")
	s.Equal("invalid email or password", err.Error())

	status, ok := Status(err)
	s.True(ok)
	s.Equal(http.StatusUnauthorized, status)
}

func (s *ErrorsSuite) TestBadRequest() {
	err := BadRequest("invalid input")
	s.Equal("invalid input", err.Error())

	status, ok := Status(err)
	s.True(ok)
	s.Equal(http.StatusBadRequest, status)
}

func (s *ErrorsSuite) TestValidation_WithoutDetails() {
	err := Validation("invalid payload", nil)
	s.Equal("invalid payload", err.Error())

	status, ok := Status(err)
	s.True(ok)
	s.Equal(http.StatusBadRequest, status)
}

func (s *ErrorsSuite) TestStatus_WithNonDomainError() {
	err := errors.New("generic error")
	status, ok := Status(err)
	s.False(ok)
	s.Equal(0, status)
}

func (s *ErrorsSuite) TestStatus_WithWrappedDomainError() {
	err := Conflict("conflict")
	wrapped := errors.Join(err, errors.New("extra"))
	status, ok := Status(wrapped)
	s.True(ok)
	s.Equal(http.StatusConflict, status)
}
