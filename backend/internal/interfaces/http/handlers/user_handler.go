package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mvcris/maya-guessr/backend/internal/core/use_cases/user"
	"github.com/mvcris/maya-guessr/backend/internal/infrastructure/gorm/repositories"
	httppkg "github.com/mvcris/maya-guessr/backend/internal/interfaces/http"
	"github.com/mvcris/maya-guessr/backend/internal/interfaces/http/dtos"
	"gorm.io/gorm"
)

type UserHandler struct {
	db *gorm.DB
	router *gin.Engine
	createUserUseCase *user.CreateUserUseCase
}

func NewUserHandler(db *gorm.DB, router *gin.Engine) *UserHandler {
	userPgRepository := repositories.NewUserPgRepository(db)
	return &UserHandler{ db: db, router: router, createUserUseCase: user.NewCreateUserUseCase(userPgRepository)}
}

func (h *UserHandler) CreateUser(c *gin.Context) {
	var input dtos.CreateUserRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}


	output, err := h.createUserUseCase.Execute(user.CreateUserInput{
		Name: input.Name,
		Email: input.Email,
		Username: input.Username,
		Password: input.Password,
	})
	if err != nil {
		httppkg.RespondError(c, err)
		return
	}
	c.JSON(http.StatusOK, output)
}

func (h *UserHandler) SetupRoutes() {
	h.router.POST("/users", h.CreateUser)
}