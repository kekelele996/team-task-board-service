package middleware

import (
	"log/slog"
	"net/http"

	"github.com/gbkanban/gbkanban/internal/constants"
	"github.com/gbkanban/gbkanban/internal/response"
	"github.com/gin-gonic/gin"
)

// Recovery converts panics into JSON 500 responses.
func Recovery(logger *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				logger.Error("panic recovered", "panic", r)
				response.Error(c, http.StatusInternalServerError, constants.CodeInternal, "internal server error")
			}
		}()
		c.Next()
	}
}
