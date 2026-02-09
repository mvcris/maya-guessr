package services

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)


type JwtService struct {}

type AccessTokenClaims struct {
	jwt.RegisteredClaims
	UserId string `json:"user_id"`
}

type RefreshTokenClaims struct {
	jwt.RegisteredClaims
	ID     string `json:"id"`
	UserId string `json:"user_id"`
}

func NewJwtService() *JwtService {
	secretKey := os.Getenv("JWT_SECRET_KEY")
	if secretKey == "" {
		panic("JWT_SECRET_KEY is not set")
	}
	return &JwtService{}
}

func (s *JwtService) GetSecretKey() []byte {
	return []byte(os.Getenv("JWT_SECRET_KEY"))
}

func (s *JwtService) GenerateAccessToken(userId string) (string, error) {
	claims := &AccessTokenClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userId,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
		},
		UserId: userId,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.GetSecretKey())
}

func (s *JwtService) GenerateRefreshToken(userId string, id string) (string, error) {
	claims := &RefreshTokenClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userId,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * 7)),
		},
		UserId: userId,
		ID: id,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.GetSecretKey())
}

func (s *JwtService) ValidateAccessToken(tokenString string) (*AccessTokenClaims, error) {
	tok, err := jwt.ParseWithClaims(tokenString, &AccessTokenClaims{}, func(t *jwt.Token) (interface{}, error) {
		return s.GetSecretKey(), nil
	})
	if err != nil {
		return nil, err
	}
	return tok.Claims.(*AccessTokenClaims), nil
}

func (s *JwtService) ValidateRefreshToken(tokenString string) (*RefreshTokenClaims, error) {
	tok, err := jwt.ParseWithClaims(tokenString, &RefreshTokenClaims{}, func(t *jwt.Token) (interface{}, error) {
		return s.GetSecretKey(), nil
	})
	if err != nil {
		return nil, err
	}
	return tok.Claims.(*RefreshTokenClaims), nil
}