package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	coreerrors "github.com/mvcris/maya-guessr/backend/internal/core/errors"
	"github.com/mvcris/maya-guessr/backend/internal/core/services"
	httppkg "github.com/mvcris/maya-guessr/backend/internal/interfaces/http"
)

type contextKey string

const UserIDKey contextKey = "user_id"

// UserIDContextKey is the string key used for storing user_id in Gin context.
const UserIDContextKey = "user_id"

// GetAuthenticatedUserID returns the authenticated user's ID from the request context, set by AuthMiddleware.
// Returns (userID, true) if present, ("", false) otherwise.
func GetAuthenticatedUserID(c *gin.Context) (string, bool) {
	userIDVal, ok := c.Get(UserIDContextKey)
	if !ok {
		return "", false
	}
	userID, ok := userIDVal.(string)
	if !ok {
		return "", false
	}
	return userID, true
}

func AuthMiddleware(jwtService *services.JwtService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			httppkg.RespondError(c, coreerrors.Unauthorized("invalid or missing token"))
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			httppkg.RespondError(c, coreerrors.Unauthorized("invalid or missing token"))
			c.Abort()
			return
		}

		tokenString := parts[1]
		claims, err := jwtService.ValidateAccessToken(tokenString)
		if err != nil {
			httppkg.RespondError(c, coreerrors.Unauthorized("invalid or missing token"))
			c.Abort()
			return
		}

		c.Set(UserIDContextKey, claims.UserId)
		c.Next()
	}
}
