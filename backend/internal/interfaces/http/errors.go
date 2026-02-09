package http

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	coreerrors "github.com/mvcris/maya-guessr/backend/internal/core/errors"
)

// RespondError writes an appropriate HTTP response for err.
// If err is a domain error (implements HTTP status), the corresponding status and message are returned.
// Otherwise, the error is logged and a 500 response with a generic message is returned.
func RespondError(c *gin.Context, err error) {
	if err == nil {
		return
	}
	if status, ok := coreerrors.Status(err); ok {
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}
	log.Printf("unhandled error: %v", err)
	c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
}
