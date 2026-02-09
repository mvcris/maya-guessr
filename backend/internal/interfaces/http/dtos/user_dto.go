package dtos

import "time"


type CreateUserRequest struct {
	Name string `json:"name" binding:"required"`
	Email string `json:"email" binding:"required,email"`
	Username string `json:"username" binding:"required,min=3,max=20"`
	Password string `json:"password" binding:"required,min=6,max=20"`
}

type CreateUserResponse struct {
	ID string `json:"id"`
	Name string `json:"name"`
	Email string `json:"email"`
	Username string `json:"username"`
	CreatedAt time.Time `json:"created_at"`
}