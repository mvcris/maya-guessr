package entities

import (
	"errors"
	"time"

	"gorm.io/gorm"
)

type RefreshToken struct {
	ID string `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	UserId string `json:"user_id" gorm:"not null,type:uuid;foreignKey:id"`
	ExpiresAt time.Time `json:"expires_at" gorm:"not null,type:timestamptz"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime;type:timestamptz"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime;type:timestamptz"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index;type:timestamptz"`
	User *User `json:"user" gorm:"foreignKey:UserId"`
}

func (RefreshToken) TableName() string {
	return "refresh_tokens"
}

func NewRefreshToken(userId string, expiresAt time.Time) *RefreshToken {
	return &RefreshToken{
		UserId: userId,
		ExpiresAt: expiresAt,
	}
}

func RestoreRefreshToken(id, userId string, expiresAt time.Time) *RefreshToken {
	return &RefreshToken{
		ID: id,
		UserId: userId,
		ExpiresAt: expiresAt,
	}
}

func (r *RefreshToken) IsExpired() bool {
	return time.Now().After(r.ExpiresAt)
}

func (r *RefreshToken) Expire() error {
	if r.IsExpired() {
		return errors.New("refresh token is expired")
	}
	r.ExpiresAt = time.Now()
	return nil
}