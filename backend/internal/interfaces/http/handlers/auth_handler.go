package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mvcris/maya-guessr/backend/internal/core/services"
	"github.com/mvcris/maya-guessr/backend/internal/core/use_cases/auth"
	"github.com/mvcris/maya-guessr/backend/internal/infrastructure/gorm/repositories"
	httppkg "github.com/mvcris/maya-guessr/backend/internal/interfaces/http"
	"github.com/mvcris/maya-guessr/backend/internal/interfaces/http/dtos"
	"gorm.io/gorm"
)

type AuthHandler struct {
	loginUseCase *auth.LoginUseCase
	router *gin.Engine
	db *gorm.DB
}

func NewAuthHandler(db *gorm.DB, router *gin.Engine) *AuthHandler {
	userRepository := repositories.NewUserPgRepository(db)
	refreshTokenRepository := repositories.NewRefreshTokenPgRepository(db)
	jwtService := services.NewJwtService()
	return &AuthHandler{db: db, router: router, loginUseCase: auth.NewLoginUseCase(userRepository, refreshTokenRepository, jwtService)}
}

func (h *AuthHandler) Login(c *gin.Context) {
	var input dtos.LoginRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	output, err := h.loginUseCase.Execute(auth.LoginInput{
		Email: input.Email,
		Password: input.Password,
	})
	if err != nil {
		httppkg.RespondError(c, err)
		return
	}
	c.JSON(http.StatusOK, output)
}

func (h *AuthHandler) SetupRoutes() {
	h.router.POST("/auth/login", h.Login)
}