package services

import (
	"os"
	"testing"

	"github.com/stretchr/testify/suite"
)

const testJWTSecret = "test-secret-key-for-jwt-service-tests"

type JwtServiceSuite struct {
	suite.Suite
}

func TestJwtServiceSuite(t *testing.T) {
	suite.Run(t, new(JwtServiceSuite))
}

func (s *JwtServiceSuite) setupJWTEnv() func() {
	prev := os.Getenv("JWT_SECRET_KEY")
	os.Setenv("JWT_SECRET_KEY", testJWTSecret)
	return func() { _ = os.Setenv("JWT_SECRET_KEY", prev) }
}

func (s *JwtServiceSuite) TestNewJwtService_WhenSecretNotSet_Panics() {
	prev := os.Getenv("JWT_SECRET_KEY")
	os.Unsetenv("JWT_SECRET_KEY")
	defer func() { _ = os.Setenv("JWT_SECRET_KEY", prev) }()

	s.Require().Panics(func() { NewJwtService() })
}

func (s *JwtServiceSuite) TestNewJwtService_WhenSecretSet_ReturnsService() {
	defer s.setupJWTEnv()()
	svc := NewJwtService()
	s.Require().NotNil(svc)
}

func (s *JwtServiceSuite) TestGetSecretKey_ReturnsEnvValue() {
	defer s.setupJWTEnv()()
	svc := NewJwtService()
	key := svc.GetSecretKey()
	s.Equal([]byte(testJWTSecret), key)
}

func (s *JwtServiceSuite) TestGenerateAccessToken_ReturnsValidToken() {
	defer s.setupJWTEnv()()
	svc := NewJwtService()
	userID := "user-123"

	token, err := svc.GenerateAccessToken(userID)

	s.Require().NoError(err)
	s.NotEmpty(token)
}

func (s *JwtServiceSuite) TestGenerateAccessToken_ValidTokenCanBeValidated() {
	defer s.setupJWTEnv()()
	svc := NewJwtService()
	userID := "user-456"

	token, err := svc.GenerateAccessToken(userID)
	s.Require().NoError(err)

	claims, err := svc.ValidateAccessToken(token)
	s.Require().NoError(err)
	s.Equal(userID, claims.UserId)
	s.Equal(userID, claims.Subject)
	s.NotNil(claims.ExpiresAt)
}

func (s *JwtServiceSuite) TestGenerateRefreshToken_ReturnsValidToken() {
	defer s.setupJWTEnv()()
	svc := NewJwtService()
	userID := "user-789"
	refreshID := "refresh-token-id-123"

	token, err := svc.GenerateRefreshToken(userID, refreshID)

	s.Require().NoError(err)
	s.NotEmpty(token)
}

func (s *JwtServiceSuite) TestGenerateRefreshToken_ValidTokenCanBeValidated() {
	defer s.setupJWTEnv()()
	svc := NewJwtService()
	userID := "user-refresh"
	refreshID := "refresh-token-id-456"

	token, err := svc.GenerateRefreshToken(userID, refreshID)
	s.Require().NoError(err)

	claims, err := svc.ValidateRefreshToken(token)
	s.Require().NoError(err)
	s.Equal(userID, claims.UserId)
	s.Equal(userID, claims.Subject)
	s.Equal(refreshID, claims.ID)
	s.NotNil(claims.ExpiresAt)
}

func (s *JwtServiceSuite) TestValidateAccessToken_WhenTokenInvalid_ReturnsError() {
	defer s.setupJWTEnv()()
	svc := NewJwtService()

	_, err := svc.ValidateAccessToken("invalid-token")

	s.Require().Error(err)
}

func (s *JwtServiceSuite) TestValidateAccessToken_WhenTokenEmpty_ReturnsError() {
	defer s.setupJWTEnv()()
	svc := NewJwtService()

	_, err := svc.ValidateAccessToken("")

	s.Require().Error(err)
}

func (s *JwtServiceSuite) TestValidateRefreshToken_WhenTokenInvalid_ReturnsError() {
	defer s.setupJWTEnv()()
	svc := NewJwtService()

	_, err := svc.ValidateRefreshToken("invalid-refresh-token")

	s.Require().Error(err)
}

func (s *JwtServiceSuite) TestValidateRefreshToken_WhenTokenEmpty_ReturnsError() {
	defer s.setupJWTEnv()()
	svc := NewJwtService()

	_, err := svc.ValidateRefreshToken("")

	s.Require().Error(err)
}

func (s *JwtServiceSuite) TestValidateAccessToken_WhenTokenSignedWithDifferentSecret_ReturnsError() {
	defer s.setupJWTEnv()()
	svc := NewJwtService()
	token, err := svc.GenerateAccessToken("user-1")
	s.Require().NoError(err)

	prev := os.Getenv("JWT_SECRET_KEY")
	os.Setenv("JWT_SECRET_KEY", "different-secret")
	defer func() { _ = os.Setenv("JWT_SECRET_KEY", prev) }()
	svcOther := NewJwtService()

	_, err = svcOther.ValidateAccessToken(token)
	s.Require().Error(err)
}

func (s *JwtServiceSuite) TestValidateRefreshToken_WhenTokenSignedWithDifferentSecret_ReturnsError() {
	defer s.setupJWTEnv()()
	svc := NewJwtService()
	token, err := svc.GenerateRefreshToken("user-1", "refresh-id-1")
	s.Require().NoError(err)

	prev := os.Getenv("JWT_SECRET_KEY")
	os.Setenv("JWT_SECRET_KEY", "different-secret")
	defer func() { _ = os.Setenv("JWT_SECRET_KEY", prev) }()
	svcOther := NewJwtService()

	_, err = svcOther.ValidateRefreshToken(token)
	s.Require().Error(err)
}

