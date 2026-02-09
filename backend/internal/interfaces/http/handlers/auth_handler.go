package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	coreerrors "github.com/mvcris/maya-guessr/backend/internal/core/errors"
	"github.com/mvcris/maya-guessr/backend/internal/core/services"
	"github.com/mvcris/maya-guessr/backend/internal/core/use_cases/auth"
	"github.com/mvcris/maya-guessr/backend/internal/core/use_cases/user"
	"github.com/mvcris/maya-guessr/backend/internal/infrastructure/gorm/repositories"
	httppkg "github.com/mvcris/maya-guessr/backend/internal/interfaces/http"
	"github.com/mvcris/maya-guessr/backend/internal/interfaces/http/dtos"
	"github.com/mvcris/maya-guessr/backend/internal/interfaces/http/middleware"
	"gorm.io/gorm"
)

type AuthHandler struct {
	loginUseCase  *auth.LoginUseCase
	getMeUseCase  *user.GetMeUseCase
	jwtService    *services.JwtService
	router        *gin.Engine
	db            *gorm.DB
}

func NewAuthHandler(db *gorm.DB, router *gin.Engine) *AuthHandler {
	userRepository := repositories.NewUserPgRepository(db)
	refreshTokenRepository := repositories.NewRefreshTokenPgRepository(db)
	jwtService := services.NewJwtService()
	return &AuthHandler{
		db:            db,
		router:        router,
		loginUseCase:  auth.NewLoginUseCase(userRepository, refreshTokenRepository, jwtService),
		getMeUseCase:  user.NewGetMeUseCase(userRepository),
		jwtService:    jwtService,
	}
}

func (h *AuthHandler) GetMe(c *gin.Context) {
	userID, ok := middleware.GetAuthenticatedUserID(c)
	if !ok {
		httppkg.RespondError(c, coreerrors.Unauthorized("invalid or missing token"))
		return
	}
	u, err := h.getMeUseCase.Execute(userID)
	if err != nil {
		httppkg.RespondError(c, err)
		return
	}
	c.JSON(http.StatusOK, dtos.CreateUserResponse{
		ID:        u.ID,
		Name:      u.Name,
		Email:     u.Email,
		Username:  u.Username,
		CreatedAt: u.CreatedAt,
	})
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
	authGroup := h.router.Group("/auth")
	authGroup.POST("/login", h.Login)
	authGroup.GET("/me", middleware.AuthMiddleware(h.jwtService), h.GetMe)
}