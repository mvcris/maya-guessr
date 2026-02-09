package entities

import (
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
)

type RefreshTokenSuite struct {
	suite.Suite
}

func TestRefreshTokenSuite(t *testing.T) {
	suite.Run(t, new(RefreshTokenSuite))
}

func (s *RefreshTokenSuite) TestTableName() {
	tableName := (RefreshToken{}).TableName()
	s.Equal("refresh_tokens", tableName)
}

func (s *RefreshTokenSuite) TestNewRefreshToken() {
	userId := "user-uuid-123"
	expiresAt := time.Now().Add(24 * time.Hour)

	rt := NewRefreshToken(userId, expiresAt)

	s.NotNil(rt)
	s.Equal(userId, rt.UserId)
	s.Equal(expiresAt, rt.ExpiresAt)
	s.Empty(rt.ID)
}

func (s *RefreshTokenSuite) TestRestoreRefreshToken() {
	id := "token-uuid-456"
	userId := "user-uuid-123"
	expiresAt := time.Now().Add(24 * time.Hour)

	rt := RestoreRefreshToken(id, userId, expiresAt)

	s.NotNil(rt)
	s.Equal(id, rt.ID)
	s.Equal(userId, rt.UserId)
	s.Equal(expiresAt, rt.ExpiresAt)
}

func (s *RefreshTokenSuite) TestIsExpired_WhenNotExpired() {
	rt := NewRefreshToken("user-id", time.Now().Add(time.Hour))
	s.False(rt.IsExpired())
}

func (s *RefreshTokenSuite) TestIsExpired_WhenExpired() {
	rt := NewRefreshToken("user-id", time.Now().Add(-time.Hour))
	s.True(rt.IsExpired())
}

func (s *RefreshTokenSuite) TestExpire_WhenNotExpired_Succeeds() {
	rt := NewRefreshToken("user-id", time.Now().Add(time.Hour))
	beforeExpire := rt.ExpiresAt

	err := rt.Expire()

	s.NoError(err)
	s.True(rt.ExpiresAt.Before(beforeExpire) || rt.ExpiresAt.Equal(beforeExpire))
	s.True(rt.IsExpired())
}

func (s *RefreshTokenSuite) TestExpire_WhenAlreadyExpired_ReturnsError() {
	rt := NewRefreshToken("user-id", time.Now().Add(-time.Hour))

	err := rt.Expire()

	s.Error(err)
	s.Equal("refresh token is expired", err.Error())
}
