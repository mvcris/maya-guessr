package entities

import (
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)


type User struct {
	ID string `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	Name string `json:"name" gorm:"not null"`
	Email string `json:"email" gorm:"not null;unique"`
	Username string `json:"username" gorm:"not null;unique"`
	Password string `json:"-" gorm:"not null"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime;type:timestamptz"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime;type:timestamptz"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index;type:timestamptz"`
	Maps []*Map `json:"maps" gorm:"foreignKey:OwnerId"`
}

func (User) TableName() string {
	return "users"
}

func NewUser(name, email, username, password string) *User {
	return &User{
		Name: name,
		Email: email,
		Password: password,
		Username: username,
	}
}

func RestoreUser(id, name, email, username, password string) *User {
	return &User{
		ID: id,
		Name: name,
		Email: email,
		Password: password,
		Username: username,
	}
}

func (u *User) EncryptPassword() error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}

func (u *User) ComparePassword(password string) error {
	 return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
}